package model

import (
	"context"
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	types "goloanservice/api/internal/type"
)

type ILoanDisbursementModel interface {
	TxDisbursement(ctx context.Context, loan *types.Loan, disbursement *types.DisbursementRequest) error
}

type LoanDisbursementModel struct {
	DB                    *sql.DB
	LoanTableName         string
	DisbursementTableName string
}

func NewLoanDisbursementModel(db *sql.DB, loanTableName, disbursementTableName string) ILoanDisbursementModel {
	return &LoanDisbursementModel{
		DB:                    db,
		LoanTableName:         loanTableName,
		DisbursementTableName: disbursementTableName,
	}
}

// TxDisbursement transact disbursement to table loan and disbursement
func (ldm *LoanDisbursementModel) TxDisbursement(ctx context.Context, loan *types.Loan, disbursement *types.DisbursementRequest) error {
	// Begin transaction
	tx, err := ldm.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Defer transaction rollback in case of error
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			log.Errorf("error while tx rollback with err: %v", rollbackErr)
		}
	}()

	// Update Loan table
	_, err = tx.ExecContext(ctx,
		fmt.Sprintf("UPDATE %s SET state = ? WHERE loan_id = ?", ldm.LoanTableName),
		loan.State, loan.LoanID)
	if err != nil {
		return err
	}

	// Update disbursement table
	_, err = tx.ExecContext(ctx,
		fmt.Sprintf("INSERT INTO %s (loan_id, signed_agreement_letter, field_officer_id, disbursement_date) VALUES (?, ?, ?, ?)", ldm.DisbursementTableName),
		disbursement.LoanID,
		disbursement.SignedAgreementLetter,
		disbursement.FieldOfficerID,
		disbursement.DisbursementDate)
	if err != nil {
		return err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil

}
