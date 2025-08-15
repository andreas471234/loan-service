package main

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LoanStatus represents the possible states of a loan
type LoanStatus string

const (
	StatusProposed  LoanStatus = "proposed"
	StatusApproved  LoanStatus = "approved"
	StatusInvested  LoanStatus = "invested"
	StatusDisbursed LoanStatus = "disbursed"
)

// StateTransition represents a valid state transition
type StateTransition struct {
	From   LoanStatus
	To     LoanStatus
	Action string
}

// FSM defines the finite state machine for loan states
type FSM struct {
	CurrentState LoanStatus
	Transitions  []StateTransition
}

// NewFSM creates a new FSM instance
func NewFSM() *FSM {
	return &FSM{
		CurrentState: StatusProposed,
		Transitions: []StateTransition{
			{From: StatusProposed, To: StatusApproved, Action: "approve"},
			{From: StatusApproved, To: StatusInvested, Action: "invest"},
			{From: StatusInvested, To: StatusDisbursed, Action: "disburse"},
		},
	}
}

// CanTransition checks if a transition is valid
func (fsm *FSM) CanTransition(to LoanStatus) bool {
	for _, transition := range fsm.Transitions {
		if transition.From == fsm.CurrentState && transition.To == to {
			return true
		}
	}
	return false
}

// Transition performs a state transition
func (fsm *FSM) Transition(to LoanStatus) error {
	if !fsm.CanTransition(to) {
		return errors.New("invalid state transition")
	}
	fsm.CurrentState = to
	return nil
}

// GetCurrentState returns the current state
func (fsm *FSM) GetCurrentState() LoanStatus {
	return fsm.CurrentState
}

// SetCurrentState sets the current state (used when loading from database)
func (fsm *FSM) SetCurrentState(state LoanStatus) {
	fsm.CurrentState = state
}

// GetValidTransitions returns all valid transitions from current state
func (fsm *FSM) GetValidTransitions() []StateTransition {
	var valid []StateTransition
	for _, transition := range fsm.Transitions {
		if transition.From == fsm.CurrentState {
			valid = append(valid, transition)
		}
	}
	return valid
}

// ApprovalDetails contains information required for loan approval
type ApprovalDetails struct {
	FieldValidatorProof string    `json:"field_validator_proof"` // Picture proof of field validator visit
	FieldValidatorID    string    `json:"field_validator_id"`    // Employee ID of field validator
	ApprovalDate        time.Time `json:"approval_date"`         // Date of approval
}

