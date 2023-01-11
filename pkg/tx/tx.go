// Package tx contains transaction parsing support.
package tx

import (
	"fmt"
	"time"
	"encoding/json"
	"io"

	"github.com/golang/glog"
)

// Employer accounts for employer contributions.
type Employer struct {
	GSUCRefund USD `json:",omitempty"`
	// Expenses
	Bonus401kPre   USD `json:",omitempty"`
	ClassCOffset   USD `json:",omitempty"`
	Dental         USD `json:",omitempty"`
	FSAHealth      USD `json:",omitempty"`
	EGroupTermLife USD `json:",omitempty"`
	InternetReim   USD `json:",omitempty"`
	LegalAccess    USD `json:",omitempty"`
	LongTermDis    USD `json:",omitempty"`
	Medical        USD `json:",omitempty"`
	TransitPreTax  USD `json:",omitempty"`
	Vision         USD `json:",omitempty"`
	VolLifeEE      USD `json:",omitempty"`
	VolLifeSpouse  USD `json:",omitempty"`
}

type DateOnly time.Time

func (d DateOnly) Equal(o DateOnly) bool {
	t1 := time.Time(d)
	t2 := time.Time(o)
	return t1.Equal(t2)
}

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("could not read date: %v", err)
	}
	p, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("while parsing date: %v: %v", s, err)
	}
	*d = DateOnly(p)
	return nil
}

// Transaction is a single transaction converted from a paystub.
type Transaction struct {
	// Date when payment was made
	Date DateOnly `json:",omitempty"`
	// DocNum is the document number (if any)
	DocNum string `json:",omitempty"`

	NetPay         USD `json:",omitempty"`
	RegularPay     USD `json:",omitempty"`
	AnnualBonus    USD `json:",omitempty"`
	IGroupTermLife USD `json:",omitempty"`
	PeerBonus      USD `json:",omitempty"`
	GoogStockUnit  USD `json:",omitempty"`
	SpotBonus      USD `json:",omitempty"`
	GSUCRefund     USD `json:",omitempty"`

	// Expenses
	Bonus401kPre   USD `json:",omitempty"`
	ClassCOffset   USD `json:",omitempty"`
	Dental         USD `json:",omitempty"`
	FSAHealth      USD `json:",omitempty"`
	EGroupTermLife USD `json:",omitempty"`
	InternetReim   USD `json:",omitempty"`
	LegalAccess    USD `json:",omitempty"`
	LongTermDis    USD `json:",omitempty"`
	Medical        USD `json:",omitempty"`
	TransitPreTax  USD `json:",omitempty"`
	Vision         USD `json:",omitempty"`
	VolLifeEE      USD `json:",omitempty"`
	VolLifeSpouse  USD `json:",omitempty"`

	// Taxes
	FederalIncomeTax            USD `json:",omitempty"`
	EmployeeMedicare            USD `json:",omitempty"`
	SocialSecurityEmployeeTax   USD `json:",omitempty"`
	CAStateIncomeTax            USD `json:",omitempty"`
	CAPrivateDisabilityEmployee USD `json:",omitempty"`


	Employer Employer `json:",omitempty"`
}

// Read gets a Transaction from a Reader.
func Read(r io.Reader) (Transaction, error) {
	d := json.NewDecoder(r)
	var tx Transaction

	if err := d.Decode(&tx); err != nil {
		return tx, fmt.Errorf("could not read from JSON file: %v", err)
	}
	return tx, nil
}

func (t *Transaction) IncomeByName(s string, amount USD) error {
	switch s {
	case "Annual Bonus":
		t.AnnualBonus += amount
	case "Group Term Life":
		t.IGroupTermLife += amount
	case "Peer Bonus":
		t.PeerBonus += amount
	case "Regular Pay":
		t.RegularPay += amount
	case "Goog Stock Unit":
		t.GoogStockUnit += amount
	case "Spot Bonus":
		t.SpotBonus += amount
	default:
		return fmt.Errorf("Income by name not found: %q to place amount %v", s, amount)
	}
	return nil
}

func (t *Transaction) ExpenseByName(s string, amount USD) error {
	glog.V(3).Infof("ExpenseByname: expense: %q, amount: %v", s, amount)
	switch s {
	case "Bonus 401K Pre":
		t.Bonus401kPre += amount
	case "Class C Offset":
		t.ClassCOffset += amount
	case "GSU C Refund":
		t.GSUCRefund += amount
	case "Dental":
		t.Dental += amount
	case "FSA Health":
		t.FSAHealth += amount
	case "Group Term Life":
		t.EGroupTermLife += amount
	case "Internet Reim":
		t.InternetReim += amount
	case "LegalAccess":
		t.LegalAccess += amount
	case "LongTerm Dis":
		t.LongTermDis += amount
	case "Medical":
		t.Medical += amount
	case "Transit PreTax":
		t.TransitPreTax += amount
	case "Vision":
		t.Vision += amount
	case "Vol Life EE":
		t.VolLifeEE += amount
	case "Vol Life Spouse":
		t.VolLifeSpouse += amount
	case "Federal Income Tax":
		t.FederalIncomeTax += amount
	case "Employee Medicare":
		t.EmployeeMedicare += amount
	case "Social Security Employee Tax":
		t.SocialSecurityEmployeeTax += amount
	case "CA State Income Tax":
		t.CAStateIncomeTax += amount
	case "CA Private Disability Employee":
		t.CAPrivateDisabilityEmployee += amount
	default:
		return fmt.Errorf("ExpenseByName: not found: %q to place amount %v", s, amount)
	}
	return nil
}

func (t *Transaction) EmployerExpenseByName(s string, amount USD) error {
	glog.V(3).Infof("EmployerExpenseByName: expense: %q, amount: %v", s, amount)
	switch s {
	case "Bonus 401K Pre":
		t.Employer.Bonus401kPre += amount
	case "Class C Offset":
		t.Employer.ClassCOffset += amount
	case "GSU C Refund":
		t.Employer.GSUCRefund += amount
	case "Dental":
		t.Employer.Dental += amount
	case "FSA Health":
		t.Employer.FSAHealth += amount
	case "Group Term Life":
		t.Employer.EGroupTermLife += amount
	case "Internet Reim":
		t.Employer.InternetReim += amount
	case "LegalAccess":
		t.Employer.LegalAccess += amount
	case "LongTerm Dis":
		t.Employer.LongTermDis += amount
	case "Medical":
		t.Employer.Medical += amount
	case "Transit PreTax":
		t.Employer.TransitPreTax += amount
	case "Vision":
		t.Employer.Vision += amount
	case "Vol Life EE":
		t.Employer.VolLifeEE += amount
	case "Vol Life Spouse":
		t.Employer.VolLifeSpouse += amount
	default:
		return fmt.Errorf("EmployerExpenseByName: not found: %q to place amount %v", s, amount)
	}
	return nil
}

// USD is the currency in this paystub.
type USD float64

func (v USD) String() string {
	return fmt.Sprintf("%.4f USD", v)
}

func (v *USD) UnmarshalJSON(b []byte) error {
	var f float64
	if err := json.Unmarshal(b, &f); err != nil {
		return fmt.Errorf("error parsing USD: %v", err)
	}
	*v = USD(f)
	return nil
}

// VACHR is the amount of vacation hours accrued.
type VACHR float64
