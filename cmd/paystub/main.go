// Package main contains the Google paystub analyzer.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/filmil/fintools/pkg/out"
	"github.com/filmil/fintools/pkg/xml"
	"github.com/golang/glog"
)

func i(s string) string {
	return fmt.Sprintf("Income:Personal:US:Google:%s", s)
}

func e(s string) string {
	return fmt.Sprintf("Expenses:Personal:Google:%s", s)
}

func t(s string) string {
	return fmt.Sprintf("Expenses:Personal:Taxes:Y%%s:%s", s)
}

var (
	dateOnly = flag.Bool("date-only", false, "If set, prints only the statement date")
)

func setFlags() {
	flag.StringVar(&cfg.NetPay, "net-pay", "Assets:Personal:BofA:Checking", "Net pay label")
	flag.StringVar(&cfg.RegularPay, "regular-pay", i("RegularPay"), "Regular pay label")
	flag.StringVar(&cfg.AnnualBonus, "annual-bonus", i("AnnualBonus"), "")
	flag.StringVar(&cfg.IGroupTermLife, "income-group-term-life", i("GroupTermLife"), "")
	flag.StringVar(&cfg.PeerBonus, "peer-bonus", i("PeerBonus"), "")
	flag.StringVar(&cfg.GoogStockUnit, "goog-stock-unit", i("GoogStockUnit"), "")
	flag.StringVar(&cfg.SpotBonus, "spot-bonus", i("SpotBonus"), "")
	flag.StringVar(&cfg.GSUCRefund, "gsu-c-refund", i("GSUCRefund"), "")

	flag.StringVar(&cfg.Bonus401kPre, "bonus-401k-pre", e("Bonus401kPre"), "")
	flag.StringVar(&cfg.ClassCOffset, "class-c-offset", e("ClassCOffset"), "")
	flag.StringVar(&cfg.Dental, "dental", e("Dental"), "")
	flag.StringVar(&cfg.FSAHealth, "fsa-health", e("FsaHealth"), "")
	flag.StringVar(&cfg.EGroupTermLife, "expense-group-term-life", e("GroupTermLife"), "")
	flag.StringVar(&cfg.InternetReim, "internet-reim", e("InternetReim"), "")
	flag.StringVar(&cfg.LegalAccess, "legal-access", e("LegalAccess"), "")
	flag.StringVar(&cfg.LongTermDis, "long-term-disability", e("LongTermDisability"), "")
	flag.StringVar(&cfg.Medical, "medical", e("Medical"), "")
	flag.StringVar(&cfg.TransitPreTax, "transit-pretax", e("TransitPreTax"), "")
	flag.StringVar(&cfg.Vision, "vision", e("Vision"), "")
	flag.StringVar(&cfg.VolLifeEE, "vol-life-ee", e("VolLifeEe"), "")
	flag.StringVar(&cfg.VolLifeSpouse, "vol-life-spouse", e("VolLifeSpouse"), "")

	flag.StringVar(&cfg.FederalIncomeTax, "federal-income-tax", t("FederalIncomeTax"), "")
	flag.StringVar(&cfg.EmployeeMedicare, "employee-medicare", t("EmployeeMedicare"), "")
	flag.StringVar(&cfg.SocialSecurityEmployeeTax, "social-security-employee-tax", t("SocialSecurityEmployeeTax"), "")
	flag.StringVar(&cfg.CAStateIncomeTax, "ca-state-income-tax", t("CaStateIncomeTax"), "")
	flag.StringVar(&cfg.CAPrivateDisabilityEmployee, "ca-private-disability-employee", t("CAPrivateDisabilityEmployee"), "")
}

var (
	inputFile = flag.String("input", "", "Name of the file to examine")

	cfg out.Config
)

func main() {
	setFlags()
	flag.Parse()

	if *inputFile == "" {
		fmt.Fprintf(os.Stderr, "flag --input is required\n")
		os.Exit(-1)
	}

	file, err := os.Open(*inputFile)
	if err != nil {
		glog.Fatalf("could not open file: %v", err)
	}
	t, err := xml.Parse(file)
	if err != nil {
		glog.Fatalf("Parse: unexpected: %v", err)
	}
	if *dateOnly {
		fmt.Printf("%s\n", out.YMD(t.Date))
	} else if err := out.Output(t, cfg, os.Stdout); err != nil {
		glog.Fatalf("Output: unexpected: %v", err)
	}
}