// Investment represents an individual investment in a loan
type Investment struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	LoanID    string    `json:"loan_id" gorm:"not null"`
	InvestorID string   `json:"investor_id" gorm:"not null"`
	Amount    float64   `json:"amount" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DisbursementDetails contains information required for loan disbursement
type DisbursementDetails struct {
	SignedAgreementLink string    `json:"signed_agreement_link"` // Signed agreement letter (PDF/JPEG)
	FieldOfficerID      string    `json:"field_officer_id"`      // Employee ID of field officer
	DisbursementDate    time.Time `json:"disbursement_date"`     // Date of disbursement
}

// Loan represents a loan entity
type Loan struct {
	ID                string               `json:"id" gorm:"primaryKey;type:varchar(36)"`
	BorrowerID        string               `json:"borrower_id" gorm:"not null"`
	PrincipalAmount   float64              `json:"principal_amount" gorm:"not null"`
	Rate              float64              `json:"rate" gorm:"not null"` // Interest rate that borrower will pay
	ROI               float64              `json:"roi" gorm:"not null"`  // Return on investment for investors
	AgreementLetterLink string             `json:"agreement_letter_link"`
	Status            LoanStatus           `json:"status" gorm:"not null;default:'proposed'"`

	// Approval details (required when status is approved or beyond)
	ApprovalDetails   *ApprovalDetails     `json:"approval_details" gorm:"embedded"`

	// Investment tracking
	Investments       []Investment         `json:"investments" gorm:"foreignKey:LoanID"`
	TotalInvested     float64              `json:"total_invested" gorm:"default:0"`

	// Disbursement details (required when status is disbursed)
	DisbursementDetails *DisbursementDetails `json:"disbursement_details" gorm:"embedded"`

	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
	DeletedAt         gorm.DeletedAt       `json:"deleted_at,omitempty" gorm:"index"`
}

// BeforeCreate is a GORM hook that sets the ID before creating a record
func (l *Loan) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	return nil
}

// GetFSM returns a FSM instance for this loan
func (l *Loan) GetFSM() *FSM {
	fsm := NewFSM()
	fsm.SetCurrentState(l.Status)
	return fsm
}

// CanUpdate checks if the loan can be updated
func (l *Loan) CanUpdate() bool {
	return l.Status == StatusProposed
}

// CanDelete checks if the loan can be deleted
func (l *Loan) CanDelete() bool {
	return l.Status == StatusProposed
}

// CanApprove checks if the loan can be approved
func (l *Loan) CanApprove() bool {
	return l.Status == StatusProposed
}

// CanInvest checks if the loan can receive investments
func (l *Loan) CanInvest() bool {
	return l.Status == StatusApproved
}

// CanDisburse checks if the loan can be disbursed
func (l *Loan) CanDisburse() bool {
	return l.Status == StatusInvested && l.TotalInvested >= l.PrincipalAmount
}

// AddInvestment adds an investment to the loan
func (l *Loan) AddInvestment(investorID string, amount float64) error {
	if !l.CanInvest() {
		return errors.New("loan is not in approved status")
	}

	if l.TotalInvested+amount > l.PrincipalAmount {
		return errors.New("total investment amount would exceed loan principal")
	}

	investment := Investment{
		ID:         uuid.New().String(),
		LoanID:     l.ID,
		InvestorID: investorID,
		Amount:     amount,
	}

	l.Investments = append(l.Investments, investment)
	l.TotalInvested += amount

	// If total invested equals principal amount, automatically transition to invested
	if l.TotalInvested >= l.PrincipalAmount {
		l.Status = StatusInvested
	}

	return nil
}

// CreateLoanRequest represents the request body for creating a loan
type CreateLoanRequest struct {
	BorrowerID        string  `json:"borrower_id" binding:"required"`
	PrincipalAmount   float64 `json:"principal_amount" binding:"required,gt=0"`
	Rate              float64 `json:"rate" binding:"required,gt=0"`
	ROI               float64 `json:"roi" binding:"required,gt=0"`
	AgreementLetterLink string `json:"agreement_letter_link"`
}

// UpdateLoanRequest represents the request body for updating a loan
type UpdateLoanRequest struct {
	PrincipalAmount   *float64 `json:"principal_amount"`
	Rate              *float64 `json:"rate"`
	ROI               *float64 `json:"roi"`
	AgreementLetterLink *string `json:"agreement_letter_link"`
}

// ApproveLoanRequest represents the request body for approving a loan
type ApproveLoanRequest struct {
	FieldValidatorProof string `json:"field_validator_proof" binding:"required"`
	FieldValidatorID    string `json:"field_validator_id" binding:"required"`
}

// InvestLoanRequest represents the request body for investing in a loan
type InvestLoanRequest struct {
	InvestorID string  `json:"investor_id" binding:"required"`
	Amount     float64 `json:"amount" binding:"required,gt=0"`
}

// DisburseLoanRequest represents the request body for disbursing a loan
type DisburseLoanRequest struct {
	SignedAgreementLink string `json:"signed_agreement_link" binding:"required"`
	FieldOfficerID      string `json:"field_officer_id" binding:"required"`
}

// LoanResponse represents the response body for loan operations
type LoanResponse struct {
	ID                string               `json:"id"`
	BorrowerID        string               `json:"borrower_id"`
	PrincipalAmount   float64              `json:"principal_amount"`
	Rate              float64              `json:"rate"`
	ROI               float64              `json:"roi"`
	AgreementLetterLink string             `json:"agreement_letter_link"`
	Status            LoanStatus           `json:"status"`
	ApprovalDetails   *ApprovalDetails     `json:"approval_details,omitempty"`
	Investments       []Investment         `json:"investments,omitempty"`
	TotalInvested     float64              `json:"total_invested"`
	DisbursementDetails *DisbursementDetails `json:"disbursement_details,omitempty"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
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