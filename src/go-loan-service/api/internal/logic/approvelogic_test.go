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

// MockLoanModelForApproval is a mock implementation of the LoanModel interface
type MockLoanModelForApproval struct {
	mock.Mock
}

// FindOne mocks the FindOne method of LoanModel
func (m *MockLoanModelForApproval) FindOne(ctx context.Context, loanID string) (*types.Loan, error) {
	args := m.Called(ctx, loanID)
	return args.Get(0).(*types.Loan), args.Error(1)
}

func (m *MockLoanModelForApproval) InsertOne(ctx context.Context, loan *types.Loan) error {
	return nil
}

// MockLoanApprovalModel is a mock implementation of the LoanApprovalModel interface
type MockLoanApprovalModel struct {
	mock.Mock
}

// TxApproval mocks the TxApproval method of LoanApprovalModel
func (m *MockLoanApprovalModel) TxApproval(ctx context.Context, loan *types.Loan, approvalInfo *types.ApprovalInfo) error {
	args := m.Called(ctx, loan, approvalInfo)
	return args.Error(0)
}

func TestApproveLogic_Approve(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Create a mock loan model
	mockLoanModel := &MockLoanModelForApproval{}

	// Mock the FindOne method of the loan model
	mockLoanModel.On("FindOne", ctx, "test_loan_id").Return(&types.Loan{
		LoanID:    "test_loan_id",
		CreatedAt: time.Now().Add(-24 * time.Hour), // Assume loan created 24 hours ago
		State:     "proposed",
	}, nil)

	// Create a mock loan approval model
	mockLoanApprovalModel := &MockLoanApprovalModel{}

	// Set up expectation for the TxApproval method
	mockLoanApprovalModel.On("TxApproval", ctx, mock.Anything, mock.Anything).Return(nil)

	// Create an instance of ApproveLogic with the mock loan model and mock loan approval model
	approveLogic := logic.NewApproveLogic(ctx, &svc.ServiceContext{
		LoanModel:         mockLoanModel,
		LoanApprovalModel: mockLoanApprovalModel,
	})

	// Create an instance of ApprovalInfoRequest with test data
	request := &types.ApprovalInfoRequest{
		LoanID:           "test_loan_id",
		PictureProofURL:  "https://example.com/picture.jpg",
		FieldValidatorID: "test_field_validator_id",
		ApprovalDate:     time.Now(),
	}

	// Call the Approve method
	response, err := approveLogic.Approve(request)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert that response is not nil
	assert.NotNil(t, response)

	// Assert that the response indicates approval
	assert.True(t, response.IsApproved)
}
