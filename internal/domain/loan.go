package domain

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

// Loan represents a loan entity
type Loan struct {
	ID                  string               `json:"id" gorm:"primaryKey;type:varchar(36)"`
	BorrowerID          string               `json:"borrower_id" gorm:"not null"`
	PrincipalAmount     float64              `json:"principal_amount" gorm:"not null"`
	Rate                float64              `json:"rate" gorm:"not null"`
	ROI                 float64              `json:"roi" gorm:"not null"`
	AgreementLetterLink string               `json:"agreement_letter_link"`
	Status              LoanStatus           `json:"status" gorm:"not null;default:'proposed'"`
	ApprovalDetails     *ApprovalDetails     `json:"approval_details" gorm:"embedded"`
	Investments         []Investment         `json:"investments" gorm:"foreignKey:LoanID"`
	TotalInvested       float64              `json:"total_invested" gorm:"default:0"`
	DisbursementDetails *DisbursementDetails `json:"disbursement_details" gorm:"embedded"`
	CreatedAt           time.Time            `json:"created_at"`
	UpdatedAt           time.Time            `json:"updated_at"`
	DeletedAt           gorm.DeletedAt       `json:"deleted_at,omitempty" gorm:"index"`
}

// ApprovalDetails contains information required for loan approval
type ApprovalDetails struct {
	FieldValidatorProof string    `json:"field_validator_proof"`
	FieldValidatorID    string    `json:"field_validator_id"`
	ApprovalDate        time.Time `json:"approval_date"`
}

// DisbursementDetails contains information required for loan disbursement
type DisbursementDetails struct {
	SignedAgreementLink string    `json:"signed_agreement_link"`
	FieldOfficerID      string    `json:"field_officer_id"`
	DisbursementDate    time.Time `json:"disbursement_date"`
}

// BeforeCreate is a GORM hook that sets the ID before creating a record
func (l *Loan) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	return nil
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
