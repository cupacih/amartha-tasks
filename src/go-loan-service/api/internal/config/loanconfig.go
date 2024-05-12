package config

type Config struct {
	Datasource            string
	ListenToPort          int
	LoansTableName        string
	ApprovalInfoTableName string
	InvestmentsTableName  string
	DisbursementTableName string
}
