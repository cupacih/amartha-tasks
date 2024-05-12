package logic

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"goloanservice/api/internal/svc"
	types "goloanservice/api/internal/type"
)

type InvestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInvestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InvestLogic {
	return &InvestLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Invest updates the database with the investment details
func (il *InvestLogic) Invest(req *types.InvestLoanRequest) (resp *types.InvestLoanResponse, err error) {
	// Step 1: Validate the request data
	if err := validateInvestmentRequest(req); err != nil {
		return nil, err
	}

	// Step 2: Calculate the total invested amount
	loan, err := il.svcCtx.LoanModel.FindOne(il.ctx, req.LoanID)
	if err != nil {
		return nil, err
	}
	totalInvested := 0.0
	for _, investor := range req.Investors {
		totalInvested += investor.Amount
	}

	// Ensure the total invested amount does not exceed the loan principal amount
	// 	ASSUMPTION:
	//		The requirement does NOT mention recurring investment.
	//		The investment happen only once. Once a loan has been invested, it can NOT have other investments.
	if totalInvested > loan.PrincipalAmount {
		return nil, fmt.Errorf("total invested amount (%f) exceeds the loan principal amount", totalInvested)
	}

	// Step 3: Check  if the loan is in the correct state for investment
	//	and update if state is valid
	if err = loan.State.Invest(); err != nil {
		return nil, fmt.Errorf("cannot invest with err: %v", err)
	}

	// Step 4: Update the loan's invested amount
	loan.InvestedAmount = totalInvested

	// Step 5: Update the database with the investment details
	err = il.svcCtx.LoanInvestmentModel.TxInvestment(il.ctx, loan, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update loan in the database: %v", err)
	}

	// Step 6: once invested all investors will receive an email containing link to agreement letter (pdf)
	sendEmail(loan, req)

	return &types.InvestLoanResponse{
		LoanID:     loan.LoanID,
		Investors:  req.Investors,
		IsInvested: true,
	}, nil
}

func sendEmail(loan *types.Loan, investments *types.InvestLoanRequest) {
	// Ideal: to integrate with SMTP mailing service or send to notification service via event streaming
	for _, investor := range investments.Investors {
		log.Infof("loanID:%s - email is sent to investorID:%s", loan.LoanID, investor.InvestorID)
	}
}

func validateInvestmentRequest(req *types.InvestLoanRequest) error {
	for idx, investor := range req.Investors {
		if len(investor.InvestorID) == 0 {
			return errors.New(fmt.Sprintf("investorID is empty at idx: %d", idx))
		}
		if investor.Amount == 0.0 {
			return errors.New(fmt.Sprintf("investor amount is zero at idx: %d", idx))
		}
	}
	return nil
}
