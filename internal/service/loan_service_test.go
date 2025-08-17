package service

import (
	"testing"

	"loan-service/internal/domain"
	"loan-service/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestService() (*loanService, *gorm.DB) {
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

	loanRepo := repository.NewLoanRepository(database)
	loanService := NewLoanService(loanRepo).(*loanService)
	return loanService, database
}

func TestCreateLoan(t *testing.T) {
	service, _ := setupTestService()

	loan := &domain.Loan{
		BorrowerID:      "user123",
		PrincipalAmount: 25000.00,
		Rate:            4.5,
		ROI:             6.0,
	}

	err := service.CreateLoan(loan)
	require.NoError(t, err)

	assert.NotEmpty(t, loan.ID)
	assert.Equal(t, domain.StatusProposed, loan.Status)
	assert.Equal(t, 0.0, loan.TotalInvested)
}

func TestGetLoan(t *testing.T) {
	service, _ := setupTestService()

	// Create a loan first
	loan := &domain.Loan{
		BorrowerID:      "user123",
		PrincipalAmount: 25000.00,
		Rate:            4.5,
		ROI:             6.0,
	}

	err := service.CreateLoan(loan)
	require.NoError(t, err)

	// Get the loan
	retrievedLoan, err := service.GetLoan(loan.ID)
	require.NoError(t, err)

	assert.Equal(t, loan.ID, retrievedLoan.ID)
	assert.Equal(t, loan.BorrowerID, retrievedLoan.BorrowerID)
	assert.Equal(t, loan.PrincipalAmount, retrievedLoan.PrincipalAmount)
}

func TestGetLoanNotFound(t *testing.T) {
	service, _ := setupTestService()

	_, err := service.GetLoan("nonexistent-id")
	assert.Error(t, err)
}

func TestGetLoans(t *testing.T) {
	service, _ := setupTestService()

	// Create multiple loans
	loan1 := &domain.Loan{
		BorrowerID:      "user123",
		PrincipalAmount: 25000.00,
		Rate:            4.5,
		ROI:             6.0,
	}

	loan2 := &domain.Loan{
		BorrowerID:      "user456",
		PrincipalAmount: 30000.00,
		Rate:            5.0,
		ROI:             7.0,
	}

	err := service.CreateLoan(loan1)
	require.NoError(t, err)

	err = service.CreateLoan(loan2)
	require.NoError(t, err)

	// Get all loans
	loans, err := service.GetLoans(nil)
	require.NoError(t, err)

	assert.Len(t, loans, 2)
}

func TestGetLoansWithFilters(t *testing.T) {
	service, db := setupTestService()

	// Create loans with different statuses
	loan1 := &domain.Loan{
		BorrowerID:      "user123",
		PrincipalAmount: 25000.00,
		Rate:            4.5,
		ROI:             6.0,
		Status:          domain.StatusProposed,
	}

	loan2 := &domain.Loan{
		BorrowerID:      "user456",
		PrincipalAmount: 30000.00,
		Rate:            5.0,
		ROI:             7.0,
		Status:          domain.StatusApproved,
	}

	// Create loans directly in database to set specific statuses
	err := db.Create(loan1).Error
	require.NoError(t, err)

	err = db.Create(loan2).Error
	require.NoError(t, err)

	// Get loans with status filter
	filters := map[string]interface{}{
		"status": domain.StatusProposed,
	}

	loans, err := service.GetLoans(filters)
	require.NoError(t, err)

	assert.Len(t, loans, 1)
	assert.Equal(t, domain.StatusProposed, loans[0].Status)
}

func TestUpdateLoan(t *testing.T) {
	service, _ := setupTestService()

	// Create a loan first
	loan := &domain.Loan{
		BorrowerID:      "user123",
		PrincipalAmount: 25000.00,
		Rate:            4.5,
		ROI:             6.0,
	}

	err := service.CreateLoan(loan)
	require.NoError(t, err)

	// Update the loan
	updates := map[string]interface{}{
		"principal_amount": 30000.00,
		"rate":             5.0,
	}

	updatedLoan, err := service.UpdateLoan(loan.ID, updates)
	require.NoError(t, err)

	assert.Equal(t, 30000.00, updatedLoan.PrincipalAmount)
	assert.Equal(t, 5.0, updatedLoan.Rate)
}

