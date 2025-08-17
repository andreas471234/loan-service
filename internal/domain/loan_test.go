package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoanCanUpdate(t *testing.T) {
	loan := &Loan{Status: StatusProposed}
	assert.True(t, loan.CanUpdate())

	loan.Status = StatusApproved
	assert.False(t, loan.CanUpdate())
}

func TestLoanCanDelete(t *testing.T) {
	loan := &Loan{Status: StatusProposed}
	assert.True(t, loan.CanDelete())

	loan.Status = StatusApproved
	assert.False(t, loan.CanDelete())
}

func TestLoanCanApprove(t *testing.T) {
	loan := &Loan{Status: StatusProposed}
	assert.True(t, loan.CanApprove())

	loan.Status = StatusApproved
	assert.False(t, loan.CanApprove())
}

func TestLoanCanInvest(t *testing.T) {
	loan := &Loan{Status: StatusApproved}
	assert.True(t, loan.CanInvest())

	loan.Status = StatusProposed
	assert.False(t, loan.CanInvest())
}

func TestLoanCanDisburse(t *testing.T) {
	loan := &Loan{
		Status:          StatusInvested,
		TotalInvested:   25000.00,
		PrincipalAmount: 25000.00,
	}
	assert.True(t, loan.CanDisburse())

	loan.TotalInvested = 20000.00
	assert.False(t, loan.CanDisburse())
}

func TestLoanAddInvestment(t *testing.T) {
	loan := &Loan{
		Status:          StatusApproved,
		PrincipalAmount: 25000.00,
		TotalInvested:   0.0,
	}

	err := loan.AddInvestment("investor_001", 10000.00)
	assert.NoError(t, err)
	assert.Equal(t, 10000.00, loan.TotalInvested)
	assert.Len(t, loan.Investments, 1)

	// Add more investment to reach full amount
	err = loan.AddInvestment("investor_002", 15000.00)
	assert.NoError(t, err)
	assert.Equal(t, 25000.00, loan.TotalInvested)
	assert.Equal(t, StatusInvested, loan.Status)
}

func TestLoanAddInvestmentExceedsLimit(t *testing.T) {
	loan := &Loan{
		Status:          StatusApproved,
		PrincipalAmount: 25000.00,
		TotalInvested:   0.0,
	}

	err := loan.AddInvestment("investor_001", 30000.00)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "total investment amount would exceed loan principal")
}

func TestLoanAddInvestmentInvalidStatus(t *testing.T) {
	loan := &Loan{
		Status:          StatusProposed,
		PrincipalAmount: 25000.00,
		TotalInvested:   0.0,
	}

	err := loan.AddInvestment("investor_001", 10000.00)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "loan is not in approved status")
}
