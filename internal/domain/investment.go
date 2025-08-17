package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Investment represents an individual investment in a loan
type Investment struct {
	ID         string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	LoanID     string    `json:"loan_id" gorm:"not null"`
	InvestorID string    `json:"investor_id" gorm:"not null"`
	Amount     float64   `json:"amount" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// BeforeCreate is a GORM hook that sets the ID before creating a record
func (i *Investment) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	return nil
}
