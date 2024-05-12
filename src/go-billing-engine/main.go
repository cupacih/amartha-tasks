package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"gobillingengine/engine"
	"gobillingengine/model"
	"os"
)

//
//// Loan represents a loan with its details.
//type Loan struct {
//	LoanID           string
//	Amount           float64 // Total loan amount
//	FlatInterestRate float64
//	Weeks            int
//}
//
//// NewLoan creates a new loan instance.
//func NewLoan(loanID string, weeks int, amount, flatInterestRate float64) *Loan {
//	return &Loan{
//		LoanID:           loanID,
//		Amount:           amount,
//		Weeks:            weeks,
//		FlatInterestRate: flatInterestRate,
//	}
//}
//
//type Payment struct {
//	Week   int
//	Amount float64
//}
//
//type Billing struct {
//	Loan           *Loan
//	PayableAmount  float64      // PayableAmount of each payment needs to be made
//	Outstanding    float64      // Outstanding balance of the loan
//	MissedPayment  int          // MissedPayment number of continuous missed payments (Some customers may miss repayments, If they miss 2 continuous repayments they are delinquent borrowers.)
//	RemainingWeeks int          // RemainingWeeks of outstanding amoutn
//	PaymentRecord  stacks.Stack // PaymentRecord of paid load
//}
//
//func NewBilling(loan *Loan) *Billing {
//	outstanding := loan.Amount * loan.FlatInterestRate
//	payableAmount := outstanding / float64(loan.Weeks)
//
//	return &Billing{
//		Loan:           loan,
//		PayableAmount:  payableAmount,
//		Outstanding:    outstanding,
//		MissedPayment:  0,
//		RemainingWeeks: loan.Weeks,
//		PaymentRecord:  lls.New(),
//	}
//}
//
//// GenerateLoanSchedule generates the billing schedule for the loan.
//func (b *Billing) GenerateLoanSchedule() *lls.Stack {
//	loanSchedule := lls.New()
//	for week := b.Loan.Weeks; week > 0; week-- {
//		loanSchedule.Push(&Payment{
//			Week:   week,
//			Amount: b.PayableAmount,
//		})
//	}
//	return loanSchedule
//}
//
//func (b *Billing) GenerateRemainingLoanSchedule() *lls.Stack {
//	loanSchedule := lls.New()
//	for week := b.Loan.Weeks; week >= (b.Loan.Weeks - b.RemainingWeeks); week-- {
//		loanSchedule.Push(&Payment{
//			Week:   week,
//			Amount: b.PayableAmount,
//		})
//	}
//	return loanSchedule
//}
//
//func (b *Billing) PrintPayment(s *lls.Stack) (string, error) {
//	var text string
//	for !s.Empty() {
//		if val, ok := s.Pop(); !ok {
//			errors.New("expect stack element is not empty")
//		} else {
//			bill := val.(*Payment)
//			text += fmt.Sprintf("Week: %d, Payable amount: %f\n", bill.Week, bill.Amount)
//		}
//	}
//	return text, nil
//}
//
//// GetOutstanding returns the current outstanding balance on the loan.
//func (b *Billing) GetOutstanding() float64 {
//	return b.Outstanding
//}
//
//// IsDelinquent checks if the borrower is delinquent (missed 2 continuous repayments).
//func (b *Billing) IsDelinquent() bool {
//	return b.MissedPayment >= 2
//}
//
//func (b *Billing) ResetMissedPayment() {
//	b.MissedPayment = 0
//}
//
//// MakePayment makes a payment of a certain amount on the loan.
//func (b *Billing) MakePayment(amount float64) error {
//	if amount > b.Outstanding {
//		return errors.New("payment exceeds outstanding amount")
//	}
//	if amount == b.PayableAmount {
//		b.makePayablePayment(amount)
//		return nil
//	}
//	if amount == 0 {
//		b.makeZeroPayment()
//		return nil
//	}
//	return errors.New(fmt.Sprintf("payment should be %f or 0", b.PayableAmount))
//}
//
//func (b *Billing) makeZeroPayment() {
//	b.MissedPayment += 1
//	b.recordPayment(0)
//}
//
//func (b *Billing) makePayablePayment(amount float64) {
//	b.Outstanding -= amount
//	b.ResetMissedPayment() // assumption: resetting the continuous missed payments when payment occur
//	b.RemainingWeeks -= 1
//	b.recordPayment(amount)
//}
//
//func (b *Billing) recordPayment(paidAmount float64) {
//	week := (b.Loan.Weeks - b.RemainingWeeks) + 1
//	payment := &Payment{
//		Week:   week,
//		Amount: paidAmount,
//	}
//	b.PaymentRecord.Push(payment)
//}

func runQuickScenario() {
	loanID := uuid.New().String()
	loanAmount := 5000000.00
	flatInterestRate := 10.0
	weeks := 50

	// Create a new loan & billing
	loan := model.NewLoan(loanID, weeks, loanAmount, flatInterestRate)
	billing := engine.NewBilling(loan)

	// Generate billing schedule
	billingSchedule := billing.GenerateLoanSchedule()
	s, err := billing.PrintPayment(billingSchedule)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nGenerateLoanSchedule: \n %s\n", s)

	var payments []float64
	payments = append(payments, 1000000.000000)
	payments = append(payments, 1000000.000000)
	payments = append(payments, 1000000.000000)
	payments = append(payments, 1000000.000000)
	payments = append(payments, 0)
	payments = append(payments, 200.00)

	for _, p := range payments {
		err := billing.MakePayment(p)
		if err != nil {
			fmt.Printf("could not make payment with amount: %f with err: %v\n", p, err)
		}
	}

	// Generate billing schedule
	remainingLoanSchedule := billing.GenerateRemainingLoanSchedule()
	s, err = billing.PrintPayment(remainingLoanSchedule)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nremaining loan schedule: \n %s\n", s)

	err = billing.MakePayment(0)
	if err != nil {
		panic(err)
	}
	err = billing.MakePayment(0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nIsDelinquent: %t\n", billing.IsDelinquent())

	err = billing.MakePayment(1000000.000000)
	fmt.Printf("\nIsDelinquent: %t\n", billing.IsDelinquent())
}

func main() {
	// runQuickScenario() // uncomment this to run a quick example scenario

	rootCmd := &cobra.Command{
		Use:   "myapp",
		Short: "A simple command-line application",
		Run: func(cmd *cobra.Command, args []string) {
			// This function will be executed when the root command is called
			fmt.Println("Welcome to myapp!")
		},
	}

	// Define a subcommand
	greetCmd := &cobra.Command{
		Use:   "greet",
		Short: "Print a greeting message",
		Run: func(cmd *cobra.Command, args []string) {
			// This function will be executed when the 'greet' subcommand is called
			name, _ := cmd.Flags().GetString("name")
			fmt.Printf("Hello, %s!\n", name)
		},
	}

	// Add flag to the subcommand
	greetCmd.Flags().StringP("name", "n", "", "The name to greet")

	// Add the subcommand to the root command
	rootCmd.AddCommand(greetCmd)

	// Execute the root command

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
