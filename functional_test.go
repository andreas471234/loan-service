package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"loan-service/internal/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== FUNCTIONAL TESTS ====================

// TestHealthEndpoint tests the health endpoint functionality
func TestHealthEndpoint(t *testing.T) {
	setup := setupTestServer()
	defer setup.server.Close()

	resp, err := makeRequest("GET", setup.server.URL+"/health", nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, "ok", response["status"])
}

// TestCreateLoan tests loan creation functionality
func TestCreateLoan(t *testing.T) {
	setup := setupTestServer()
	defer setup.server.Close()

	createReq := dto.CreateLoanRequest{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
		AgreementLetterLink: "https://example.com/agreement/user123.pdf",
	}

	resp, err := makeRequest("POST", setup.server.URL+"/api/v1/loans/", createReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var response dto.SuccessResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	loanData := response.Data.(map[string]interface{})
	assert.NotEmpty(t, loanData["id"])
	assert.Equal(t, "user123", loanData["borrower_id"])
	assert.Equal(t, 25000.0, loanData["principal_amount"])
	assert.Equal(t, "proposed", loanData["status"])
}

// TestGetLoan tests getting a specific loan functionality
func TestGetLoan(t *testing.T) {
	setup := setupTestServer()
	defer setup.server.Close()

	// First create a loan
	createReq := dto.CreateLoanRequest{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
		AgreementLetterLink: "https://example.com/agreement/user123.pdf",
	}

	createResp, err := makeRequest("POST", setup.server.URL+"/api/v1/loans/", createReq)
	require.NoError(t, err)
	defer createResp.Body.Close()

	var createResponse dto.SuccessResponse
	err = json.NewDecoder(createResp.Body).Decode(&createResponse)
	require.NoError(t, err)

	loanData := createResponse.Data.(map[string]interface{})
	loanID := loanData["id"].(string)

	// Now get the specific loan
	resp, err := makeRequest("GET", setup.server.URL+"/api/v1/loans/"+loanID, nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response dto.SuccessResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	retrievedLoan := response.Data.(map[string]interface{})
	assert.Equal(t, loanID, retrievedLoan["id"])
	assert.Equal(t, "user123", retrievedLoan["borrower_id"])
}

// TestGetAllLoans tests getting all loans functionality
func TestGetAllLoans(t *testing.T) {
	setup := setupTestServer()
	defer setup.server.Close()

	// Create multiple loans
	createReq1 := dto.CreateLoanRequest{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
		AgreementLetterLink: "https://example.com/agreement/user123.pdf",
	}

	createReq2 := dto.CreateLoanRequest{
		BorrowerID:          "user456",
		PrincipalAmount:     30000.00,
		Rate:                5.0,
		ROI:                 7.0,
		AgreementLetterLink: "https://example.com/agreement/user456.pdf",
	}

	_, err := makeRequest("POST", setup.server.URL+"/api/v1/loans/", createReq1)
	require.NoError(t, err)

	_, err = makeRequest("POST", setup.server.URL+"/api/v1/loans/", createReq2)
	require.NoError(t, err)

	// Get all loans
	resp, err := makeRequest("GET", setup.server.URL+"/api/v1/loans/", nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response dto.SuccessResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	loans := response.Data.([]interface{})
	assert.GreaterOrEqual(t, len(loans), 2)
}

// TestUpdateLoan tests loan update functionality
func TestUpdateLoan(t *testing.T) {
	setup := setupTestServer()
	defer setup.server.Close()

	// First create a loan
	createReq := dto.CreateLoanRequest{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
		AgreementLetterLink: "https://example.com/agreement/user123.pdf",
	}

	createResp, err := makeRequest("POST", setup.server.URL+"/api/v1/loans/", createReq)
	require.NoError(t, err)
	defer createResp.Body.Close()

	var createResponse dto.SuccessResponse
	err = json.NewDecoder(createResp.Body).Decode(&createResponse)
	require.NoError(t, err)

	loanData := createResponse.Data.(map[string]interface{})
	loanID := loanData["id"].(string)

	// Update the loan
	updateReq := dto.UpdateLoanRequest{
		PrincipalAmount:     float64Ptr(18000.00),
		Rate:                float64Ptr(5.5),
		ROI:                 float64Ptr(8.0),
		AgreementLetterLink: stringPtr("https://example.com/agreement/user123_updated.pdf"),
	}

	resp, err := makeRequest("PUT", setup.server.URL+"/api/v1/loans/"+loanID, updateReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response dto.SuccessResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	updatedLoan := response.Data.(map[string]interface{})
	assert.Equal(t, 18000.0, updatedLoan["principal_amount"])
	assert.Equal(t, 5.5, updatedLoan["rate"])
	assert.Equal(t, 8.0, updatedLoan["roi"])
}

// TestGetLoanTransitions tests getting loan transitions functionality
func TestGetLoanTransitions(t *testing.T) {
	setup := setupTestServer()
	defer setup.server.Close()

	// First create a loan
	createReq := dto.CreateLoanRequest{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
		AgreementLetterLink: "https://example.com/agreement/user123.pdf",
	}

	createResp, err := makeRequest("POST", setup.server.URL+"/api/v1/loans/", createReq)
	require.NoError(t, err)
	defer createResp.Body.Close()

	var createResponse dto.SuccessResponse
	err = json.NewDecoder(createResp.Body).Decode(&createResponse)
	require.NoError(t, err)

	loanData := createResponse.Data.(map[string]interface{})
	loanID := loanData["id"].(string)

	// Get transitions
	resp, err := makeRequest("GET", setup.server.URL+"/api/v1/loans/"+loanID+"/transitions", nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response dto.SuccessResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	transitions := response.Data.(map[string]interface{})
	assert.Equal(t, "proposed", transitions["current_state"])
	assert.NotEmpty(t, transitions["transitions"])
}

// TestApproveLoan tests loan approval functionality
func TestApproveLoan(t *testing.T) {
	setup := setupTestServer()
	defer setup.server.Close()

	// First create a loan
	createReq := dto.CreateLoanRequest{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
		AgreementLetterLink: "https://example.com/agreement/user123.pdf",
	}

	createResp, err := makeRequest("POST", setup.server.URL+"/api/v1/loans/", createReq)
	require.NoError(t, err)
	defer createResp.Body.Close()

	var createResponse dto.SuccessResponse
	err = json.NewDecoder(createResp.Body).Decode(&createResponse)
	require.NoError(t, err)

	loanData := createResponse.Data.(map[string]interface{})
	loanID := loanData["id"].(string)

	// Approve the loan
	approveReq := dto.ApproveLoanRequest{
		FieldValidatorProof: "https://example.com/proofs/field_visit_123.jpg",
		FieldValidatorID:    "validator_001",
	}

	resp, err := makeRequest("PUT", setup.server.URL+"/api/v1/loans/"+loanID+"/approve", approveReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response dto.SuccessResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	approvedLoan := response.Data.(map[string]interface{})
	assert.Equal(t, "approved", approvedLoan["status"])
}

// TestInvestInLoan tests loan investment functionality
func TestInvestInLoan(t *testing.T) {
	setup := setupTestServer()
	defer setup.server.Close()

	// First create and approve a loan
	createReq := dto.CreateLoanRequest{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
		AgreementLetterLink: "https://example.com/agreement/user123.pdf",
	}

	createResp, err := makeRequest("POST", setup.server.URL+"/api/v1/loans/", createReq)
	require.NoError(t, err)
	defer createResp.Body.Close()

	var createResponse dto.SuccessResponse
	err = json.NewDecoder(createResp.Body).Decode(&createResponse)
	require.NoError(t, err)

	loanData := createResponse.Data.(map[string]interface{})
	loanID := loanData["id"].(string)

	// Approve the loan
	approveReq := dto.ApproveLoanRequest{
		FieldValidatorProof: "https://example.com/proofs/field_visit_123.jpg",
		FieldValidatorID:    "validator_001",
	}

	_, err = makeRequest("PUT", setup.server.URL+"/api/v1/loans/"+loanID+"/approve", approveReq)
	require.NoError(t, err)

	// Invest in the loan
	investReq := dto.InvestLoanRequest{
		InvestorID: "investor_001",
		Amount:     15000.00,
	}

	resp, err := makeRequest("PUT", setup.server.URL+"/api/v1/loans/"+loanID+"/invest", investReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response dto.SuccessResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	investedLoan := response.Data.(map[string]interface{})
	assert.Equal(t, 15000.0, investedLoan["total_invested"])
}

// TestDisburseLoan tests loan disbursement functionality
func TestDisburseLoan(t *testing.T) {
	setup := setupTestServer()
	defer setup.server.Close()

	// First create, approve, and fully invest in a loan
	createReq := dto.CreateLoanRequest{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
		AgreementLetterLink: "https://example.com/agreement/user123.pdf",
	}

	createResp, err := makeRequest("POST", setup.server.URL+"/api/v1/loans/", createReq)
	require.NoError(t, err)
	defer createResp.Body.Close()

	var createResponse dto.SuccessResponse
	err = json.NewDecoder(createResp.Body).Decode(&createResponse)
	require.NoError(t, err)

	loanData := createResponse.Data.(map[string]interface{})
	loanID := loanData["id"].(string)

	// Approve the loan
	approveReq := dto.ApproveLoanRequest{
		FieldValidatorProof: "https://example.com/proofs/field_visit_123.jpg",
		FieldValidatorID:    "validator_001",
	}

	_, err = makeRequest("PUT", setup.server.URL+"/api/v1/loans/"+loanID+"/approve", approveReq)
	require.NoError(t, err)

	// Fully invest in the loan
	investReq1 := dto.InvestLoanRequest{
		InvestorID: "investor_001",
		Amount:     15000.00,
	}

	investReq2 := dto.InvestLoanRequest{
		InvestorID: "investor_002",
		Amount:     10000.00,
	}

	_, err = makeRequest("PUT", setup.server.URL+"/api/v1/loans/"+loanID+"/invest", investReq1)
	require.NoError(t, err)

	_, err = makeRequest("PUT", setup.server.URL+"/api/v1/loans/"+loanID+"/invest", investReq2)
	require.NoError(t, err)

	// Disburse the loan
	disburseReq := dto.DisburseLoanRequest{
		SignedAgreementLink: "https://example.com/agreements/signed_agreement_123.pdf",
		FieldOfficerID:      "officer_001",
	}

	resp, err := makeRequest("PUT", setup.server.URL+"/api/v1/loans/"+loanID+"/disburse", disburseReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response dto.SuccessResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	disbursedLoan := response.Data.(map[string]interface{})
	assert.Equal(t, "disbursed", disbursedLoan["status"])
}

// TestGetLoansByStatus tests filtering loans by status functionality
func TestGetLoansByStatus(t *testing.T) {
	setup := setupTestServer()
	defer setup.server.Close()

	// Create multiple loans with different statuses
	createReq1 := dto.CreateLoanRequest{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
		AgreementLetterLink: "https://example.com/agreement/user123.pdf",
	}

	createReq2 := dto.CreateLoanRequest{
		BorrowerID:          "user456",
		PrincipalAmount:     30000.00,
		Rate:                5.0,
		ROI:                 7.0,
		AgreementLetterLink: "https://example.com/agreement/user456.pdf",
	}

	_, err := makeRequest("POST", setup.server.URL+"/api/v1/loans/", createReq1)
	require.NoError(t, err)

	_, err = makeRequest("POST", setup.server.URL+"/api/v1/loans/", createReq2)
	require.NoError(t, err)

	// Get loans by status
	resp, err := makeRequest("GET", setup.server.URL+"/api/v1/loans/?status=proposed", nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response dto.SuccessResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	loans := response.Data.([]interface{})
	assert.GreaterOrEqual(t, len(loans), 2)
}

// TestGetLoansByBorrower tests filtering loans by borrower functionality
func TestGetLoansByBorrower(t *testing.T) {
	setup := setupTestServer()
	defer setup.server.Close()

	// Create loans for different borrowers
	createReq1 := dto.CreateLoanRequest{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
		AgreementLetterLink: "https://example.com/agreement/user123.pdf",
	}

	createReq2 := dto.CreateLoanRequest{
		BorrowerID:          "user456",
		PrincipalAmount:     30000.00,
		Rate:                5.0,
		ROI:                 7.0,
		AgreementLetterLink: "https://example.com/agreement/user456.pdf",
	}

	_, err := makeRequest("POST", setup.server.URL+"/api/v1/loans/", createReq1)
	require.NoError(t, err)

	_, err = makeRequest("POST", setup.server.URL+"/api/v1/loans/", createReq2)
	require.NoError(t, err)

	// Get loans by borrower
	resp, err := makeRequest("GET", setup.server.URL+"/api/v1/loans/?borrower_id=user123", nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response dto.SuccessResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	loans := response.Data.([]interface{})
	assert.GreaterOrEqual(t, len(loans), 1)

	// Verify all returned loans belong to user123
	for _, loan := range loans {
		loanData := loan.(map[string]interface{})
		assert.Equal(t, "user123", loanData["borrower_id"])
	}
}