func TestUpdateLoanNotFound(t *testing.T) {
	service, _ := setupTestService()

	updates := map[string]interface{}{
		"principal_amount": 30000.00,
	}

	_, err := service.UpdateLoan("nonexistent-id", updates)
	assert.Error(t, err)
}

func TestUpdateLoanInvalidState(t *testing.T) {
	service, _ := setupTestService()

	// Create and approve a loan
	loan := &domain.Loan{
		BorrowerID:      "user123",
		PrincipalAmount: 25000.00,
		Rate:            4.5,
		ROI:             6.0,
	}

	err := service.CreateLoan(loan)
	require.NoError(t, err)

	// Approve the loan
	approvalDetails := &domain.ApprovalDetails{
		FieldValidatorProof: "proof",
		FieldValidatorID:    "validator_001",
	}

	_, err = service.ApproveLoan(loan.ID, approvalDetails)
	require.NoError(t, err)

	// Try to update approved loan (should fail)
	updates := map[string]interface{}{
		"principal_amount": 30000.00,
	}

	_, err = service.UpdateLoan(loan.ID, updates)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "can only update loans in proposed status")
}

func TestDeleteLoan(t *testing.T) {
	service, _ := setupTestService()

	// Create a loan
	loan := &domain.Loan{
		BorrowerID:      "user123",
		PrincipalAmount: 25000.00,
		Rate:            4.5,
		ROI:             6.0,
	}

	err := service.CreateLoan(loan)
	require.NoError(t, err)

	// Delete the loan
	err = service.DeleteLoan(loan.ID)
	require.NoError(t, err)

	// Verify the loan is deleted
	_, err = service.GetLoan(loan.ID)
	assert.Error(t, err)
}

func TestDeleteLoanNotFound(t *testing.T) {
	service, _ := setupTestService()

	err := service.DeleteLoan("nonexistent-id")
	assert.Error(t, err)
}

func TestDeleteLoanInvalidState(t *testing.T) {
	service, _ := setupTestService()

	// Create and approve a loan
	loan := &domain.Loan{
		BorrowerID:      "user123",
		PrincipalAmount: 25000.00,
		Rate:            4.5,
		ROI:             6.0,
	}

	err := service.CreateLoan(loan)
	require.NoError(t, err)

	// Approve the loan
	approvalDetails := &domain.ApprovalDetails{
		FieldValidatorProof: "proof",
		FieldValidatorID:    "validator_001",
	}

	_, err = service.ApproveLoan(loan.ID, approvalDetails)
	require.NoError(t, err)

	// Try to delete approved loan (should fail)
	err = service.DeleteLoan(loan.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "can only delete loans in proposed status")
}

func TestApproveLoan(t *testing.T) {
	service, _ := setupTestService()

	// Create a loan
	loan := &domain.Loan{
		BorrowerID:      "user123",
		PrincipalAmount: 25000.00,
		Rate:            4.5,
		ROI:             6.0,
	}

	err := service.CreateLoan(loan)
	require.NoError(t, err)

	// Approve the loan
	approvalDetails := &domain.ApprovalDetails{
		FieldValidatorProof: "proof",
		FieldValidatorID:    "validator_001",
	}

	approvedLoan, err := service.ApproveLoan(loan.ID, approvalDetails)
	require.NoError(t, err)

	assert.Equal(t, domain.StatusApproved, approvedLoan.Status)
	assert.NotNil(t, approvedLoan.ApprovalDetails)
	assert.Equal(t, "proof", approvedLoan.ApprovalDetails.FieldValidatorProof)
	assert.Equal(t, "validator_001", approvedLoan.ApprovalDetails.FieldValidatorID)
}

func TestApproveLoanNotFound(t *testing.T) {
	service, _ := setupTestService()

	approvalDetails := &domain.ApprovalDetails{
		FieldValidatorProof: "proof",
		FieldValidatorID:    "validator_001",
	}

	_, err := service.ApproveLoan("nonexistent-id", approvalDetails)
	assert.Error(t, err)
}

