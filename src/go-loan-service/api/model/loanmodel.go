package model

import (
	"context"
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	types "goloanservice/api/internal/type"
)

type ILoanModel interface {
	FindOne(ctx context.Context, loanId string) (*types.Loan, error)
	InsertOne(ctx context.Context, loan *types.Loan) error
}

type LoanModel struct {
	DB        *sql.DB
	TableName string
}

func NewLoanModel(db *sql.DB, tableName string) ILoanModel {
	return &LoanModel{
		DB:        db,
		TableName: tableName,
	}
}

func (lm *LoanModel) FindOne(ctx context.Context, loanId string) (*types.Loan, error) {
	var loan types.Loan
	err := lm.DB.QueryRowContext(ctx,
		fmt.Sprintf(
			"SELECT id, loan_id, borrower_id, principal_amount, rate, roi, agreement_link, state, invested_amount, created_at FROM %s WHERE loan_id = ?",
			lm.TableName), loanId).
		Scan(&loan.ID,
			&loan.LoanID,
			&loan.BorrowerID,
			&loan.PrincipalAmount,
			&loan.Rate,
			&loan.ROI,
			&loan.AgreementLink,
			&loan.State,
			&loan.InvestedAmount,
			&loan.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &loan, nil
}

func (lm *LoanModel) InsertOne(ctx context.Context, loan *types.Loan) error {
	stmt, err := lm.DB.PrepareContext(ctx,
		fmt.Sprintf("INSERT INTO %s (loan_id, borrower_id, principal_amount, rate, roi, agreement_link, state) VALUES (?, ?, ?, ?, ?, ?, ?)", lm.TableName))
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Errorln(err)
		}
	}()
	_, err = stmt.Exec(loan.LoanID, loan.BorrowerID, loan.PrincipalAmount, loan.Rate, loan.ROI, loan.AgreementLink, loan.State)
	if err != nil {
		return err
	}
	return nil
}
