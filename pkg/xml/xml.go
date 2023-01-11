// Package xml contains routines for parsing the paystub data
// from an xml format produced by pdf2txt.
package xml

import (
	goxml "encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/filmil/fintools/pkg/tx"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

const eps = 1e-6

// USDate parses a US format date from the date string.
func USDate(s string) (tx.DateOnly, error) {
	d, err := time.Parse("01/02/2006", s)
	if err != nil {
		return tx.DateOnly{}, fmt.Errorf("unparseable date: %q", s)
	}
	return tx.DateOnly(d), nil
}

// storeAmounts stores the amounts in amountCol based on the row label values stored
// in textCol, using the setter function provided.  Returns error if it could not
// set anything.
//
// Example:
//
//	textCol       amountCol
//	Some Text     $1,123.44
//	Some other    $2,345.66
//	text
//
// (note there is no value to the right of "text", so it is joined to be "Some other text")
func storeAmounts(textCol, amountCol []Textline, setter func(string, tx.USD) error) error {
	// Current amount line index.
	var amLine int
	glog.V(3).Infof("textCol=%+v\n", textCol)
	glog.V(3).Infof("amountCol=%+v\n", amountCol)
	glog.V(3).Infof("len(textCol)=%v; len(amountCol)=%v\n", len(textCol), len(amountCol))
	for i := 0; i < len(textCol); i++ {
		pt := textCol[i]
		// Try to collect the lines that don't have any value associated
		name := pt.Text()
		c := amLine + 1
		glog.V(5).Infof("c=%v\n", c)
		for j := i + 1; j < len(textCol); j++ {
			// Keep appending lines that have no associated amount.  We determine
			// which lines have no associated amount by checking if there is a
			// value bounding box to the right of the line.
			tb := textCol[j].BBox.ExtendRight()
			// If there is an amount to the right of the line, then that line
			// should not be part of the previous one, and break.  Similarly,
			// if there are no more amounts.
			if c >= len(amountCol) {
				// Handles the case where the multiline is at the end of
				// report:
				// Some Text  $1
				// Some other $2
				// text
				name = strings.Join([]string{name, textCol[j].Text()}, " ")
				i++
				break
			}
			ab := amountCol[c].BBox
			r := IntersectingBBox(tb, ab)
			if r {
				break
			}
			name = strings.Join([]string{name, textCol[j].Text()}, " ")
			i++
		}
		valueStr := amountCol[amLine].Text()
		val, err := parseAmount(valueStr)
		if err != nil {
			return errors.Wrapf(err, "while parsing %q for row %q", valueStr, name)
		}
		if err := setter(name, val); err != nil {
			return errors.Wrapf(err, "could not set value %v for row %q", val, name)
		}
		amLine++
	}
	return nil
}

// Convert turns a Paystub parsed XML into a Transaction.
// TODO(filmil): This function is horrifying. Refactor.
func Convert(p Paystub) (tx.Transaction, error) {
	var t tx.Transaction

	tls := Textlines(p.Pages[0])

	// First establish some anchors.
	payDateTxt, err := FindOneTL(tls, "Pay Date")
	if err != nil {
		return t, errors.Wrapf(err, "while finding date")
	}

	earningsTl, err := FindOneTL(tls, "Earnings")
	if err != nil {
		return t, errors.Wrapf(err, "could not find earnings")
	}

	// Find the line containing "Taxes" which is below the "earnings" label.
	taxesTl, err := FindOneTLInBBox(tls, "Taxes", earningsTl.BBox.ExtendBottom())
	if err != nil {
		return t, errors.Wrapf(err, "could not find taxes")
	}

	deductionsTl, err := FindOneTLInBBox(tls, "Deductions", taxesTl.BBox.ExtendTop().ExtendRight())
	if err != nil {
		return t, errors.Wrapf(err, "could not find deductions")
	}

	paidTimeOff, err := FindOneTL(tls, "Paid Time Off")
	if err != nil {
		return t, errors.Wrapf(err, "could not find paid time off")
	}

	// Next make boxes
	// Earnings box
	earningsB := earningsTl.BBox.ExtendRight().ExtendBottom()
	earningsB.Right = deductionsTl.BBox.Left // Limit from the left
	earningsB.Bottom = taxesTl.BBox.Top + eps

	// Deductions box
	deductionsB := deductionsTl.BBox.ExtendRight().ExtendBottom()
	deductionsB.Bottom = taxesTl.BBox.Top + eps

	// Taxes box
	taxesB := taxesTl.BBox.ExtendRight().ExtendBottom()
	taxesB.Bottom = paidTimeOff.BBox.Top + eps

	// Parse from the "Pay Statement" box.
	// Now try parsing some stuff out.
	rightOfPayDate := payDateTxt.BBox.RightOf()
	payDate, err := OneTextline(SortLeft(FindInBBox(tls, rightOfPayDate)))
	if err != nil {
		return t, errors.Wrapf(err, "could not find date in %v", payDate)
	}
	t.Date, err = USDate(payDate.Text())
	if err != nil {
		return t, errors.Wrapf(err, "while parsing date")
	}

	documentTl, err := FindOneTL(tls, "Document")
	if err != nil {
		return t, errors.Wrapf(err, "while looking for document")
	}
	rightOfDocumentTl := documentTl.BBox.RightOf()
	docNumTl, err := OneTextline(SortLeft(FindInBBox(tls, rightOfDocumentTl)))
	if err != nil {
		return t, errors.Wrapf(err, "no docnumTL")
	}
	t.DocNum = docNumTl.Text()

	// Parse from the "Earnings" box.
	// Narrow only to the earnings textlines.
	earningsTls := FindInBBox(tls, earningsB)

	totalHoursWorkedTl, err := FindOneTLPrefix(tls, "Total Hours Worked")
	if err != nil {
		return t, errors.Wrapf(err, "could not find Total Hours Worked")
	}
	earningsB.Bottom = totalHoursWorkedTl.BBox.Top + eps

	earningsPayTypeTl, err := FindOneTL(earningsTls, "Pay Type")
	if err != nil {
		return t, errors.Wrapf(err, "could not find pay type in earnings box")
	}

	// Get a box that extends from pay type down to total hours worked
	payTypesBox := earningsPayTypeTl.BBox
	payTypesBox.Top = payTypesBox.Bottom - eps
	payTypesBox.Bottom = totalHoursWorkedTl.BBox.Top + eps

	currentTl, err := FindOneTL(earningsTls, "Current")
	if err != nil {
		return t, errors.Wrapf(err, "could not find current in earnings box")
	}

	// Get a box that extends from Current down to total hours worked
	currentBox := currentTl.BBox
	currentBox.Top = currentBox.Bottom - eps
	currentBox.Bottom = totalHoursWorkedTl.BBox.Top + eps

	payTypesTls := SortTop(FindInBBox(earningsTls, payTypesBox))
	currentTls := SortTop(FindInBBox(earningsTls, currentBox))
	if err := storeAmounts(payTypesTls, currentTls, t.IncomeByName); err != nil {
		return t, errors.Wrapf(err, "while setting income")
	}

	// Parse net pay.  Net pay amount is right of the label "Net Pay" which is
	// located above the fixed "Earnings" label.
	netPayTl, err := FindOneTLInBBox(tls, "Net Pay", earningsTl.BBox.ExtendTop().ExtendRight())
	if err != nil {
		return t, errors.Wrapf(err, "while finding Net Pay in the upper right corner")
	}
	netPayAmtTl, err := OneTextline(FindInBBox(tls, netPayTl.BBox.RightOf()))
	if err != nil {
		return t, errors.Wrapf(err, "could not find amount right of 'Net pay'")
	}
	t.NetPay, err = parseAmount(netPayAmtTl.Text())
	if err != nil {
		return t, errors.Wrapf(err, "could not set Net Pay")
	}

	// Parse deductions.  Uses deductionB which was established way up.
	// For deduction name, we look for strings under a fixed label "Deduction".
	// For contribution, we look for strings under label "Current" which is
	// in itself under label "Employee" in the deduction box.

	// Only the textlines in the deductions box.
	dedsTls := FindInBBox(tls, deductionsB)

	// Find the textline "Deduction"
	dedTl, err := FindOneTL(dedsTls, "Deduction")
	if err != nil {
		return t, errors.Wrapf(err, "could not find Deduction in Deductions")
	}

	// Parse Deductions->Employee->Current
	// Find the textline "Employee"
	employeeTl, err := FindOneTL(dedsTls, "Employee")
	if err != nil {
		return t, errors.Wrapf(err, "could not find Employee in Deductions")
	}
	// Find the textline "Current" below Employee.
	employeeCurrentTl, err := FindOneTLInBBox(dedsTls,
		"Current",
		employeeTl.BBox.ExtendDownTo(deductionsB.Bottom).Below(employeeTl.BBox))
	if err != nil {
		return t, errors.Wrapf(err, "could not find Current below Employee in Deductions")
	}

	// Isolate the column where we expect the deduction texts to be, and get
	// the text column.
	dedsCol := dedTl.BBox.ExtendDownTo(deductionsB.Bottom).Below(dedTl.BBox)
	dedColTls := SortTop(MatchPredicate(dedsTls, BindBBox(dedsCol, IntersectingBBoxTextline)))

	// Isolate the column where we expect the current deds amounts to be.
	dedsAmountsCol := employeeCurrentTl.BBox.
		ExtendDownTo(deductionsB.Bottom).
		Below(employeeCurrentTl.BBox)
	dedAmountColTls := SortTop(FindInBBox(dedsTls, dedsAmountsCol))
	if err := storeAmounts(dedColTls, dedAmountColTls, t.ExpenseByName); err != nil {
		return t, errors.Wrapf(err, "while parsing expenses")
	}

	// Parse Deductions->Employer->Current
	employerTl, err := FindOneTL(dedsTls, "Employer")
	if err != nil {
		return t, errors.Wrapf(err, "could not find Employer in Deductions")
	}

	// Find the textline "Current" below Employer.  Sadly, that text is
	// "merged" with the "YTD" from previous column and only SOMETIMES., for no
	// reason that I can see and understand.
	employerTlBBox := employerTl.BBox.ExtendDownTo(deductionsB.Bottom).Below(employerTl.BBox)
	employerCurrentTl, err := OneTextline(MatchPredicate(dedsTls,
		BindText("Current", MatchingSuffix),
		BindBBox(employerTlBBox, IntersectingBBoxTextline)))
	if err != nil {
		return t, errors.Wrapf(err, "could not find Current below Employer in Deductions")
	}

	emplDedsAmountsCol := employerCurrentTl.BBox.
		ExtendDownTo(deductionsB.Bottom).
		Below(employerCurrentTl.BBox)
	emplDedAmountColTls := SortTop(
		MatchPredicate(
			dedsTls,
			BindBBox(emplDedsAmountsCol, IntersectingBBoxTextline),
			BindBBox(employerTlBBox, IntersectingBBoxTextline)))
	if err := storeAmounts(dedColTls, emplDedAmountColTls, t.EmployerExpenseByName); err != nil {
		return t, errors.Wrapf(err, "while parsing employer deductions")
	}

	// Parse taxes. Uses taxesB whihc was established way up.  For tax name,
	// we look for strings under a fixed label "Tax" in this section. For the
	// tax amounts, we look for the amounts under a fixed label "Current" in
	// this section.

	// Limit to the textlines in this section.
	taxesTls := FindInBBox(tls, taxesB)
	taxTl, err := FindOneTL(taxesTls, "Tax")
	if err != nil {
		return t, errors.Wrapf(err, "could not find Tax in Taxes box")
	}
	taxCurrentTl, err := FindOneTL(taxesTls, "Current")
	if err != nil {
		return t, errors.Wrapf(err, "could not find Current in Taxes box")
	}

	// Find the column for the taxes labels.
	taxColB := taxTl.BBox.ExtendDownTo(taxesB.Bottom).Below(taxTl.BBox)
	taxAmountB := taxCurrentTl.BBox.ExtendDownTo(taxesB.Bottom).Below(taxCurrentTl.BBox)

	taxColTls := SortTop(FindInBBox(taxesTls, taxColB))
	taxAmountTls := SortTop(FindInBBox(taxesTls, taxAmountB))

	if err := storeAmounts(taxColTls, taxAmountTls, t.ExpenseByName); err != nil {
		return t, errors.Wrapf(err, "while setting tax expenses")
	}
	return t, nil
}

func parseAmount(s string) (tx.USD, error) {
	glog.V(3).Infof("parseAmount(%v)", s)
	r := strings.NewReplacer("$", "", ",", "", "(", "", ")", "")
	t := r.Replace(s)
	f, err := strconv.ParseFloat(t, 64)
	if len(s) == 0 {
		return tx.USD(0), fmt.Errorf("Nothing to parse: %q", s)
	}
	// Negative amounts are denoted like so: ($10).  Makes no sense, but here
	// we are.
	if s[0] == '(' {
		f = -f
	}
	return tx.USD(f), err

}

// Decode decodes the reader into Paystub data.
func Decode(r io.Reader) (Paystub, error) {
	var p Paystub
	d := goxml.NewDecoder(r)
	if err := d.Decode(&p); err != nil {
		return p, errors.Wrapf(err, "while parsing XML")
	}
	return p, nil
}

// Parse parses passed reader into a Paystub
func Parse(r io.Reader) (tx.Transaction, error) {
	var t tx.Transaction
	p, err := Decode(r)
	if err != nil {
		return t, errors.Wrapf(err, "while decoding input to Paystub")
	}
	t, err = Convert(p)
	if err != nil {
		return t, errors.Wrapf(err, "while converting Paystub to Transaction")
	}
	return t, nil
}
