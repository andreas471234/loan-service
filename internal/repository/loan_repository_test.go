package repository

import (
	"testing"

	"loan-service/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestRepository() (LoanRepository, *gorm.DB) {
	// Create test database
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}

	// Auto migrate the schema
	err = database.AutoMigrate(&domain.Loan{}, &domain.Investment{})
	if err != nil {
		panic("failed to migrate test database")
	}
	
	repo := NewLoanRepository(database)
	return repo, database
}

func TestCreateLoan(t *testing.T) {
	repo, _ := setupTestRepository()
	
	loan := &domain.Loan{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
		AgreementLetterLink: "https://example.com/agreement/user123.pdf",
	}
	
	err := repo.Create(loan)
	require.NoError(t, err)
	
	assert.NotEmpty(t, loan.ID)
	assert.Equal(t, domain.StatusProposed, loan.Status)
}

func TestFindByID(t *testing.T) {
	repo, _ := setupTestRepository()
	
	// Create a loan first
	loan := &domain.Loan{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
		AgreementLetterLink: "https://example.com/agreement/user123.pdf",
	}
	
	err := repo.Create(loan)
	require.NoError(t, err)
	
	// Find the loan
	foundLoan, err := repo.FindByID(loan.ID)
	require.NoError(t, err)
	
	assert.Equal(t, loan.ID, foundLoan.ID)
	assert.Equal(t, loan.BorrowerID, foundLoan.BorrowerID)
	assert.Equal(t, loan.PrincipalAmount, foundLoan.PrincipalAmount)
}

func TestFindAll(t *testing.T) {
	repo, _ := setupTestRepository()
	
	// Create multiple loans
	loan1 := &domain.Loan{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
		AgreementLetterLink: "https://example.com/agreement/user123.pdf",
	}
	
	loan2 := &domain.Loan{
		BorrowerID:          "user456",
		PrincipalAmount:     30000.00,
		Rate:                5.0,
		ROI:                 7.0,
		AgreementLetterLink: "https://example.com/agreement/user456.pdf",
	}
	
	err := repo.Create(loan1)
	require.NoError(t, err)
	
	err = repo.Create(loan2)
	require.NoError(t, err)
	
	// Find all loans
	loans, err := repo.FindAll(nil)
	require.NoError(t, err)
	
	assert.Len(t, loans, 2)
}

func TestFindAllWithFilters(t *testing.T) {
	repo, _ := setupTestRepository()
	
	// Create loans with different statuses
	loan1 := &domain.Loan{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
		Status:              domain.StatusProposed,
		AgreementLetterLink: "https://example.com/agreement/user123.pdf",
	}
	
	loan2 := &domain.Loan{
		BorrowerID:          "user456",
		PrincipalAmount:     30000.00,
		Rate:                5.0,
		ROI:                 7.0,
		Status:              domain.StatusApproved,
		AgreementLetterLink: "https://example.com/agreement/user456.pdf",
	}
	
	err := repo.Create(loan1)
	require.NoError(t, err)
	
	err = repo.Create(loan2)
	require.NoError(t, err)
	
	// Find loans with status filter
	filters := map[string]interface{}{
		"status": domain.StatusProposed,
	}
	
	loans, err := repo.FindAll(filters)
	require.NoError(t, err)
	
	assert.Len(t, loans, 1)
	assert.Equal(t, domain.StatusProposed, loans[0].Status)
}

func TestUpdateLoan(t *testing.T) {
	repo, _ := setupTestRepository()
	
	// Create a loan
	loan := &domain.Loan{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
		AgreementLetterLink: "https://example.com/agreement/user123.pdf",
	}
	
	err := repo.Create(loan)
	require.NoError(t, err)
	
	// Update the loan
	loan.PrincipalAmount = 30000.00
	loan.Rate = 5.0
	
	err = repo.Update(loan)
	require.NoError(t, err)
	
	// Verify the update
	updatedLoan, err := repo.FindByID(loan.ID)
	require.NoError(t, err)
	
	assert.Equal(t, 30000.00, updatedLoan.PrincipalAmount)
	assert.Equal(t, 5.0, updatedLoan.Rate)
}

func TestDeleteLoan(t *testing.T) {
	repo, _ := setupTestRepository()
	
	// Create a loan
	loan := &domain.Loan{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
		AgreementLetterLink: "https://example.com/agreement/user123.pdf",
	}
	
	err := repo.Create(loan)
	require.NoError(t, err)
	
	// Delete the loan
	err = repo.Delete(loan.ID)
	require.NoError(t, err)
	
	// Verify the loan is deleted
	_, err = repo.FindByID(loan.ID)
	assert.Error(t, err)
}
