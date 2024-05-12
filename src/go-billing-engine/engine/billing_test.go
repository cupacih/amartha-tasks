package engine

import (
	"errors"
	"testing"

	lls "github.com/emirpasic/gods/stacks/linkedliststack"
	"github.com/stretchr/testify/assert"

	"gobillingengine/model"
)

func TestNewBilling(t *testing.T) {
	loan := &model.Loan{
		LoanID:           "1001",
		Amount:           5000,
		FlatInterestRate: 0.1,
		Weeks:            50,
	}
	billing := NewBilling(loan)

	assert.NotNil(t, billing)
	assert.Equal(t, loan, billing.Loan)
	assert.Equal(t, 5000*0.1, billing.Outstanding)
	assert.Equal(t, 5000*0.1/float64(50), billing.PayableAmount)
	assert.Zero(t, billing.MissedPayment)
	assert.Equal(t, loan.Weeks, billing.RemainingWeeks)
	assert.NotNil(t, billing.PaymentRecord)
	assert.IsType(t, &lls.Stack{}, billing.PaymentRecord)
}

func TestGenerateLoanSchedule(t *testing.T) {
	loan := &model.Loan{
		LoanID:           "1001",
		Amount:           5000,
		FlatInterestRate: 0.1,
		Weeks:            50,
	}
	billing := NewBilling(loan)
	loanSchedule := billing.GenerateLoanSchedule()

	assert.NotNil(t, loanSchedule)
	assert.Equal(t, loan.Weeks, loanSchedule.Size())
}

func TestGenerateRemainingLoanSchedule(t *testing.T) {
	loan := &model.Loan{
		LoanID:           "1001",
		Amount:           5000,
		FlatInterestRate: 0.1,
		Weeks:            50,
	}
	billing := NewBilling(loan)
	billing.RemainingWeeks = 10
	loanSchedule := billing.GenerateRemainingLoanSchedule()

	assert.NotNil(t, loanSchedule)
	assert.Equal(t, 10, loanSchedule.Size())
}

func TestPrintPayment(t *testing.T) {
	payment := &Payment{
		Week:   1,
		Amount: 100,
	}
	stack := lls.New()
	stack.Push(payment)

	billing := &Billing{}
	text, err := billing.PrintPayment(stack)

	assert.NoError(t, err)
	assert.Contains(t, text, "Week: 1")
	assert.Contains(t, text, "100")
}

func TestMakePayment(t *testing.T) {
	loan := &model.Loan{
		LoanID:           "1001",
		Amount:           5000,
		FlatInterestRate: 0.1,
		Weeks:            50,
	}
	billing := NewBilling(loan)

	// Test case 1: Payment exceeds outstanding amount
	err := billing.MakePayment(6000)
	assert.Error(t, err)
	assert.Equal(t, errors.New("payment exceeds outstanding amount"), err)

	// Test case 2: Invalid payment amount
	err = billing.MakePayment(20)
	assert.Error(t, err)
	assert.Equal(t, errors.New("payment should be 10.000000 or 0"), err)

	// Test case 3: Payable payment
	err = billing.MakePayment(billing.PayableAmount)
	assert.NoError(t, err)

	// Test case 4: Zero payment
	err = billing.MakePayment(0)
	assert.NoError(t, err)
}

func TestMakePayment_MultiplePayments(t *testing.T) {
	loan := &model.Loan{
		LoanID:           "1001",
		Amount:           5000,
		FlatInterestRate: 0.1,
		Weeks:            5,
	}
	billing := NewBilling(loan)

	// Make multiple payments
	paymentAmount := billing.PayableAmount
	for i := 0; i < loan.Weeks; i++ {
		err := billing.MakePayment(paymentAmount)
		assert.NoError(t, err)
	}

	// Verify outstanding amount is zero
	assert.Equal(t, float64(0), billing.Outstanding)
	assert.Equal(t, 0, billing.RemainingWeeks)

	// Verify that making additional payment fails
	err := billing.MakePayment(paymentAmount)
	assert.Error(t, err)
	assert.Equal(t, errors.New("payment exceeds outstanding amount"), err)
}

func TestIsDelinquent(t *testing.T) {
	// Test case 1: Not delinquent
	billing := &Billing{}
	assert.False(t, billing.IsDelinquent())

	// Test case 2: Delinquent
	billing.MissedPayment = 2
	assert.True(t, billing.IsDelinquent())
}
