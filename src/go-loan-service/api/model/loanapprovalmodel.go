package model

import (
	"context"
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	types "goloanservice/api/internal/type"
)

type ILoanApprovalModel interface {
	TxApproval(ctx context.Context, loan *types.Loan, approvalInfo *types.ApprovalInfo) error
}

type LoanApprovalModel struct {
	DB                    *sql.DB
	LoanTableName         string
	ApprovalInfoTableName string
}

func NewLoanApprovalModel(db *sql.DB, loanTableName, approvalInfoTableName string) ILoanApprovalModel {
	return &LoanApprovalModel{
		DB:                    db,
		LoanTableName:         loanTableName,
		ApprovalInfoTableName: approvalInfoTableName,
	}
}

func (lam *LoanApprovalModel) TxApproval(ctx context.Context, loan *types.Loan, approvalInfo *types.ApprovalInfo) error {
	// Begin transaction
	tx, err := lam.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Update Loan table
	_, err = tx.ExecContext(ctx,
		fmt.Sprintf("UPDATE %s SET state = ? WHERE loan_id = ?", lam.LoanTableName),
		loan.State, loan.LoanID)
	if err != nil {
		errRollback := tx.Rollback()
		log.Errorln(errRollback)
		return err
	}

	// Update ApprovalInfo table
	_, err = tx.ExecContext(ctx,
		fmt.Sprintf("INSERT INTO %s (loan_id, picture_proof_url, field_validator_id, approval_date) VALUES (?, ?, ?, ?)", lam.ApprovalInfoTableName),
		approvalInfo.LoanID,
		approvalInfo.PictureProofURL,
		approvalInfo.FieldValidatorID,
		approvalInfo.ApprovalDate)
	if err != nil {
		errRollback := tx.Rollback()
		log.Errorln(errRollback)
		return err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
