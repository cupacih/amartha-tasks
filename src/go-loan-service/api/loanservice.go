package main

import (
	"database/sql"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zeromicro/go-zero/core/conf"
	"goloanservice/api/internal/config"
	"goloanservice/api/internal/handler"
	"goloanservice/api/internal/svc"
	"goloanservice/api/model"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var (
	configFile = flag.String("f", "etc/loanservice.yaml", "the config file")
)

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	log.Printf("config: %v", c)
	db, err := sql.Open("mysql", fmt.Sprintf("%s?parseTime=true", c.Datasource))
	if err != nil {
		panic(err)
	}

	loanModel := model.NewLoanModel(db, c.LoansTableName)
	loanApprovalModel := model.NewLoanApprovalModel(db, c.LoansTableName, c.ApprovalInfoTableName)
	loanInvestmentModel := model.NewLoanInvestmentModel(db, c.LoansTableName, c.InvestmentsTableName)
	loanDisbursementModel := model.NewLoanDisbursementModel(db, c.LoansTableName, c.DisbursementTableName)

	svctx := svc.NewServiceContext(db,
		loanModel,
		loanApprovalModel,
		loanInvestmentModel,
		loanDisbursementModel)
	r := handler.CreateRouter(svctx)

	fmt.Printf("Starting server at %s:%d...\n", "localhost", c.ListenToPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", c.ListenToPort), r))
}
