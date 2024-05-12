package engine

import (
	"errors"
	"fmt"
	"github.com/emirpasic/gods/stacks"
	lls "github.com/emirpasic/gods/stacks/linkedliststack"
	"gobillingengine/model"
)

type Payment struct {
	Week   int
	Amount float64
}

type Billing struct {
	Loan           *model.Loan
	PayableAmount  float64      // PayableAmount of each payment needs to be made
	Outstanding    float64      // Outstanding balance of the loan
	MissedPayment  int          // MissedPayment number of continuous missed payments (Some customers may miss repayments, If they miss 2 continuous repayments they are delinquent borrowers.)
	RemainingWeeks int          // RemainingWeeks of outstanding amoutn
	PaymentRecord  stacks.Stack // PaymentRecord of paid load
}

func NewBilling(loan *model.Loan) *Billing {
	outstanding := loan.Amount * loan.FlatInterestRate
	payableAmount := outstanding / float64(loan.Weeks)

	return &Billing{
		Loan:           loan,
		PayableAmount:  payableAmount,
		Outstanding:    outstanding,
		MissedPayment:  0,
		RemainingWeeks: loan.Weeks,
		PaymentRecord:  lls.New(),
	}
}

// GenerateLoanSchedule generates the billing schedule for the loan.
func (b *Billing) GenerateLoanSchedule() *lls.Stack {
	loanSchedule := lls.New()
	for week := b.Loan.Weeks; week > 0; week-- {
		loanSchedule.Push(&Payment{
			Week:   week,
			Amount: b.PayableAmount,
		})
	}
	return loanSchedule
}

func (b *Billing) GenerateRemainingLoanSchedule() *lls.Stack {
	loanSchedule := lls.New()
	for week := b.Loan.Weeks; week > (b.Loan.Weeks - b.RemainingWeeks); week-- {
		loanSchedule.Push(&Payment{
			Week:   week,
			Amount: b.PayableAmount,
		})
	}
	return loanSchedule
}

func (b *Billing) PrintPayment(s *lls.Stack) (string, error) {
	var text string
	for !s.Empty() {
		if val, ok := s.Pop(); !ok {
			errors.New("expect stack element is not empty")
		} else {
			bill := val.(*Payment)
			text += fmt.Sprintf("Week: %d, Payable amount: %f\n", bill.Week, bill.Amount)
		}
	}
	return text, nil
}

// GetOutstanding returns the current outstanding balance on the loan.
func (b *Billing) GetOutstanding() float64 {
	return b.Outstanding
}

// IsDelinquent checks if the borrower is delinquent (missed 2 continuous repayments).
func (b *Billing) IsDelinquent() bool {
	return b.MissedPayment >= 2
}

func (b *Billing) ResetMissedPayment() {
	b.MissedPayment = 0
}

// MakePayment makes a payment of a certain amount on the loan.
func (b *Billing) MakePayment(amount float64) error {
	if amount > b.Outstanding {
		return errors.New("payment exceeds outstanding amount")
	}
	if amount == b.PayableAmount {
		b.makePayablePayment(amount)
		return nil
	}
	if amount == 0 {
		b.makeZeroPayment()
		return nil
	}
	return errors.New(fmt.Sprintf("payment should be %f or 0", b.PayableAmount))
}

func (b *Billing) makeZeroPayment() {
	b.MissedPayment += 1
	b.recordPayment(0)
}

func (b *Billing) makePayablePayment(amount float64) {
	b.Outstanding -= amount
	b.ResetMissedPayment() // assumption: resetting the continuous missed payments when payment occur
	b.RemainingWeeks -= 1
	b.recordPayment(amount)
}

func (b *Billing) recordPayment(paidAmount float64) {
	week := (b.Loan.Weeks - b.RemainingWeeks) + 1
	payment := &Payment{
		Week:   week,
		Amount: paidAmount,
	}
	b.PaymentRecord.Push(payment)
}
