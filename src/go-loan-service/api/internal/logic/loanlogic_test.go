package logic_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"goloanservice/api/internal/logic"
	"goloanservice/api/internal/svc"
	types "goloanservice/api/internal/type"
	"goloanservice/states"
	"testing"
	"time"
)

// MockLoanModel is a mock implementation of the LoanModel interface for testing purposes.
type MockLoanModel struct {
	mock.Mock
}

func (m *MockLoanModel) FindOne(ctx context.Context, loanID string) (*types.Loan, error) {
	if loanID == "valid_loan_id" {
		return &types.Loan{
			LoanID:          "valid_loan_id",
			BorrowerID:      "borrower_id",
			PrincipalAmount: 1000,
			Rate:            0.05,
			AgreementLink:   "http://example.com/agreement.pdf",
			CreatedAt:       time.Now(),
			State:           states.InitialState(),
		}, nil
	}
	return nil, errors.New("loan not found")
}

func (m *MockLoanModel) InsertOne(ctx context.Context, loan *types.Loan) error {
	return nil
}

func TestGetLoan(t *testing.T) {
	// Create a mock ServiceContext with a mock LoanModel
	mockLoanModel := &MockLoanModel{}
	svcCtx := &svc.ServiceContext{
		LoanModel: mockLoanModel,
	}

	// Create a new instance of the logic with the mock ServiceContext
	l := logic.NewLoanLogic(context.Background(), svcCtx)

	// Test case: Successful retrieval of loan
	req := &types.GetLoanRequest{LoanID: "valid_loan_id"}
	resp, err := l.GetLoan(req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := &types.LoanResponse{
		LoanID:          "valid_loan_id",
		BorrowerID:      "borrower_id",
		PrincipalAmount: 1000,
		Rate:            0.05,
		AgreementLink:   "http://example.com/agreement.pdf",
		CreatedAt:       resp.CreatedAt, // Assuming CreatedAt is the only dynamic field
	}
	if resp.LoanID != expected.LoanID || resp.BorrowerID != expected.BorrowerID || resp.PrincipalAmount != expected.PrincipalAmount ||
		resp.Rate != expected.Rate || resp.AgreementLink != expected.AgreementLink || !resp.CreatedAt.Equal(expected.CreatedAt) {
		t.Errorf("Unexpected response. Got: %+v, Expected: %+v", resp, expected)
	}

	// Test case: Loan not found
	req = &types.GetLoanRequest{LoanID: "non_existent_loan_id"}
	_, err = l.GetLoan(req)
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}
}

func TestCreateLoan(t *testing.T) {
	// Create a mock ServiceContext with a mock LoanModel
	mockLoanModel := &MockLoanModel{}
	svcCtx := &svc.ServiceContext{
		LoanModel: mockLoanModel,
	}

	// Create a new instance of the logic with the mock ServiceContext
	l := logic.NewLoanLogic(context.Background(), svcCtx)

	// Test case: Successful creation of loan
	req := &types.LoanRequest{
		BorrowerID:      "borrower_id",
		PrincipalAmount: 1000,
		Rate:            0.05,
	}
	resp, err := l.CreateLoan(req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if resp.BorrowerID != req.BorrowerID || resp.PrincipalAmount != req.PrincipalAmount ||
		resp.Rate != req.Rate || resp.AgreementLink == "" {
		t.Errorf("Unexpected response: %+v from request: %+v", resp, req)
	}

	// Test case: Rate above 1.0
	req = &types.LoanRequest{
		BorrowerID:      "borrower_id",
		PrincipalAmount: 1000,
		Rate:            1.1,
	}
	_, err = l.CreateLoan(req)
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}
}
