package handler

import (
	"github.com/gorilla/mux"
	"goloanservice/api/internal/svc"
)

func CreateRouter(serverCtx *svc.ServiceContext) *mux.Router {
	r := mux.NewRouter()

	r.Handle("/loans/{loan_id}", GetLoanHandler(serverCtx)).Methods("GET")
	r.Handle("/loan", LoanHandler(serverCtx)).Methods("POST")
	r.Handle("/loan/approve", ApprovalHandler(serverCtx)).Methods("POST")
	r.Handle("/loan/invest", InvestHandler(serverCtx)).Methods("POST")
	r.Handle("/loan/disburse", DisburseHandler(serverCtx)).Methods("POST")

	return r
}
