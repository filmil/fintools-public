// Package out contains the code that outputs transaction information in the
// format compatible with beancount.
package out

import (
	"io"
	"text/template"
	"time"

	"github.com/filmil/fintools-public/pkg/tx"
)

// Employer contains text labels for Deductions->Employer
type Employer struct {
	GSUCRefund string
	// Expenses
	Bonus401kPre   string
	ClassCOffset   string
	Dental         string
	FSAHealth      string
	EGroupTermLife string
	InternetReim   string
	LegalAccess    string
	LongTermDis    string
	Medical        string
	TransitPreTax  string
	Vision         string
	VolLifeEE      string
	VolLifeSpouse  string
}

// Config contains the settings.
type Config struct {
	// Incomes
	NetPay         string
	RegularPay     string
	AnnualBonus    string
	IGroupTermLife string
	PeerBonus      string
	GoogStockUnit  string
	SpotBonus      string
	GSUCRefund     string
	// Expenses
	Bonus401kPre   string
	ClassCOffset   string
	Dental         string
	FSAHealth      string
	EGroupTermLife string
	InternetReim   string
	LegalAccess    string
	LongTermDis    string
	Medical        string
	TransitPreTax  string
	Vision         string
	VolLifeEE      string
	VolLifeSpouse  string
	// Taxes
	FederalIncomeTax            string
	EmployeeMedicare            string
	SocialSecurityEmployeeTax   string
	CAStateIncomeTax            string
	CAPrivateDisabilityEmployee string

	Employer Employer
}

// Out is the structure used to output transaction intormation.
type Out struct {
	T tx.Transaction
	C Config
}

// YMD formats time in the ISO YYYY-MM-DD format.
func YMD(t tx.DateOnly) string {
	tt := time.Time(t)
	return tt.Format("2006-01-02")
}

func year(t tx.DateOnly) string {
	tt := time.Time(t)
	return tt.Format("2006")
}

// The weird formatting is so that we can omit zero items without messing up
// the output.
var outTpl = template.Must(template.New("tx").Funcs(
	template.FuncMap{
		"ymd":  YMD,
		"year": year,
	},
).Parse(`{{ymd .T.Date}} ! "GOOGLE LLC Payroll {{.T.DocNum}}"{{if .T.RegularPay}}
   {{.C.RegularPay}} -{{.T.RegularPay}}{{end}}{{if .T.IGroupTermLife}}
   {{.C.IGroupTermLife}} -{{.T.IGroupTermLife}}{{end}}{{if .T.AnnualBonus}}
   {{.C.AnnualBonus}} -{{.T.AnnualBonus}}{{end}}{{if .T.PeerBonus}}
   {{.C.PeerBonus}} -{{.T.PeerBonus}}{{end}}{{if .T.GoogStockUnit}}
   {{.C.GoogStockUnit}} -{{.T.GoogStockUnit}}{{end}}{{if .T.SpotBonus}}
   {{.C.SpotBonus}} -{{.T.SpotBonus}}{{end}}{{if .T.GSUCRefund}}
   {{.C.GSUCRefund}} -{{.T.GSUCRefund}}{{end}}{{if .T.Bonus401kPre}}
   {{.C.Bonus401kPre}} {{.T.Bonus401kPre}}{{end}}{{if .T.ClassCOffset}}
   {{.C.ClassCOffset}} {{.T.ClassCOffset}}{{end}}{{if .T.Dental}}
   {{.C.Dental}} {{.T.Dental}}{{end}}{{if .T.FSAHealth}}
   {{.C.FSAHealth}} {{.T.FSAHealth}}{{end}}{{if .T.EGroupTermLife}}
   {{.C.EGroupTermLife}} {{.T.EGroupTermLife}}{{end}}{{if .T.InternetReim}}
   {{.C.InternetReim}} {{.T.InternetReim}}{{end}}{{if .T.LegalAccess}}
   {{.C.LegalAccess}} {{.T.LegalAccess}}{{end}}{{if .T.LongTermDis}}
   {{.C.LongTermDis}} {{.T.LongTermDis}}{{end}}{{if .T.Medical}}
   {{.C.Medical}} {{.T.Medical}}{{end}}{{if .T.TransitPreTax}}
   {{.C.TransitPreTax}} {{.T.TransitPreTax}}{{end}}{{if .T.Vision}}
   {{.C.Vision}} {{.T.Vision}}{{end}}{{if .T.VolLifeEE}}
   {{.C.VolLifeEE}} {{.T.VolLifeEE}}{{end}}{{if .T.VolLifeSpouse}}
   {{.C.VolLifeSpouse}} {{.T.VolLifeSpouse}}{{end}}{{if .T.FederalIncomeTax}}
   {{year .T.Date | printf .C.FederalIncomeTax}} {{.T.FederalIncomeTax}}{{end}}{{if .T.EmployeeMedicare}}
   {{year .T.Date | printf .C.EmployeeMedicare}} {{.T.EmployeeMedicare}}{{end}}{{if .T.SocialSecurityEmployeeTax}}
   {{year .T.Date | printf .C.SocialSecurityEmployeeTax}} {{.T.SocialSecurityEmployeeTax}}{{end}}{{if .T.CAStateIncomeTax}}
   {{year .T.Date | printf .C.CAStateIncomeTax}} {{.T.CAStateIncomeTax}}{{end}}{{if .T.CAPrivateDisabilityEmployee}}
   {{year .T.Date | printf .C.CAPrivateDisabilityEmployee}} {{.T.CAPrivateDisabilityEmployee}}{{end}}{{if .T.NetPay}}
   {{.C.NetPay}} {{.T.NetPay}}{{end}}
`))

func Output(t tx.Transaction, cfg Config, w io.Writer) error {
	o := Out{T: t, C: cfg}
	if err := outTpl.Execute(w, o); err != nil {
		return err
	}
	return nil
}
