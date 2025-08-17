package repository

import (
	"loan-service/internal/domain"

	"gorm.io/gorm"
)

// LoanRepository defines the interface for loan data operations
type LoanRepository interface {
	Create(loan *domain.Loan) error
	FindByID(id string) (*domain.Loan, error)
	FindAll(filters map[string]interface{}) ([]domain.Loan, error)
	Update(loan *domain.Loan) error
	Delete(id string) error
}

// loanRepository implements LoanRepository
type loanRepository struct {
	db *gorm.DB
}

// NewLoanRepository creates a new loan repository
func NewLoanRepository(db *gorm.DB) LoanRepository {
	return &loanRepository{db: db}
}

// Create creates a new loan
func (r *loanRepository) Create(loan *domain.Loan) error {
	return r.db.Create(loan).Error
}

// FindByID finds a loan by ID
func (r *loanRepository) FindByID(id string) (*domain.Loan, error) {
	var loan domain.Loan
	err := r.db.Preload("Investments").First(&loan, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &loan, nil
}

// FindAll finds all loans with optional filters
func (r *loanRepository) FindAll(filters map[string]interface{}) ([]domain.Loan, error) {
	var loans []domain.Loan
	query := r.db.Preload("Investments")

	if status, ok := filters["status"]; ok {
		query = query.Where("status = ?", status)
	}

	if borrowerID, ok := filters["borrower_id"]; ok {
		query = query.Where("borrower_id = ?", borrowerID)
	}

	err := query.Find(&loans).Error
	return loans, err
}

// Update updates a loan
func (r *loanRepository) Update(loan *domain.Loan) error {
	return r.db.Save(loan).Error
}

// Delete deletes a loan
func (r *loanRepository) Delete(id string) error {
	return r.db.Delete(&domain.Loan{}, "id = ?", id).Error
}
