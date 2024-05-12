package svc

import (
	"database/sql"
	"goloanservice/api/model"
)

type ServiceContext struct {
	DB                    *sql.DB
	LoanModel             model.ILoanModel
	LoanApprovalModel     model.ILoanApprovalModel
	LoanInvestmentModel   model.ILoanInvestmentModel
	LoanDisbursementModel model.ILoanDisbursementModel
}

func NewServiceContext(db *sql.DB,
	loanModel model.ILoanModel,
	loanApprovalModel model.ILoanApprovalModel,
	loanInvestmentModel model.ILoanInvestmentModel,
	loanDisbursementModel model.ILoanDisbursementModel) *ServiceContext {

	return &ServiceContext{
		DB:                    db,
		LoanModel:             loanModel,
		LoanApprovalModel:     loanApprovalModel,
		LoanInvestmentModel:   loanInvestmentModel,
		LoanDisbursementModel: loanDisbursementModel,
	}
}
