package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"loan-service/internal/domain"
	"loan-service/internal/dto"
	"loan-service/internal/repository"
	"loan-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Helper function to create float64 pointer
func float64Ptr(v float64) *float64 {
	return &v
}

func setupTestHandler() (*LoanHandler, *gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
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
	
	// Create dependencies
	loanRepo := repository.NewLoanRepository(database)
	loanService := service.NewLoanService(loanRepo)
	loanHandler := NewLoanHandler(loanService)
	
	return loanHandler, router, database
}

func TestGetLoans(t *testing.T) {
	handler, router, _ := setupTestHandler()
	
	router.GET("/loans", handler.GetLoans)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/loans", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Loans retrieved successfully", response.Message)
}

func TestGetLoansWithFilters(t *testing.T) {
	handler, router, _ := setupTestHandler()
	
	router.GET("/loans", handler.GetLoans)
	
	// Test with status filter
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/loans?status=proposed", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Loans retrieved successfully", response.Message)
	
	// Test with borrower_id filter
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/loans?borrower_id=user123", nil)
	router.ServeHTTP(w2, req2)
	
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestCreateLoan(t *testing.T) {
	handler, router, _ := setupTestHandler()
	
	router.POST("/loans", handler.CreateLoan)
	
	createReq := dto.CreateLoanRequest{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
	}
	
	reqBody, _ := json.Marshal(createReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/loans", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Loan created successfully", response.Message)
}

