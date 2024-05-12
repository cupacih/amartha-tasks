package logic

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"goloanservice/api/internal/svc"
	types "goloanservice/api/internal/type"
	"goloanservice/states"
)

type LoanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoanLogic {
	return &LoanLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoanLogic) GetLoan(req *types.GetLoanRequest) (resp *types.LoanResponse, err error) {
	loan, err := l.svcCtx.LoanModel.FindOne(l.ctx, req.LoanID)
	if err != nil {
		return nil, err
	}
	return &types.LoanResponse{
		LoanID:          loan.LoanID,
		BorrowerID:      loan.BorrowerID,
		PrincipalAmount: loan.PrincipalAmount,
		Rate:            loan.Rate,
		AgreementLink:   loan.AgreementLink,
		CreatedAt:       loan.CreatedAt,
	}, nil
}

func (l *LoanLogic) CreateLoan(req *types.LoanRequest) (resp *types.LoanResponse, err error) {
	if req.Rate >= 1.0 {
		return nil, fmt.Errorf("the rate must be below 1.0, inclusive")
	}

	loanID := uuid.New().String()
	createdLoan := &types.Loan{
		LoanID:          loanID,
		BorrowerID:      req.BorrowerID,
		PrincipalAmount: req.PrincipalAmount,
		Rate:            req.Rate,
		ROI:             req.PrincipalAmount * req.Rate,
		AgreementLink:   fmt.Sprintf("http://s3.amazonaws.com/loans/%v-agreement.pdf", loanID), // mock agreement letter URL
		State:           states.InitialState(),
		InvestedAmount:  0,
	}

	err = l.svcCtx.LoanModel.InsertOne(l.ctx, createdLoan)

	return &types.LoanResponse{
		LoanID:          createdLoan.LoanID,
		BorrowerID:      createdLoan.BorrowerID,
		PrincipalAmount: createdLoan.PrincipalAmount,
		Rate:            createdLoan.Rate,
		AgreementLink:   createdLoan.AgreementLink,
		CreatedAt:       createdLoan.CreatedAt,
	}, nil
}
