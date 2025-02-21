// Package buildium contains utilities to read Buildium CSV files.
package buildium

const (
	// Various index fields in the buildium "Income Statement Detailed" report.

	// Transaction identifier. Multiple lines may have the same indexID.
	IndexId = 0

	// Dollar amount.
	IndexAmount = 6

	// Account identifier.
	IndexAccountID = 8

	// Transaction date, as US string.
	IndexEntryDate = 2

	IndexJournalMemo     = 3
	IndexCheckNumber     = 4
	IndexPostingMemo     = 5
	IndexReferenceNumber = 9
	IndexBuildingId      = 10
	IndexPayeeName       = 16
	IndexPayeeNameRaw    = 17 // What is this?

	// "Income", "Expense"
	IndexAccountTypeName = 29

	// "3" => "4000 RENT INCOME" etc.
	IndexGlAccountId = IndexAccountID
	// "4000 RENT INCOME"
	IndexGlAccountName = 11
)
