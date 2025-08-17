package dto

// CreateLoanRequest represents the request body for creating a loan
type CreateLoanRequest struct {
	BorrowerID      string  `json:"borrower_id" binding:"required"`
	PrincipalAmount float64 `json:"principal_amount" binding:"required,gt=0"`
	Rate            float64 `json:"rate" binding:"required,gt=0"`
	ROI             float64 `json:"roi" binding:"required,gt=0"`
}

// UpdateLoanRequest represents the request body for updating a loan
type UpdateLoanRequest struct {
	PrincipalAmount     *float64 `json:"principal_amount"`
	Rate                *float64 `json:"rate"`
	ROI                 *float64 `json:"roi"`
	AgreementLetterLink *string  `json:"agreement_letter_link"`
}

// ApproveLoanRequest represents the request body for approving a loan
type ApproveLoanRequest struct {
	FieldValidatorProof string `json:"field_validator_proof" binding:"required"`
	FieldValidatorID    string `json:"field_validator_id" binding:"required"`
}

// InvestLoanRequest represents the request body for investing in a loan
type InvestLoanRequest struct {
	InvestorID          string  `json:"investor_id" binding:"required"`
	Amount              float64 `json:"amount" binding:"required,gt=0"`
	AgreementLetterLink string  `json:"agreement_letter_link"`
}

// DisburseLoanRequest represents the request body for disbursing a loan
type DisburseLoanRequest struct {
	SignedAgreementLink string `json:"signed_agreement_link" binding:"required"`
	FieldOfficerID      string `json:"field_officer_id" binding:"required"`
}