func TestApproveLoanInvalidState(t *testing.T) {
	service, _ := setupTestService()

	// Create and approve a loan
	loan := &domain.Loan{
		BorrowerID:      "user123",
		PrincipalAmount: 25000.00,
		Rate:            4.5,
		ROI:             6.0,
	}

	err := service.CreateLoan(loan)
	require.NoError(t, err)

	// Approve the loan
	approvalDetails := &domain.ApprovalDetails{
		FieldValidatorProof: "proof",
		FieldValidatorID:    "validator_001",
	}

	_, err = service.ApproveLoan(loan.ID, approvalDetails)
	require.NoError(t, err)

	// Try to approve again (should fail)
	_, err = service.ApproveLoan(loan.ID, approvalDetails)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "can only approve loans in proposed status")
}

func TestInvestInLoan(t *testing.T) {
	service, _ := setupTestService()

	// Create and approve a loan
	loan := &domain.Loan{
		BorrowerID:      "user123",
		PrincipalAmount: 25000.00,
		Rate:            4.5,
		ROI:             6.0,
	}

	err := service.CreateLoan(loan)
	require.NoError(t, err)

	approvalDetails := &domain.ApprovalDetails{
		FieldValidatorProof: "proof",
		FieldValidatorID:    "validator_001",
	}

	_, err = service.ApproveLoan(loan.ID, approvalDetails)
	require.NoError(t, err)

	// Invest in the loan
	investedLoan, err := service.InvestInLoan(loan.ID, "investor_001", 10000.00)
	require.NoError(t, err)

	assert.Equal(t, 10000.00, investedLoan.TotalInvested)
	assert.Len(t, investedLoan.Investments, 1)
	assert.Equal(t, "investor_001", investedLoan.Investments[0].InvestorID)
	assert.Equal(t, 10000.00, investedLoan.Investments[0].Amount)
}

func TestInvestInLoanNotFound(t *testing.T) {
	service, _ := setupTestService()

	_, err := service.InvestInLoan("nonexistent-id", "investor_001", 10000.00)
	assert.Error(t, err)
}

func TestInvestInLoanInvalidState(t *testing.T) {
	service, _ := setupTestService()

	// Create a loan (not approved)
	loan := &domain.Loan{
		BorrowerID:      "user123",
		PrincipalAmount: 25000.00,
		Rate:            4.5,
		ROI:             6.0,
	}

	err := service.CreateLoan(loan)
	require.NoError(t, err)

	// Try to invest in unapproved loan (should fail)
	_, err = service.InvestInLoan(loan.ID, "investor_001", 10000.00)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "loan is not in approved status")
}

func TestInvestInLoanExceedsLimit(t *testing.T) {
	service, _ := setupTestService()

	// Create and approve a loan
	loan := &domain.Loan{
		BorrowerID:      "user123",
		PrincipalAmount: 25000.00,
		Rate:            4.5,
		ROI:             6.0,
	}

	err := service.CreateLoan(loan)
	require.NoError(t, err)

	approvalDetails := &domain.ApprovalDetails{
		FieldValidatorProof: "proof",
		FieldValidatorID:    "validator_001",
	}

	_, err = service.ApproveLoan(loan.ID, approvalDetails)
	require.NoError(t, err)

	// Try to invest more than principal amount (should fail)
	_, err = service.InvestInLoan(loan.ID, "investor_001", 30000.00)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "total investment amount would exceed loan principal")
}

func TestDisburseLoan(t *testing.T) {
	service, _ := setupTestService()

	// Create, approve, and fully invest in a loan
	loan := &domain.Loan{
		BorrowerID:      "user123",
		PrincipalAmount: 25000.00,
		Rate:            4.5,
		ROI:             6.0,
	}

	err := service.CreateLoan(loan)
	require.NoError(t, err)

	approvalDetails := &domain.ApprovalDetails{
		FieldValidatorProof: "proof",
		FieldValidatorID:    "validator_001",
	}

	_, err = service.ApproveLoan(loan.ID, approvalDetails)
	require.NoError(t, err)

	_, err = service.InvestInLoan(loan.ID, "investor_001", 25000.00)
	require.NoError(t, err)

	// Disburse the loan
	disbursementDetails := &domain.DisbursementDetails{
		SignedAgreementLink: "https://example.com/signed-agreement.pdf",
		FieldOfficerID:      "officer_001",
	}

	disbursedLoan, err := service.DisburseLoan(loan.ID, disbursementDetails)
	require.NoError(t, err)

	assert.Equal(t, domain.StatusDisbursed, disbursedLoan.Status)
	assert.NotNil(t, disbursedLoan.DisbursementDetails)
	assert.Equal(t, "https://example.com/signed-agreement.pdf", disbursedLoan.DisbursementDetails.SignedAgreementLink)
	assert.Equal(t, "officer_001", disbursedLoan.DisbursementDetails.FieldOfficerID)
}

