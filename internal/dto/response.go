package dto

import (
	"time"

	"loan-service/internal/domain"
)

// LoanResponse represents the response body for loan operations
type LoanResponse struct {
	ID                  string                      `json:"id"`
	BorrowerID          string                      `json:"borrower_id"`
	PrincipalAmount     float64                     `json:"principal_amount"`
	Rate                float64                     `json:"rate"`
	ROI                 float64                     `json:"roi"`
	AgreementLetterLink string                      `json:"agreement_letter_link"`
	Status              domain.LoanStatus           `json:"status"`
	ApprovalDetails     *domain.ApprovalDetails     `json:"approval_details,omitempty"`
	Investments         []domain.Investment         `json:"investments,omitempty"`
	TotalInvested       float64                     `json:"total_invested"`
	DisbursementDetails *domain.DisbursementDetails `json:"disbursement_details,omitempty"`
	CreatedAt           time.Time                   `json:"created_at"`
	UpdatedAt           time.Time                   `json:"updated_at"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// TransitionResponse represents a transition response
type TransitionResponse struct {
	CurrentState domain.LoanStatus        `json:"current_state"`
	Transitions  []domain.StateTransition `json:"transitions"`
}

// ToLoanResponse converts a domain.Loan to LoanResponse
func ToLoanResponse(loan domain.Loan) LoanResponse {
	return LoanResponse{
		ID:                  loan.ID,
		BorrowerID:          loan.BorrowerID,
		PrincipalAmount:     loan.PrincipalAmount,
		Rate:                loan.Rate,
		ROI:                 loan.ROI,
		AgreementLetterLink: loan.AgreementLetterLink,
		Status:              loan.Status,
		ApprovalDetails:     loan.ApprovalDetails,
		Investments:         loan.Investments,
		TotalInvested:       loan.TotalInvested,
		DisbursementDetails: loan.DisbursementDetails,
		CreatedAt:           loan.CreatedAt,
		UpdatedAt:           loan.UpdatedAt,
	}
}
