package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// getLoans retrieves all loans with optional filtering
func getLoans(c *gin.Context) {
	var loans []Loan

	// Get query parameters for filtering
	status := c.Query("status")
	borrowerID := c.Query("borrower_id")

	query := db

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if borrowerID != "" {
		query = query.Where("borrower_id = ?", borrowerID)
	}

	if err := query.Find(&loans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}

	// Convert to response format
	var responses []LoanResponse
	for _, loan := range loans {
		responses = append(responses, loanToResponse(loan))
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Loans retrieved successfully",
		Data:    responses,
	})
}

// getLoan retrieves a specific loan by ID
func getLoan(c *gin.Context) {
	id := c.Param("id")

	var loan Loan
	if err := db.First(&loan, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Not found",
				Message: "Loan not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Loan retrieved successfully",
		Data:    loanToResponse(loan),
	})
}

// createLoan creates a new loan
func createLoan(c *gin.Context) {
	var req CreateLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation error",
			Message: err.Error(),
		})
		return
	}

	loan := Loan{
		BorrowerID:        req.BorrowerID,
		PrincipalAmount:   req.PrincipalAmount,
		Rate:              req.Rate,
		ROI:               req.ROI,
		AgreementLetterLink: req.AgreementLetterLink,
		Status:            StatusProposed,
	}

	if err := db.Create(&loan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Loan created successfully",
		Data:    loanToResponse(loan),
	})
}

// updateLoan updates an existing loan
func updateLoan(c *gin.Context) {
	id := c.Param("id")

	var req UpdateLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation error",
			Message: err.Error(),
		})
		return
	}

	var loan Loan
	if err := db.First(&loan, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Not found",
				Message: "Loan not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}

	// Check if loan can be updated using FSM
	if !loan.CanUpdate() {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid operation",
			Message: "Can only update loans in proposed status",
		})
		return
	}

	// Update fields if provided
	if req.PrincipalAmount != nil {
		loan.PrincipalAmount = *req.PrincipalAmount
	}
	if req.Rate != nil {
		loan.Rate = *req.Rate
	}
	if req.ROI != nil {
		loan.ROI = *req.ROI
	}
	if req.AgreementLetterLink != nil {
		loan.AgreementLetterLink = *req.AgreementLetterLink
	}

	if err := db.Save(&loan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Loan updated successfully",
		Data:    loanToResponse(loan),
	})
}

// deleteLoan deletes a loan
func deleteLoan(c *gin.Context) {
	id := c.Param("id")

	var loan Loan
	if err := db.First(&loan, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Not found",
				Message: "Loan not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}

	// Check if loan can be deleted using FSM
	if !loan.CanDelete() {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid operation",
			Message: "Can only delete loans in proposed status",
		})
		return
	}

	if err := db.Delete(&loan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Loan deleted successfully",
	})
}

// approveLoan approves a loan using FSM
func approveLoan(c *gin.Context) {
	id := c.Param("id")

	var loan Loan
	if err := db.First(&loan, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Not found",
				Message: "Loan not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}

	// Use FSM to perform state transition
	fsm := loan.GetFSM()
	if err := fsm.Transition(StatusApproved); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid operation",
			Message: "Can only approve loans in proposed status",
		})
		return
	}

	// Update loan status
	loan.Status = fsm.GetCurrentState()
	if err := db.Save(&loan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Loan approved successfully",
		Data:    loanToResponse(loan),
	})
}

// investLoan marks a loan as invested using FSM
func investLoan(c *gin.Context) {
	id := c.Param("id")

	var loan Loan
	if err := db.First(&loan, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Not found",
				Message: "Loan not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}

	// Use FSM to perform state transition
	fsm := loan.GetFSM()
	if err := fsm.Transition(StatusInvested); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid operation",
			Message: "Can only invest in approved loans",
		})
		return
	}

	// Update loan status
	loan.Status = fsm.GetCurrentState()
	if err := db.Save(&loan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Loan marked as invested successfully",
		Data:    loanToResponse(loan),
	})
}

// disburseLoan disburses a loan using FSM
func disburseLoan(c *gin.Context) {
	id := c.Param("id")

	var loan Loan
	if err := db.First(&loan, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Not found",
				Message: "Loan not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}

	// Use FSM to perform state transition
	fsm := loan.GetFSM()
	if err := fsm.Transition(StatusDisbursed); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid operation",
			Message: "Can only disburse invested loans",
		})
		return
	}

	// Update loan status
	loan.Status = fsm.GetCurrentState()
	if err := db.Save(&loan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Loan disbursed successfully",
		Data:    loanToResponse(loan),
	})
}

// getLoanTransitions returns valid transitions for a loan
func getLoanTransitions(c *gin.Context) {
	id := c.Param("id")

	var loan Loan
	if err := db.First(&loan, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Not found",
				Message: "Loan not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}

	fsm := loan.GetFSM()
	validTransitions := fsm.GetValidTransitions()

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Valid transitions retrieved successfully",
		Data: gin.H{
			"current_state": fsm.GetCurrentState(),
			"transitions":   validTransitions,
		},
	})
}

// loanToResponse converts a Loan to LoanResponse
func loanToResponse(loan Loan) LoanResponse {
	return LoanResponse{
		ID:                loan.ID,
		BorrowerID:        loan.BorrowerID,
		PrincipalAmount:   loan.PrincipalAmount,
		Rate:              loan.Rate,
		ROI:               loan.ROI,
		AgreementLetterLink: loan.AgreementLetterLink,
		Status:            loan.Status,
		CreatedAt:         loan.CreatedAt,
		UpdatedAt:         loan.UpdatedAt,
	}
}