func TestDisburseLoanNotFound(t *testing.T) {
	service, _ := setupTestService()

	disbursementDetails := &domain.DisbursementDetails{
		SignedAgreementLink: "https://example.com/signed-agreement.pdf",
		FieldOfficerID:      "officer_001",
	}

	_, err := service.DisburseLoan("nonexistent-id", disbursementDetails)
	assert.Error(t, err)
}

func TestDisburseLoanNotFullyInvested(t *testing.T) {
	service, _ := setupTestService()

	// Create and approve a loan
	loan := &domain.Loan{
		BorrowerID:      "user123",
		PrincipalAmount: 25000.00,
		Rate:            4.5,
		ROI:             6.0,
	}

	err := service.CreateLoan(loan)
	require.NoError(t, err)

	approvalDetails := &domain.ApprovalDetails{
		FieldValidatorProof: "proof",
		FieldValidatorID:    "validator_001",
	}

	_, err = service.ApproveLoan(loan.ID, approvalDetails)
	require.NoError(t, err)

	// Invest partially
	_, err = service.InvestInLoan(loan.ID, "investor_001", 10000.00)
	require.NoError(t, err)

	// Try to disburse partially invested loan (should fail)
	disbursementDetails := &domain.DisbursementDetails{
		SignedAgreementLink: "https://example.com/signed-agreement.pdf",
		FieldOfficerID:      "officer_001",
	}

	_, err = service.DisburseLoan(loan.ID, disbursementDetails)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "can only disburse fully invested loans")
}

func TestGetLoanTransitions(t *testing.T) {
	service, _ := setupTestService()

	// Create a loan
	loan := &domain.Loan{
		BorrowerID:      "user123",
		PrincipalAmount: 25000.00,
		Rate:            4.5,
		ROI:             6.0,
	}

	err := service.CreateLoan(loan)
	require.NoError(t, err)

	// Get transitions for proposed loan
	transitions, err := service.GetLoanTransitions(loan.ID)
	require.NoError(t, err)

	assert.Len(t, transitions, 1)
	assert.Equal(t, domain.StatusApproved, transitions[0].To)
	assert.Equal(t, "approve", transitions[0].Action)
}

func TestGetLoanTransitionsNotFound(t *testing.T) {
	service, _ := setupTestService()

	_, err := service.GetLoanTransitions("nonexistent-id")
	assert.Error(t, err)
}

func TestInvestInLoanAutoGeneratesAgreementLink(t *testing.T) {
	service, _ := setupTestService()

	// Create and approve a loan
	loan := &domain.Loan{
		BorrowerID:      "user123",
		PrincipalAmount: 10000.00,
		Rate:            4.5,
		ROI:             6.0,
	}

	err := service.CreateLoan(loan)
	require.NoError(t, err)

	approvalDetails := &domain.ApprovalDetails{
		FieldValidatorProof: "proof",
		FieldValidatorID:    "validator_001",
	}

	_, err = service.ApproveLoan(loan.ID, approvalDetails)
	require.NoError(t, err)

	// Invest fully in the loan (should auto-generate agreement letter link)
	investedLoan, err := service.InvestInLoan(loan.ID, "investor_001", 10000.00)
	require.NoError(t, err)

	// Verify the loan is now invested and has auto-generated agreement letter link
	assert.Equal(t, domain.StatusInvested, investedLoan.Status)
	assert.Equal(t, 10000.00, investedLoan.TotalInvested)
	assert.NotEmpty(t, investedLoan.AgreementLetterLink)
	assert.Contains(t, investedLoan.AgreementLetterLink, "https://example.com/agreements/loan_")
	assert.Contains(t, investedLoan.AgreementLetterLink, "_agreement.pdf")
}
