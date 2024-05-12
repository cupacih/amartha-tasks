package types

import (
	"goloanservice/states"
	"time"
)

// Loan represents a loan with its details.
type Loan struct {
	ID              int64             `json:"id"`
	LoanID          string            `json:"loan_id"`
	BorrowerID      string            `json:"borrower_id"`
	PrincipalAmount float64           `json:"principal_amount"`
	Rate            float64           `json:"rate"`
	ROI             float64           `json:"roi"`
	AgreementLink   string            `json:"agreement_link"`
	State           states.LoanState  `json:"state"`
	ApprovedInfo    *ApprovalInfo     `json:"approved_info,omitempty"`
	InvestedAmount  float64           `json:"invested_amount"`
	Investors       []*Investor       `json:"investors"`
	DisbursedInfo   *DisbursementInfo `json:"disbursed_info,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
}

type LoanRequest struct {
	BorrowerID      string  `json:"borrower_id"`
	PrincipalAmount float64 `json:"principal_amount"`
	Rate            float64 `json:"rate"`
}

type GetLoanRequest struct {
	LoanID string `json:"loan_id"`
}

type LoanResponse struct {
	LoanID          string    `json:"loan_id"`
	BorrowerID      string    `json:"borrower_id"`
	PrincipalAmount float64   `json:"principal_amount"`
	Rate            float64   `json:"rate"`
	AgreementLink   string    `json:"agreement_link"`
	CreatedAt       time.Time `json:"created_at"`
}

// ApprovalInfo contains the details of the loan approval.
type ApprovalInfo struct {
	ID               int64     `json:"id"`
	LoanID           string    `json:"loan_id"`
	PictureProofURL  string    `json:"picture_proof_url"`
	FieldValidatorID string    `json:"field_validator_id"`
	ApprovalDate     time.Time `json:"approval_date"`
	CreatedAt        time.Time `json:"created_at"`
}

type ApprovalInfoRequest struct {
	LoanID           string    `json:"loan_id"`
	PictureProofURL  string    `json:"picture_proof_url"`
	FieldValidatorID string    `json:"field_validator_id"`
	ApprovalDate     time.Time `json:"approval_date"`
}

type ApprovalInfoResponse struct {
	LoanID           string    `json:"loan_id"`
	PictureProofURL  string    `json:"picture_proof_url"`
	FieldValidatorID string    `json:"field_validator_id"`
	ApprovalDate     time.Time `json:"approval_date"`
	IsApproved       bool      `json:"is_approved"`
}

// Investor represents an investor in the loan.
type Investor struct {
	InvestorID string  `json:"investor_id"`
	Amount     float64 `json:"amount"`
}

// Investment represents an investment in the loan.
type Investment struct {
	ID         int64     `json:"id"`
	LoanID     string    `json:"loan_id"`
	InvestorID string    `json:"investor_id"`
	Amount     float64   `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}

type InvestLoanRequest struct {
	LoanID    string     `json:"loan_id"`
	Investors []Investor `json:"investors"`
}

type InvestLoanResponse struct {
	LoanID     string     `json:"loan_id"`
	Investors  []Investor `json:"investors"`
	IsInvested bool       `json:"is_invested"`
}

// DisbursementInfo contains the details of the loan disbursement.
type DisbursementInfo struct {
	ID                    int64     `json:"id"`
	LoanID                string    `json:"loan_id"`
	SignedAgreementLetter string    `json:"signed_agreement_letter"`
	FieldOfficerID        string    `json:"field_officer_id"`
	DisbursementDate      time.Time `json:"disbursement_date"`
	CreatedAt             time.Time `json:"created_at"`
}

type DisbursementRequest struct {
	LoanID                string    `json:"loan_id"`
	SignedAgreementLetter string    `json:"signed_agreement_letter"`
	FieldOfficerID        string    `json:"field_officer_id"`
	DisbursementDate      time.Time `json:"disbursement_date"`
}

type DisbursementResponse struct {
	DisbursementRequest
	IsDisbursed bool `json:"is_disbursed"`
}
