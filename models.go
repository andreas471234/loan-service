package main

import (
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
	StatusRejected  LoanStatus = "rejected"
)

// Loan represents a loan entity
type Loan struct {
	ID          string     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	BorrowerID  string     `json:"borrower_id" gorm:"not null"`
	Amount      float64    `json:"amount" gorm:"not null"`
	InterestRate float64   `json:"interest_rate" gorm:"not null"`
	Term        int        `json:"term" gorm:"not null"` // in months
	Purpose     string     `json:"purpose" gorm:"not null"`
	Status      LoanStatus `json:"status" gorm:"not null;default:'proposed'"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// BeforeCreate is a GORM hook that sets the ID before creating a record
func (l *Loan) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	return nil
}

// CreateLoanRequest represents the request body for creating a loan
type CreateLoanRequest struct {
	BorrowerID  string  `json:"borrower_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	InterestRate float64 `json:"interest_rate" binding:"required,gt=0"`
	Term        int     `json:"term" binding:"required,gt=0"`
	Purpose     string  `json:"purpose" binding:"required"`
	Description string  `json:"description"`
}

// UpdateLoanRequest represents the request body for updating a loan
type UpdateLoanRequest struct {
	Amount      *float64 `json:"amount"`
	InterestRate *float64 `json:"interest_rate"`
	Term        *int     `json:"term"`
	Purpose     *string  `json:"purpose"`
	Description *string  `json:"description"`
}

// LoanResponse represents the response body for loan operations
type LoanResponse struct {
	ID          string     `json:"id"`
	BorrowerID  string     `json:"borrower_id"`
	Amount      float64    `json:"amount"`
	InterestRate float64   `json:"interest_rate"`
	Term        int        `json:"term"`
	Purpose     string     `json:"purpose"`
	Status      LoanStatus `json:"status"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
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