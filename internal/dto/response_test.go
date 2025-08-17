package dto

import (
	"testing"
	"time"

	"loan-service/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestToLoanResponse(t *testing.T) {
	loan := domain.Loan{
		ID:                  "test-id",
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
		Status:              domain.StatusProposed,
		TotalInvested:       0.0,
		AgreementLetterLink: "https://example.com/agreement.pdf",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	
	response := ToLoanResponse(loan)
	
	assert.Equal(t, loan.ID, response.ID)
	assert.Equal(t, loan.BorrowerID, response.BorrowerID)
	assert.Equal(t, loan.PrincipalAmount, response.PrincipalAmount)
	assert.Equal(t, loan.Rate, response.Rate)
	assert.Equal(t, loan.ROI, response.ROI)
	assert.Equal(t, loan.Status, response.Status)
	assert.Equal(t, loan.TotalInvested, response.TotalInvested)
	assert.Equal(t, loan.AgreementLetterLink, response.AgreementLetterLink)
}
