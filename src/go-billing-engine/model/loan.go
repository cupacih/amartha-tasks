package model

// Loan represents a loan with its details.
type Loan struct {
	LoanID           string
	Amount           float64 // Total loan amount
	FlatInterestRate float64
	Weeks            int
}

// NewLoan creates a new loan instance.
func NewLoan(loanID string, weeks int, amount, flatInterestRate float64) *Loan {
	return &Loan{
		LoanID:           loanID,
		Amount:           amount,
		Weeks:            weeks,
		FlatInterestRate: flatInterestRate,
	}
}
