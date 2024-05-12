package handler

import (
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/zeromicro/go-zero/rest/httpx"
	"goloanservice/api/internal/logic"
	"goloanservice/api/internal/svc"
	types "goloanservice/api/internal/type"
	"net/http"
)

func GetLoanHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract loan ID from URL path parameter
		// Get loan_id path parameter from the request
		vars := mux.Vars(r)
		loanID := vars["loan_id"]

		req := &types.GetLoanRequest{
			LoanID: loanID,
		}

		if err := uuid.Validate(req.LoanID); err != nil { // to check empty string, invalid format, SQL injection, etc
			httpx.ErrorCtx(r.Context(), w, err)
		}

		l := logic.NewLoanLogic(r.Context(), svcCtx)
		resp, err := l.GetLoan(req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func LoanHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoanRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewLoanLogic(r.Context(), svcCtx)
		resp, err := l.CreateLoan(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func ApprovalHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ApprovalInfoRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		if err := uuid.Validate(req.LoanID); err != nil { // to check empty string, invalid format, SQL injection, etc
			httpx.ErrorCtx(r.Context(), w, err)
		}

		l := logic.NewApproveLogic(r.Context(), svcCtx)
		resp, err := l.Approve(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func InvestHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.InvestLoanRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		if err := uuid.Validate(req.LoanID); err != nil { // to check empty string, invalid format, SQL injection, etc
			httpx.ErrorCtx(r.Context(), w, err)
		}

		l := logic.NewInvestLogic(r.Context(), svcCtx)
		resp, err := l.Invest(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func DisburseHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DisbursementRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		if err := uuid.Validate(req.LoanID); err != nil { // to check empty string, invalid format, SQL injection, etc
			httpx.ErrorCtx(r.Context(), w, err)
		}

		l := logic.NewDisbursementLogic(r.Context(), svcCtx)
		resp, err := l.Disburse(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
