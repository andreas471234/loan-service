package handler

import (
	"net/http"

	"loan-service/internal/domain"
	"loan-service/internal/dto"
	"loan-service/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LoanHandler handles HTTP requests for loan operations
type LoanHandler struct {
	loanService service.LoanService
}

// NewLoanHandler creates a new loan handler
func NewLoanHandler(loanService service.LoanService) *LoanHandler {
	return &LoanHandler{
		loanService: loanService,
	}
}

// GetLoans retrieves all loans with optional filtering
func (h *LoanHandler) GetLoans(c *gin.Context) {
	filters := make(map[string]interface{})

	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	if borrowerID := c.Query("borrower_id"); borrowerID != "" {
		filters["borrower_id"] = borrowerID
	}

	loans, err := h.loanService.GetLoans(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}

	var responses []dto.LoanResponse
	for _, loan := range loans {
		responses = append(responses, dto.ToLoanResponse(loan))
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Loans retrieved successfully",
		Data:    responses,
	})
}

// GetLoan retrieves a specific loan by ID
func (h *LoanHandler) GetLoan(c *gin.Context) {
	id := c.Param("id")

	loan, err := h.loanService.GetLoan(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not found",
				Message: "Loan not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Loan retrieved successfully",
		Data:    dto.ToLoanResponse(*loan),
	})
}

// CreateLoan creates a new loan
func (h *LoanHandler) CreateLoan(c *gin.Context) {
	var req dto.CreateLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Validation error",
			Message: err.Error(),
		})
		return
	}

	loan := &domain.Loan{
		BorrowerID:      req.BorrowerID,
		PrincipalAmount: req.PrincipalAmount,
		Rate:            req.Rate,
		ROI:             req.ROI,
	}

	if err := h.loanService.CreateLoan(loan); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Message: "Loan created successfully",
		Data:    dto.ToLoanResponse(*loan),
	})
}

// UpdateLoan updates an existing loan
func (h *LoanHandler) UpdateLoan(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Validation error",
			Message: err.Error(),
		})
		return
	}

	updates := make(map[string]interface{})
	if req.PrincipalAmount != nil {
		updates["principal_amount"] = *req.PrincipalAmount
	}
	if req.Rate != nil {
		updates["rate"] = *req.Rate
	}
	if req.ROI != nil {
		updates["roi"] = *req.ROI
	}
	if req.AgreementLetterLink != nil {
		updates["agreement_letter_link"] = *req.AgreementLetterLink
	}

	loan, err := h.loanService.UpdateLoan(id, updates)
	if err != nil {
		if err.Error() == "can only update loans in proposed status" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid operation",
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Loan updated successfully",
		Data:    dto.ToLoanResponse(*loan),
	})
}

// DeleteLoan deletes a loan
func (h *LoanHandler) DeleteLoan(c *gin.Context) {
	id := c.Param("id")

	if err := h.loanService.DeleteLoan(id); err != nil {
		if err.Error() == "can only delete loans in proposed status" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid operation",
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Loan deleted successfully",
	})
}

// ApproveLoan approves a loan
func (h *LoanHandler) ApproveLoan(c *gin.Context) {
	id := c.Param("id")

	var req dto.ApproveLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Validation error",
			Message: err.Error(),
		})
		return
	}

	approvalDetails := &domain.ApprovalDetails{
		FieldValidatorProof: req.FieldValidatorProof,
		FieldValidatorID:    req.FieldValidatorID,
	}

	loan, err := h.loanService.ApproveLoan(id, approvalDetails)
	if err != nil {
		if err.Error() == "can only approve loans in proposed status" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid operation",
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Loan approved successfully",
		Data:    dto.ToLoanResponse(*loan),
	})
}

// InvestLoan adds an investment to a loan
func (h *LoanHandler) InvestLoan(c *gin.Context) {
	id := c.Param("id")

	var req dto.InvestLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Validation error",
			Message: err.Error(),
		})
		return
	}

	loan, err := h.loanService.InvestInLoan(id, req.InvestorID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Investment error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Investment added successfully",
		Data:    dto.ToLoanResponse(*loan),
	})
}

// DisburseLoan disburses a loan
func (h *LoanHandler) DisburseLoan(c *gin.Context) {
	id := c.Param("id")

	var req dto.DisburseLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Validation error",
			Message: err.Error(),
		})
		return
	}

	disbursementDetails := &domain.DisbursementDetails{
		SignedAgreementLink: req.SignedAgreementLink,
		FieldOfficerID:      req.FieldOfficerID,
	}

	loan, err := h.loanService.DisburseLoan(id, disbursementDetails)
	if err != nil {
		if err.Error() == "can only disburse fully invested loans" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid operation",
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Loan disbursed successfully",
		Data:    dto.ToLoanResponse(*loan),
	})
}

// GetLoanTransitions returns valid transitions for a loan
func (h *LoanHandler) GetLoanTransitions(c *gin.Context) {
	id := c.Param("id")

	transitions, err := h.loanService.GetLoanTransitions(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not found",
				Message: "Loan not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}

	loan, _ := h.loanService.GetLoan(id)
	fsm := domain.NewFSM()
	fsm.SetCurrentState(loan.Status)

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Valid transitions retrieved successfully",
		Data: dto.TransitionResponse{
			CurrentState: fsm.GetCurrentState(),
			Transitions:  transitions,
		},
	})
}
