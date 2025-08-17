package service

import (
	"errors"
	"fmt"
	"time"

	"loan-service/internal/domain"
	"loan-service/internal/repository"
)

// generateAgreementLetterLink creates a dummy agreement letter link
func generateAgreementLetterLink(loanID string) string {
	return fmt.Sprintf("https://example.com/agreements/loan_%s_agreement.pdf", loanID)
}

// LoanService defines the interface for loan business logic
type LoanService interface {
	CreateLoan(loan *domain.Loan) error
	GetLoan(id string) (*domain.Loan, error)
	GetLoans(filters map[string]interface{}) ([]domain.Loan, error)
	UpdateLoan(id string, updates map[string]interface{}) (*domain.Loan, error)
	DeleteLoan(id string) error
	ApproveLoan(id string, approvalDetails *domain.ApprovalDetails) (*domain.Loan, error)
	InvestInLoan(id string, investorID string, amount float64) (*domain.Loan, error)
	DisburseLoan(id string, disbursementDetails *domain.DisbursementDetails) (*domain.Loan, error)
	GetLoanTransitions(id string) ([]domain.StateTransition, error)
}

// loanService implements LoanService
type loanService struct {
	repo repository.LoanRepository
}

// NewLoanService creates a new loan service
func NewLoanService(repo repository.LoanRepository) LoanService {
	return &loanService{repo: repo}
}

// CreateLoan creates a new loan
func (s *loanService) CreateLoan(loan *domain.Loan) error {
	loan.Status = domain.StatusProposed
	loan.TotalInvested = 0
	return s.repo.Create(loan)
}

// GetLoan retrieves a loan by ID
func (s *loanService) GetLoan(id string) (*domain.Loan, error) {
	return s.repo.FindByID(id)
}

// GetLoans retrieves all loans with optional filters
func (s *loanService) GetLoans(filters map[string]interface{}) ([]domain.Loan, error) {
	return s.repo.FindAll(filters)
}

// UpdateLoan updates a loan
func (s *loanService) UpdateLoan(id string, updates map[string]interface{}) (*domain.Loan, error) {
	loan, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if !loan.CanUpdate() {
		return nil, errors.New("can only update loans in proposed status")
	}

	// Apply updates
	if principalAmount, ok := updates["principal_amount"].(float64); ok {
		loan.PrincipalAmount = principalAmount
	}
	if rate, ok := updates["rate"].(float64); ok {
		loan.Rate = rate
	}
	if roi, ok := updates["roi"].(float64); ok {
		loan.ROI = roi
	}
	if agreementLetterLink, ok := updates["agreement_letter_link"].(string); ok {
		loan.AgreementLetterLink = agreementLetterLink
	}

	err = s.repo.Update(loan)
	if err != nil {
		return nil, err
	}

	return loan, nil
}

// DeleteLoan deletes a loan
func (s *loanService) DeleteLoan(id string) error {
	loan, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if !loan.CanDelete() {
		return errors.New("can only delete loans in proposed status")
	}

	return s.repo.Delete(id)
}

// ApproveLoan approves a loan
func (s *loanService) ApproveLoan(id string, approvalDetails *domain.ApprovalDetails) (*domain.Loan, error) {
	loan, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if !loan.CanApprove() {
		return nil, errors.New("can only approve loans in proposed status")
	}

	fsm := domain.NewFSM()
	fsm.SetCurrentState(loan.Status)
	if err := fsm.Transition(domain.StatusApproved); err != nil {
		return nil, err
	}

	loan.Status = fsm.GetCurrentState()
	loan.ApprovalDetails = approvalDetails
	loan.ApprovalDetails.ApprovalDate = time.Now()

	err = s.repo.Update(loan)
	if err != nil {
		return nil, err
	}

	return loan, nil
}

// InvestInLoan adds an investment to a loan
func (s *loanService) InvestInLoan(id string, investorID string, amount float64) (*domain.Loan, error) {
	loan, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err := loan.AddInvestment(investorID, amount); err != nil {
		return nil, err
	}

	// Auto-generate the agreement letter link when the loan becomes invested
	if loan.Status == domain.StatusInvested {
		loan.AgreementLetterLink = generateAgreementLetterLink(loan.ID)
	}

	err = s.repo.Update(loan)
	if err != nil {
		return nil, err
	}

	return loan, nil
}

// DisburseLoan disburses a loan
func (s *loanService) DisburseLoan(id string, disbursementDetails *domain.DisbursementDetails) (*domain.Loan, error) {
	loan, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if !loan.CanDisburse() {
		return nil, errors.New("can only disburse fully invested loans")
	}

	fsm := domain.NewFSM()
	fsm.SetCurrentState(loan.Status)
	if err := fsm.Transition(domain.StatusDisbursed); err != nil {
		return nil, err
	}

	loan.Status = fsm.GetCurrentState()
	loan.DisbursementDetails = disbursementDetails
	loan.DisbursementDetails.DisbursementDate = time.Now()

	err = s.repo.Update(loan)
	if err != nil {
		return nil, err
	}

	return loan, nil
}

// GetLoanTransitions returns valid transitions for a loan
func (s *loanService) GetLoanTransitions(id string) ([]domain.StateTransition, error) {
	loan, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	fsm := domain.NewFSM()
	fsm.SetCurrentState(loan.Status)
	return fsm.GetValidTransitions(), nil
}
