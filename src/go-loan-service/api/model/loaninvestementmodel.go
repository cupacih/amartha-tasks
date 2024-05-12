package model

import (
	"context"
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	types "goloanservice/api/internal/type"
)

type ILoanInvestmentModel interface {
	TxInvestment(ctx context.Context, loan *types.Loan, investments *types.InvestLoanRequest) error
}

type LoanInvestmentModel struct {
	DB                  *sql.DB
	LoanTableName       string
	InvestmentTableName string
}

func NewLoanInvestmentModel(db *sql.DB, loanTableName, investmentTableName string) ILoanInvestmentModel {
	return &LoanInvestmentModel{
		DB:                  db,
		LoanTableName:       loanTableName,
		InvestmentTableName: investmentTableName,
	}
}

func (lim *LoanInvestmentModel) TxInvestment(ctx context.Context, loan *types.Loan, investments *types.InvestLoanRequest) error {
	// Step 1: Begin a transaction
	tx, err := lim.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Defer transaction rollback in case of error
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			log.Errorf("error while tx rollback with err: %v", rollbackErr)
		}
	}()

	// Step 2: Update the loan table
	_, err = tx.ExecContext(ctx,
		fmt.Sprintf("UPDATE %s SET state = ?, invested_amount = ? WHERE loan_id = ?", lim.LoanTableName),
		loan.State,
		loan.InvestedAmount,
		loan.LoanID)
	if err != nil {
		return err
	}

	// Step 3: Update the investment table
	query := fmt.Sprintf("INSERT INTO %s (loan_id, investor_id, amount) VALUES (?, ?, ?)", lim.InvestmentTableName)
	for _, investor := range investments.Investors {
		_, err := tx.Exec(query, investments.LoanID, investor.InvestorID, investor.Amount)
		if err != nil {
			return fmt.Errorf("failed to update investment: %v", err)
		}
	}

	// Step 4: Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
