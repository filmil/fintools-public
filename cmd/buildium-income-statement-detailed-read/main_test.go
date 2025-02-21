package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/filmil/fintools-public/pkg/cfg"
	"github.com/filmil/fintools-public/pkg/csv2"
)

func TestConvert(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    string
		expected string
		schema   string
	}{
		{
			input: `id,journalCodeId,entryDate,journalMemo,checkNumber,postingMemo,amount,isCashPosting,glAccountId,referenceNumber,buildingId,GLAccountName,parentGLAccountName,isParent,GLSubTypeId,parentGLAccountId,payeeName,payeeNameRaw,payeeUserId,payeeNames,userTypeId,leaseId,unitNumber,unitId,buildingName,vendorId,fileId,isBankAccount,GLAccountTypeId,GLAccountTypeName,journalCodeDescription,AccountingBookId,isEFT,eftResultDate,eftResultCode,attributeId,reversedJournalId,BankAccountName,BuildingStatusId,depositJournalId,isActive,ParentIsActive,GLAccountNameRaw,period,sortedName,periodSortOrder,PeriodStartDate,PeriodEndDate
212729100,4,12/31/2024,by Jonathan Martinez,,by Jonathan Martinez,-1025.00,1,3,,625353,4000 Rent Income,,0,40,0,Unit 3 - Jonathan Martinez,Jonathan Martinez,5855978,Jonathan Martinez,1,2399056,3,1665680,300 Berry Avenue,625353,0,0,4,Income,Payment,,1,1/2/2025,S01,10001,0,,1,0,True,,4000 Rent Income,Q4-2024,,3,10/1/2024,12/31/2024
212729100,4,12/31/2024,by Jonathan Martinez,,by Jonathan Martinez,-30.00,1,120327,,625353,Utility Income,,0,40,0,Unit 3 - Jonathan Martinez,Jonathan Martinez,5855978,Jonathan Martinez,1,2399056,3,1665680,300 Berry Avenue,625353,0,0,4,Income,Payment,,1,1/2/2025,S01,10001,0,,1,0,True,,Utility Income,Q4-2024,,3,10/1/2024,12/31/2024
`,
			expected: `T,Date,Description,Category,Amount,Tags,Account,Account #,Institution,Month,Week,Transaction ID,Account ID,Check Number,Full Description,Date Added,Categorized Date
,12/31/2024,Unit 3 - Jonathan Martinez,Category 3,1025.00,Tax,manual:some_account_1,,,1/1/0001,12/31/0000,,manual:some_account_1,,,1/1/0001,
,12/31/2024,Unit 3 - Jonathan Martinez,Category 120327,30.00,Tax,manual:some_account_1,,,1/1/0001,12/31/0000,,manual:some_account_1,,,1/1/0001,
`,
			schema: `
            {
              "account_id": "manual:some_account_1",
              "account_map": [
                {
                  "original": "3",
                  "category": "Category 3"
                },
                {
                  "original": "120327",
                  "category": "Category 120327",
                  "id": "manual:8ca0919d2dc7bb8b9af46389ef987ae762c9b77c2913e3c205b04aaf697f3d07"
                }
              ]
            }
            `,
		},
	}

	for i, test := range tests {
		test := test
		i := i
		t.Run(fmt.Sprintf("test:%d", i), func(t *testing.T) {
			sr := strings.NewReader(test.schema)
			lsx, err := cfg.LoadSchema(sr)
			r := strings.NewReader(test.input)
			if err != nil {
				t.Fatalf("could not load JSON schema: %v", err)
			}
			cf := cfg.New(lsx)
			var b strings.Builder

			c, err := csv2.NewCSVData(r)
			if err != nil {
				t.Fatalf("could not read CSV: %v", err)
			}

			tRows, err := ConvertTiller(c, cf)
			if err != nil {
				t.Fatalf("could not convert tiller: %v", err)
			}
			if err != nil {
				t.Fatalf("could not read CSV: %v", err)
			}
			if err := WriteWriter(&b, tRows); err != nil {
				t.Fatalf("could not write CSV: %v", err)
			}
			actual := b.String()
			if test.expected != actual {
				t.Errorf("want: %v; \n\t got: %v", test.expected, actual)
			}
		})
	}

}
