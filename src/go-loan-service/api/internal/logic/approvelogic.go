package logic

import (
	"context"
	"errors"
	"goloanservice/api/internal/svc"
	types "goloanservice/api/internal/type"
	"goloanservice/util"
	"net/url"
)

type ApproveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApproveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveLogic {
	return &ApproveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (al *ApproveLogic) Approve(req *types.ApprovalInfoRequest) (resp *types.ApprovalInfoResponse, err error) {
	if len(req.FieldValidatorID) <= 0 {
		return nil, errors.New("loan must been validated with a valid employee ID of field validator")
	}

	if !util.ValidateDate(req.ApprovalDate) {
		return nil, errors.New("invalid date")
	}

	if _, err = url.ParseRequestURI(req.PictureProofURL); err != nil {
		return nil, errors.Join(errors.New("invalid picture proof"), err)
	}

	loan, err := al.svcCtx.LoanModel.FindOne(al.ctx, req.LoanID)
	if err != nil {
		return nil, err
	}

	if req.ApprovalDate.Before(loan.CreatedAt) {
		return nil, errors.New("approval date must be later than loan creation date")
	}

	if err = loan.State.Approve(); err != nil {
		return nil, err
	}

	approvalInfo := &types.ApprovalInfo{
		LoanID:           req.LoanID,
		PictureProofURL:  req.PictureProofURL,
		FieldValidatorID: req.FieldValidatorID,
		ApprovalDate:     req.ApprovalDate,
	}

	err = al.svcCtx.LoanApprovalModel.TxApproval(al.ctx, loan, approvalInfo)
	if err != nil {
		return nil, err
	}

	return &types.ApprovalInfoResponse{
		LoanID:           approvalInfo.LoanID,
		PictureProofURL:  approvalInfo.PictureProofURL,
		FieldValidatorID: approvalInfo.FieldValidatorID,
		ApprovalDate:     approvalInfo.ApprovalDate,
		IsApproved:       true,
	}, nil
}
