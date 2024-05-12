package model

import (
	"testing"
)

func TestNewLoan(t *testing.T) {
	// Define test cases
	tests := []struct {
		loanID           string
		weeks            int
		amount           float64
		flatInterestRate float64
	}{
		{"1001", 50, 5000, 0.1},
		{"1002", 30, 3000, 0.05},
	}

	// Iterate over test cases
	for _, tc := range tests {
		t.Run(tc.loanID, func(t *testing.T) {
			// Create a new loan instance
			loan := NewLoan(tc.loanID, tc.weeks, tc.amount, tc.flatInterestRate)

			// Check loan ID
			if loan.LoanID != tc.loanID {
				t.Errorf("Expected loan ID %s, but got %s", tc.loanID, loan.LoanID)
			}

			// Check weeks
			if loan.Weeks != tc.weeks {
				t.Errorf("Expected weeks %d, but got %d", tc.weeks, loan.Weeks)
			}

			// Check amount
			if loan.Amount != tc.amount {
				t.Errorf("Expected amount %.2f, but got %.2f", tc.amount, loan.Amount)
			}

			// Check flat interest rate
			if loan.FlatInterestRate != tc.flatInterestRate {
				t.Errorf("Expected flat interest rate %.2f, but got %.2f", tc.flatInterestRate, loan.FlatInterestRate)
			}
		})
	}
}
