package logic_test

import (
	"context"
	"goloanservice/api/internal/logic"
	"goloanservice/api/internal/svc"
	"goloanservice/api/internal/type"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLoanModel is a mock implementation of the LoanModel interface
type MockLoanModelForDisburse struct {
	mock.Mock
}

// FindOne mocks the FindOne method of LoanModel
func (m *MockLoanModelForDisburse) FindOne(ctx context.Context, loanID string) (*types.Loan, error) {
	args := m.Called(ctx, loanID)
	return args.Get(0).(*types.Loan), args.Error(1)
}

func (m *MockLoanModelForDisburse) InsertOne(ctx context.Context, loan *types.Loan) error {
	return nil
}

// MockLoanDisbursementModel is a mock implementation of the LoanDisbursementModel interface
type MockLoanDisbursementModel struct {
	mock.Mock
}

// TxDisbursement mocks the TxDisbursement method of LoanDisbursementModel
func (m *MockLoanDisbursementModel) TxDisbursement(ctx context.Context, loan *types.Loan, req *types.DisbursementRequest) error {
	args := m.Called(ctx, loan, req)
	return args.Error(0)
}

func TestDisbursementLogic_Disburse(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Create a mock loan model
	mockLoanModel := &MockLoanModelForDisburse{}

	// Mock the FindOne method of the loan model to return a valid loan
	loanCreateAt := time.Now().Add(-24 * time.Hour)
	mockLoanModel.On("FindOne", ctx, "test_loan_id").Return(&types.Loan{
		LoanID:    "test_loan_id",
		CreatedAt: loanCreateAt,
		State:     "invested",
	}, nil)

	// Create a mock loan disbursement model
	mockLoanDisbursementModel := &MockLoanDisbursementModel{}

	// Set up expectation for the TxDisbursement method
	mockLoanDisbursementModel.On("TxDisbursement", ctx, mock.Anything, mock.Anything).Return(nil)

	// Create an instance of DisbursementLogic with the mock loan model and mock loan disbursement model
	disbursementLogic := logic.NewDisbursementLogic(ctx, &svc.ServiceContext{
		LoanModel:             mockLoanModel,
		LoanDisbursementModel: mockLoanDisbursementModel,
	})

	// Create a valid disbursement request
	disbursementDate := time.Now().Add(24 * time.Hour)
	request := &types.DisbursementRequest{
		LoanID:                "test_loan_id",
		DisbursementDate:      disbursementDate,
		SignedAgreementLetter: "https://example.com/agreement.pdf",
		FieldOfficerID:        "officer123",
	}

	// Call the Disburse method
	response, err := disbursementLogic.Disburse(request)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert that response is not nil
	assert.NotNil(t, response)

	// Assert that the loan's state has been updated to "disbursed"
	assert.True(t, response.IsDisbursed)

	// Assert that the disbursement date is later than the loan creation date
	assert.True(t, disbursementDate.After(loanCreateAt))

	// Assert that the loan's state has been updated to "disbursed"
	assert.True(t, response.IsDisbursed)
}

func TestDisbursementLogic_Disburse_InvalidRequest(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Create a mock loan model
	mockLoanModel := &MockLoanModelForDisburse{}

	// Mock the FindOne method of the loan model to return a valid loan
	loanCreateAt := time.Now().Add(+24 * time.Hour)
	mockLoanModel.On("FindOne", ctx, "test_loan_id").Return(&types.Loan{
		ID:              123,
		LoanID:          "test_loan_id",
		BorrowerID:      "123",
		PrincipalAmount: 1000,
		Rate:            0.05,
		ROI:             50,
		AgreementLink:   "https://example.com/agreement.pdf",
		State:           "invested",
		InvestedAmount:  1000,
		//ApprovedInfo:    nil,
		//Investors:       nil,
		//DisbursedInfo:   nil,
		CreatedAt: loanCreateAt,
	}, nil)

	// Create an instance of DisbursementLogic with the mock loan model
	disbursementLogic := logic.NewDisbursementLogic(ctx, &svc.ServiceContext{
		LoanModel: mockLoanModel,
	})

	// Create an invalid disbursement request (invalid date)
	request := &types.DisbursementRequest{
		LoanID:                "test_loan_id",
		DisbursementDate:      time.Now().Add(-24 * time.Hour), // Use a past date
		SignedAgreementLetter: "https://example.com/agreement.pdf",
		FieldOfficerID:        "officer123",
	}

	// Call the Disburse method
	response, err := disbursementLogic.Disburse(request)

	// Assert that an error occurred
	assert.Error(t, err)

	// Assert that response is nil
	assert.Nil(t, response)

	// Assert the error message
	expectedErrorMsg := "disbursement date must be later than the loan creation date"
	assert.EqualError(t, err, expectedErrorMsg)
}

func TestDisbursementLogic_Disburse_InvalidAgreementLetterURL(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Create a mock loan model
	mockLoanModel := &MockLoanModelForDisburse{}

	// Mock the FindOne method of the loan model to return a valid loan
	loanCreateAt := time.Now().Add(-24 * time.Hour)
	mockLoanModel.On("FindOne", ctx, "test_loan_id").Return(&types.Loan{
		ID:              123,
		LoanID:          "test_loan_id",
		BorrowerID:      "123",
		PrincipalAmount: 1000,
		Rate:            0.05,
		ROI:             50,
		AgreementLink:   "https://example.com/agreement.pdf",
		State:           "invested",
		InvestedAmount:  1000,
		//ApprovedInfo:    nil,
		//Investors:       nil,
		//DisbursedInfo:   nil,
		CreatedAt: loanCreateAt,
	}, nil)

	// Create an instance of DisbursementLogic with the mock loan model
	disbursementLogic := logic.NewDisbursementLogic(ctx, &svc.ServiceContext{
		LoanModel: mockLoanModel,
	})

	// Create an invalid disbursement request (invalid agreement letter URL)
	request := &types.DisbursementRequest{
		LoanID:                "test_loan_id",
		DisbursementDate:      time.Now().Add(24 * time.Hour),
		SignedAgreementLetter: "invalid-url", // Use an invalid URL
		FieldOfficerID:        "officer123",
	}

	// Call the Disburse method
	response, err := disbursementLogic.Disburse(request)

	// Assert that an error occurred
	assert.Error(t, err)

	// Assert that response is nil
	assert.Nil(t, response)

	// Assert the error message
	expectedErrorMsg := "invalid signed agreement letter URL with err: parse \"invalid-url\": invalid URI for request"
	assert.EqualError(t, err, expectedErrorMsg)
}
