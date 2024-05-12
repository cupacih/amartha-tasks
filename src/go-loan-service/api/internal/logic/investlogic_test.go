package logic_test

import (
	"context"
	"goloanservice/api/internal/logic"
	"goloanservice/api/internal/svc"
	"goloanservice/api/internal/type"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLoanModel is a mock implementation of the LoanModel interface
type mockLoanModel struct {
	mock.Mock
}

// FindOne mocks the FindOne method of LoanModel
func (m *mockLoanModel) FindOne(ctx context.Context, loanID string) (*types.Loan, error) {
	args := m.Called(ctx, loanID)
	return args.Get(0).(*types.Loan), args.Error(1)
}

func (m *mockLoanModel) InsertOne(ctx context.Context, loan *types.Loan) error {
	return nil
}

// MockLoanInvestmentModel is a mock implementation of the LoanInvestmentModel interface
type MockLoanInvestmentModel struct {
	mock.Mock
}

// TxInvestment mocks the TxInvestment method of LoanInvestmentModel
func (m *MockLoanInvestmentModel) TxInvestment(ctx context.Context, loan *types.Loan, req *types.InvestLoanRequest) error {
	args := m.Called(ctx, loan, req)
	return args.Error(0)
}

func TestInvestLogic_Invest(t *testing.T) {
	// Create a context
	ctx := context.Background()

	// Create a mock loan model
	mockLoanModel := &mockLoanModel{}

	// Mock the FindOne method of the loan model to return a valid loan state
	mockLoanModel.On("FindOne", ctx, "test_loan_id").Return(&types.Loan{
		LoanID:          "test_loan_id",
		PrincipalAmount: 1000.0,
		State:           "approved", // Ensure the state is "approved"
		InvestedAmount:  0.0,
	}, nil)

	// Create a mock loan investment model
	mockLoanInvestmentModel := &MockLoanInvestmentModel{}

	// Set up expectation for the TxInvestment method
	mockLoanInvestmentModel.On("TxInvestment", ctx, mock.Anything, mock.Anything).Return(nil)

	// Create an instance of InvestLogic with the mock loan model and mock loan investment model
	investLogic := logic.NewInvestLogic(ctx, &svc.ServiceContext{
		LoanModel:           mockLoanModel,
		LoanInvestmentModel: mockLoanInvestmentModel,
	})

	// Create an instance of InvestLoanRequest with test data
	request := &types.InvestLoanRequest{
		LoanID: "test_loan_id",
		Investors: []types.Investor{
			{InvestorID: "investor1", Amount: 500.0},
			{InvestorID: "investor2", Amount: 500.0},
		},
	}

	// Call the Invest method
	response, err := investLogic.Invest(request)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert that response is not nil
	assert.NotNil(t, response)

	// Assert that the loan's invested amount has been updated
	assert.Equal(t, 1000.0, response.Investors[0].Amount+response.Investors[1].Amount)
}
