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
		responses = append(responses, LoanResponse{
			ID:           loan.ID,
			BorrowerID:   loan.BorrowerID,
			Amount:       loan.Amount,
			InterestRate: loan.InterestRate,
			Term:         loan.Term,
			Purpose:      loan.Purpose,
			Status:       loan.Status,
			Description:  loan.Description,
			CreatedAt:    loan.CreatedAt,
			UpdatedAt:    loan.UpdatedAt,
		})
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
	
	response := LoanResponse{
		ID:           loan.ID,
		BorrowerID:   loan.BorrowerID,
		Amount:       loan.Amount,
		InterestRate: loan.InterestRate,
		Term:         loan.Term,
		Purpose:      loan.Purpose,
		Status:       loan.Status,
		Description:  loan.Description,
		CreatedAt:    loan.CreatedAt,
		UpdatedAt:    loan.UpdatedAt,
	}
	
	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Loan retrieved successfully",
		Data:    response,
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
		BorrowerID:   req.BorrowerID,
		Amount:       req.Amount,
		InterestRate: req.InterestRate,
		Term:         req.Term,
		Purpose:      req.Purpose,
		Description:  req.Description,
		Status:       StatusProposed,
	}
	
	if err := db.Create(&loan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}
	
	response := LoanResponse{
		ID:           loan.ID,
		BorrowerID:   loan.BorrowerID,
		Amount:       loan.Amount,
		InterestRate: loan.InterestRate,
		Term:         loan.Term,
		Purpose:      loan.Purpose,
		Status:       loan.Status,
		Description:  loan.Description,
		CreatedAt:    loan.CreatedAt,
		UpdatedAt:    loan.UpdatedAt,
	}
	
	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Loan created successfully",
		Data:    response,
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
	
	// Only allow updates if loan is in proposed status
	if loan.Status != StatusProposed {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid operation",
			Message: "Can only update loans in proposed status",
		})
		return
	}
	
	// Update fields if provided
	if req.Amount != nil {
		loan.Amount = *req.Amount
	}
	if req.InterestRate != nil {
		loan.InterestRate = *req.InterestRate
	}
	if req.Term != nil {
		loan.Term = *req.Term
	}
	if req.Purpose != nil {
		loan.Purpose = *req.Purpose
	}
	if req.Description != nil {
		loan.Description = *req.Description
	}
	
	if err := db.Save(&loan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}
	
	response := LoanResponse{
		ID:           loan.ID,
		BorrowerID:   loan.BorrowerID,
		Amount:       loan.Amount,
		InterestRate: loan.InterestRate,
		Term:         loan.Term,
		Purpose:      loan.Purpose,
		Status:       loan.Status,
		Description:  loan.Description,
		CreatedAt:    loan.CreatedAt,
		UpdatedAt:    loan.UpdatedAt,
	}
	
	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Loan updated successfully",
		Data:    response,
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
	
	// Only allow deletion if loan is in proposed status
	if loan.Status != StatusProposed {
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

// approveLoan approves a loan
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
	
	// Only allow approval if loan is in proposed status
	if loan.Status != StatusProposed {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid operation",
			Message: "Can only approve loans in proposed status",
		})
		return
	}
	
	loan.Status = StatusApproved
	if err := db.Save(&loan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}
	
	response := LoanResponse{
		ID:           loan.ID,
		BorrowerID:   loan.BorrowerID,
		Amount:       loan.Amount,
		InterestRate: loan.InterestRate,
		Term:         loan.Term,
		Purpose:      loan.Purpose,
		Status:       loan.Status,
		Description:  loan.Description,
		CreatedAt:    loan.CreatedAt,
		UpdatedAt:    loan.UpdatedAt,
	}
	
	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Loan approved successfully",
		Data:    response,
	})
}

// investLoan marks a loan as invested
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
	
	// Only allow investment if loan is in approved status
	if loan.Status != StatusApproved {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid operation",
			Message: "Can only invest in approved loans",
		})
		return
	}
	
	loan.Status = StatusInvested
	if err := db.Save(&loan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}
	
	response := LoanResponse{
		ID:           loan.ID,
		BorrowerID:   loan.BorrowerID,
		Amount:       loan.Amount,
		InterestRate: loan.InterestRate,
		Term:         loan.Term,
		Purpose:      loan.Purpose,
		Status:       loan.Status,
		Description:  loan.Description,
		CreatedAt:    loan.CreatedAt,
		UpdatedAt:    loan.UpdatedAt,
	}
	
	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Loan marked as invested successfully",
		Data:    response,
	})
}

// disburseLoan disburses a loan
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
	
	// Only allow disbursement if loan is in invested status
	if loan.Status != StatusInvested {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid operation",
			Message: "Can only disburse invested loans",
		})
		return
	}
	
	loan.Status = StatusDisbursed
	if err := db.Save(&loan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Database error",
			Message: err.Error(),
		})
		return
	}
	
	response := LoanResponse{
		ID:           loan.ID,
		BorrowerID:   loan.BorrowerID,
		Amount:       loan.Amount,
		InterestRate: loan.InterestRate,
		Term:         loan.Term,
		Purpose:      loan.Purpose,
		Status:       loan.Status,
		Description:  loan.Description,
		CreatedAt:    loan.CreatedAt,
		UpdatedAt:    loan.UpdatedAt,
	}
	
	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Loan disbursed successfully",
		Data:    response,
	})
} 