func TestCreateLoanInvalidRequest(t *testing.T) {
	handler, router, _ := setupTestHandler()
	
	router.POST("/loans", handler.CreateLoan)
	
	// Test with invalid JSON
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/loans", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetLoan(t *testing.T) {
	handler, router, _ := setupTestHandler()
	
	// Set up routes
	router.POST("/loans", handler.CreateLoan)
	router.GET("/loans/:id", handler.GetLoan)
	
	// First create a loan
	createReq := dto.CreateLoanRequest{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
	}
	
	reqBody, _ := json.Marshal(createReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/loans", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	var createResponse dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &createResponse)
	require.NoError(t, err)
	
	// Extract loan ID from response
	loanData := createResponse.Data.(map[string]interface{})
	loanID := loanData["id"].(string)
	require.NotEmpty(t, loanID)
	
	// Now get the specific loan
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/loans/"+loanID, nil)
	router.ServeHTTP(w2, req2)
	
	assert.Equal(t, http.StatusOK, w2.Code)
	
	var response dto.SuccessResponse
	err = json.Unmarshal(w2.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Loan retrieved successfully", response.Message)
}

func TestGetLoanNotFound(t *testing.T) {
	handler, router, _ := setupTestHandler()
	
	router.GET("/loans/:id", handler.GetLoan)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/loans/nonexistent-id", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNotFound, w.Code)
	
	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Not found", response.Error)
}

func TestUpdateLoan(t *testing.T) {
	handler, router, _ := setupTestHandler()
	
	// Set up routes
	router.POST("/loans", handler.CreateLoan)
	router.PUT("/loans/:id", handler.UpdateLoan)
	
	// First create a loan
	createReq := dto.CreateLoanRequest{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
	}
	
	reqBody, _ := json.Marshal(createReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/loans", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	var createResponse dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &createResponse)
	require.NoError(t, err)
	
	loanData := createResponse.Data.(map[string]interface{})
	loanID := loanData["id"].(string)
	
	// Update the loan
	updateReq := dto.UpdateLoanRequest{
		PrincipalAmount: float64Ptr(30000.00),
		Rate:            float64Ptr(5.0),
	}
	
	reqBody2, _ := json.Marshal(updateReq)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("PUT", "/loans/"+loanID, bytes.NewBuffer(reqBody2))
	req2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w2, req2)
	
	assert.Equal(t, http.StatusOK, w2.Code)
	
	var response dto.SuccessResponse
	err = json.Unmarshal(w2.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Loan updated successfully", response.Message)
}

func TestUpdateLoanInvalidRequest(t *testing.T) {
	handler, router, _ := setupTestHandler()
	
	router.PUT("/loans/:id", handler.UpdateLoan)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/loans/test-id", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteLoan(t *testing.T) {
	handler, router, _ := setupTestHandler()
	
	// Set up routes
	router.POST("/loans", handler.CreateLoan)
	router.DELETE("/loans/:id", handler.DeleteLoan)
	
	// First create a loan
	createReq := dto.CreateLoanRequest{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
	}
	
	reqBody, _ := json.Marshal(createReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/loans", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	var createResponse dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &createResponse)
	require.NoError(t, err)
	
	loanData := createResponse.Data.(map[string]interface{})
	loanID := loanData["id"].(string)
	
	// Delete the loan
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("DELETE", "/loans/"+loanID, nil)
	router.ServeHTTP(w2, req2)
	
	assert.Equal(t, http.StatusOK, w2.Code)
	
	var response dto.SuccessResponse
	err = json.Unmarshal(w2.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Loan deleted successfully", response.Message)
}

func TestApproveLoan(t *testing.T) {
	handler, router, _ := setupTestHandler()
	
	// Set up routes
	router.POST("/loans", handler.CreateLoan)
	router.PUT("/loans/:id/approve", handler.ApproveLoan)
	
	// First create a loan
	createReq := dto.CreateLoanRequest{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
	}
	
	reqBody, _ := json.Marshal(createReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/loans", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	var createResponse dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &createResponse)
	require.NoError(t, err)
	
	loanData := createResponse.Data.(map[string]interface{})
	loanID := loanData["id"].(string)
	
	// Approve the loan
	approveReq := dto.ApproveLoanRequest{
		FieldValidatorProof: "https://example.com/proofs/field_visit_123.jpg",
		FieldValidatorID:    "validator_001",
	}
	
	reqBody2, _ := json.Marshal(approveReq)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("PUT", "/loans/"+loanID+"/approve", bytes.NewBuffer(reqBody2))
	req2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w2, req2)
	
	assert.Equal(t, http.StatusOK, w2.Code)
	
	var response dto.SuccessResponse
	err = json.Unmarshal(w2.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Loan approved successfully", response.Message)
}

func TestInvestLoan(t *testing.T) {
	handler, router, _ := setupTestHandler()
	
	// Set up routes
	router.POST("/loans", handler.CreateLoan)
	router.PUT("/loans/:id/approve", handler.ApproveLoan)
	router.PUT("/loans/:id/invest", handler.InvestLoan)
	
	// First create a loan
	createReq := dto.CreateLoanRequest{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
	}
	
	reqBody, _ := json.Marshal(createReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/loans", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	var createResponse dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &createResponse)
	require.NoError(t, err)
	
	loanData := createResponse.Data.(map[string]interface{})
	loanID := loanData["id"].(string)
	
	// Approve the loan first
	approveReq := dto.ApproveLoanRequest{
		FieldValidatorProof: "https://example.com/proofs/field_visit_123.jpg",
		FieldValidatorID:    "validator_001",
	}
	
	reqBody2, _ := json.Marshal(approveReq)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("PUT", "/loans/"+loanID+"/approve", bytes.NewBuffer(reqBody2))
	req2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w2, req2)
	
	// Invest in the loan
	investReq := dto.InvestLoanRequest{
		InvestorID:          "investor_001",
		Amount:              10000.00,
		AgreementLetterLink: "https://example.com/agreement.pdf",
	}
	
	reqBody3, _ := json.Marshal(investReq)
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("PUT", "/loans/"+loanID+"/invest", bytes.NewBuffer(reqBody3))
	req3.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w3, req3)
	
	assert.Equal(t, http.StatusOK, w3.Code)
	
	var response dto.SuccessResponse
	err = json.Unmarshal(w3.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Investment added successfully", response.Message)
}

func TestDisburseLoan(t *testing.T) {
	handler, router, _ := setupTestHandler()
	
	// Set up routes
	router.POST("/loans", handler.CreateLoan)
	router.PUT("/loans/:id/approve", handler.ApproveLoan)
	router.PUT("/loans/:id/invest", handler.InvestLoan)
	router.PUT("/loans/:id/disburse", handler.DisburseLoan)
	
	// First create a loan
	createReq := dto.CreateLoanRequest{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
	}
	
	reqBody, _ := json.Marshal(createReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/loans", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	var createResponse dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &createResponse)
	require.NoError(t, err)
	
	loanData := createResponse.Data.(map[string]interface{})
	loanID := loanData["id"].(string)
	
	// Approve the loan
	approveReq := dto.ApproveLoanRequest{
		FieldValidatorProof: "https://example.com/proofs/field_visit_123.jpg",
		FieldValidatorID:    "validator_001",
	}
	
	reqBody2, _ := json.Marshal(approveReq)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("PUT", "/loans/"+loanID+"/approve", bytes.NewBuffer(reqBody2))
	req2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w2, req2)
	
	// Invest fully in the loan
	investReq := dto.InvestLoanRequest{
		InvestorID:          "investor_001",
		Amount:              25000.00,
		AgreementLetterLink: "https://example.com/agreement.pdf",
	}
	
	reqBody3, _ := json.Marshal(investReq)
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("PUT", "/loans/"+loanID+"/invest", bytes.NewBuffer(reqBody3))
	req3.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w3, req3)
	
	// Disburse the loan
	disburseReq := dto.DisburseLoanRequest{
		SignedAgreementLink: "https://example.com/signed-agreement.pdf",
		FieldOfficerID:      "officer_001",
	}
	
	reqBody4, _ := json.Marshal(disburseReq)
	w4 := httptest.NewRecorder()
	req4, _ := http.NewRequest("PUT", "/loans/"+loanID+"/disburse", bytes.NewBuffer(reqBody4))
	req4.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w4, req4)
	
	assert.Equal(t, http.StatusOK, w4.Code)
	
	var response dto.SuccessResponse
	err = json.Unmarshal(w4.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Loan disbursed successfully", response.Message)
}

func TestGetLoanTransitions(t *testing.T) {
	handler, router, _ := setupTestHandler()
	
	// Set up routes
	router.POST("/loans", handler.CreateLoan)
	router.GET("/loans/:id/transitions", handler.GetLoanTransitions)
	
	// First create a loan
	createReq := dto.CreateLoanRequest{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
	}
	
	reqBody, _ := json.Marshal(createReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/loans", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	var createResponse dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &createResponse)
	require.NoError(t, err)
	
	loanData := createResponse.Data.(map[string]interface{})
	loanID := loanData["id"].(string)
	
	// Get loan transitions
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/loans/"+loanID+"/transitions", nil)
	router.ServeHTTP(w2, req2)
	
	assert.Equal(t, http.StatusOK, w2.Code)
	
	var response dto.SuccessResponse
	err = json.Unmarshal(w2.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Valid transitions retrieved successfully", response.Message)
}
