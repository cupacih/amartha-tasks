package logic

import (
	"context"
	"errors"
	"fmt"
	"goloanservice/api/internal/svc"
	types "goloanservice/api/internal/type"
	"goloanservice/util"
	"net/url"
)

type DisbursementLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDisbursementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DisbursementLogic {
	return &DisbursementLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (dl *DisbursementLogic) Disburse(req *types.DisbursementRequest) (resp *types.DisbursementResponse, err error) {
	if err := validateDisbursementRequest(req); err != nil {
		return nil, err
	}

	loan, err := dl.svcCtx.LoanModel.FindOne(dl.ctx, req.LoanID)
	if err != nil {
		return nil, err
	}

	if req.DisbursementDate.Before(loan.CreatedAt) {
		return nil, errors.New("disbursement date must be later than the loan creation date")
	}

	// Update the loan state to "disbursed" in the database
	if err = loan.State.Disburse(); err != nil {
		return nil, fmt.Errorf("cannot disburse with err: %v", err)
	}

	// Update the database
	if err := dl.svcCtx.LoanDisbursementModel.TxDisbursement(dl.ctx, loan, req); err != nil {
		return nil, err
	}

	return &types.DisbursementResponse{
		DisbursementRequest: *req,
		IsDisbursed:         true,
	}, nil
}

func validateDisbursementRequest(req *types.DisbursementRequest) error {
	if !util.ValidateDate(req.DisbursementDate) {
		return errors.New("invalid date")
	}

	if _, err := url.ParseRequestURI(req.SignedAgreementLetter); err != nil {
		return fmt.Errorf("invalid signed agreement letter URL with err: %v", err)
	}

	if len(req.FieldOfficerID) == 0 {
		return fmt.Errorf("FieldOfficerID is empty")
	}

	return nil
}
