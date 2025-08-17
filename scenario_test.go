package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"loan-service/internal/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== SCENARIO TESTS ====================

// TestLoanLifecycle tests the complete loan lifecycle scenario
func TestLoanLifecycle(t *testing.T) {
	fmt.Println("=== Testing Complete Loan Lifecycle Scenario ===")

	setup := setupTestServer()
	defer setup.server.Close()

	baseURL := setup.server.URL

	// 1. Create a loan
	fmt.Println("1. Creating a new loan...")
	createReq := dto.CreateLoanRequest{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
		AgreementLetterLink: "https://example.com/agreement/user123.pdf",
	}

	createResp, err := makeRequest("POST", baseURL+"/api/v1/loans/", createReq)
	require.NoError(t, err)
	defer createResp.Body.Close()

	assert.Equal(t, http.StatusCreated, createResp.StatusCode)

	var createResponse dto.SuccessResponse
	err = json.NewDecoder(createResp.Body).Decode(&createResponse)
	require.NoError(t, err)

	loanData := createResponse.Data.(map[string]interface{})
	loanID := loanData["id"].(string)

	fmt.Printf("Created loan ID: %s\n", loanID)
	assert.Equal(t, "proposed", loanData["status"])

	// 2. Get transitions for proposed loan
	fmt.Println("2. Getting valid transitions for proposed loan...")
	transitionsResp, err := makeRequest("GET", baseURL+"/api/v1/loans/"+loanID+"/transitions", nil)
	require.NoError(t, err)
	defer transitionsResp.Body.Close()

	assert.Equal(t, http.StatusOK, transitionsResp.StatusCode)

	var transitionsResponse dto.SuccessResponse
	err = json.NewDecoder(transitionsResp.Body).Decode(&transitionsResponse)
	require.NoError(t, err)

	transitions := transitionsResponse.Data.(map[string]interface{})
	assert.Equal(t, "proposed", transitions["current_state"])

	// 3. Approve the loan
	fmt.Println("3. Approving the loan...")
	approveReq := dto.ApproveLoanRequest{
		FieldValidatorProof: "https://example.com/proofs/field_visit_123.jpg",
		FieldValidatorID:    "validator_001",
	}

	approveResp, err := makeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/approve", approveReq)
	require.NoError(t, err)
	defer approveResp.Body.Close()

	assert.Equal(t, http.StatusOK, approveResp.StatusCode)

	var approveResponse dto.SuccessResponse
	err = json.NewDecoder(approveResp.Body).Decode(&approveResponse)
	require.NoError(t, err)

	approvedLoan := approveResponse.Data.(map[string]interface{})
	assert.Equal(t, "approved", approvedLoan["status"])

	// 4. Get transitions for approved loan
	fmt.Println("4. Getting valid transitions for approved loan...")
	transitionsResp2, err := makeRequest("GET", baseURL+"/api/v1/loans/"+loanID+"/transitions", nil)
	require.NoError(t, err)
	defer transitionsResp2.Body.Close()

	assert.Equal(t, http.StatusOK, transitionsResp2.StatusCode)

	var transitionsResponse2 dto.SuccessResponse
	err = json.NewDecoder(transitionsResp2.Body).Decode(&transitionsResponse2)
	require.NoError(t, err)

	transitions2 := transitionsResponse2.Data.(map[string]interface{})
	assert.Equal(t, "approved", transitions2["current_state"])

	// 5. Add first investment
	fmt.Println("5. Adding first investment...")
	investReq1 := dto.InvestLoanRequest{
		InvestorID: "investor_001",
		Amount:     15000.00,
	}

	investResp1, err := makeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq1)
	require.NoError(t, err)
	defer investResp1.Body.Close()

	assert.Equal(t, http.StatusOK, investResp1.StatusCode)

	var investResponse1 dto.SuccessResponse
	err = json.NewDecoder(investResp1.Body).Decode(&investResponse1)
	require.NoError(t, err)

	investedLoan1 := investResponse1.Data.(map[string]interface{})
	assert.Equal(t, 15000.0, investedLoan1["total_invested"])
	assert.Equal(t, "approved", investedLoan1["status"])

	// 6. Add second investment to complete the loan
	fmt.Println("6. Adding second investment to complete the loan...")
	investReq2 := dto.InvestLoanRequest{
		InvestorID: "investor_002",
		Amount:     10000.00,
	}

	investResp2, err := makeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq2)
	require.NoError(t, err)
	defer investResp2.Body.Close()

	assert.Equal(t, http.StatusOK, investResp2.StatusCode)

	var investResponse2 dto.SuccessResponse
	err = json.NewDecoder(investResp2.Body).Decode(&investResponse2)
	require.NoError(t, err)

	investedLoan2 := investResponse2.Data.(map[string]interface{})
	assert.Equal(t, 25000.0, investedLoan2["total_invested"])
	assert.Equal(t, "invested", investedLoan2["status"])

	// 7. Get transitions for invested loan
	fmt.Println("7. Getting valid transitions for invested loan...")
	transitionsResp3, err := makeRequest("GET", baseURL+"/api/v1/loans/"+loanID+"/transitions", nil)
	require.NoError(t, err)
	defer transitionsResp3.Body.Close()

	assert.Equal(t, http.StatusOK, transitionsResp3.StatusCode)

	var transitionsResponse3 dto.SuccessResponse
	err = json.NewDecoder(transitionsResp3.Body).Decode(&transitionsResponse3)
	require.NoError(t, err)

	transitions3 := transitionsResponse3.Data.(map[string]interface{})
	assert.Equal(t, "invested", transitions3["current_state"])

	// 8. Disburse the loan
	fmt.Println("8. Disbursing the loan...")
	disburseReq := dto.DisburseLoanRequest{
		SignedAgreementLink: "https://example.com/agreements/signed_agreement_123.pdf",
		FieldOfficerID:      "officer_001",
	}

	disburseResp, err := makeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/disburse", disburseReq)
	require.NoError(t, err)
	defer disburseResp.Body.Close()

	assert.Equal(t, http.StatusOK, disburseResp.StatusCode)

	var disburseResponse dto.SuccessResponse
	err = json.NewDecoder(disburseResp.Body).Decode(&disburseResponse)
	require.NoError(t, err)

	disbursedLoan := disburseResponse.Data.(map[string]interface{})
	assert.Equal(t, "disbursed", disbursedLoan["status"])

	// 9. Get transitions for disbursed loan
	fmt.Println("9. Getting valid transitions for disbursed loan...")
	transitionsResp4, err := makeRequest("GET", baseURL+"/api/v1/loans/"+loanID+"/transitions", nil)
	require.NoError(t, err)
	defer transitionsResp4.Body.Close()

	assert.Equal(t, http.StatusOK, transitionsResp4.StatusCode)

	var transitionsResponse4 dto.SuccessResponse
	err = json.NewDecoder(transitionsResp4.Body).Decode(&transitionsResponse4)
	require.NoError(t, err)

	transitions4 := transitionsResponse4.Data.(map[string]interface{})
	assert.Equal(t, "disbursed", transitions4["current_state"])

	fmt.Println("=== Loan Lifecycle Scenario Completed Successfully ===")
}

// TestInvalidTransitions tests invalid state transition scenarios
func TestInvalidTransitions(t *testing.T) {
	fmt.Println("=== Testing Invalid State Transition Scenarios ===")

	setup := setupTestServer()
	defer setup.server.Close()

	baseURL := setup.server.URL

	// Create and complete a loan lifecycle
	createReq := dto.CreateLoanRequest{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
		AgreementLetterLink: "https://example.com/agreement/user123.pdf",
	}

	createResp, err := makeRequest("POST", baseURL+"/api/v1/loans/", createReq)
	require.NoError(t, err)
	defer createResp.Body.Close()

	var createResponse dto.SuccessResponse
	err = json.NewDecoder(createResp.Body).Decode(&createResponse)
	require.NoError(t, err)

	loanData := createResponse.Data.(map[string]interface{})
	loanID := loanData["id"].(string)

	// Complete the loan lifecycle
	approveReq := dto.ApproveLoanRequest{
		FieldValidatorProof: "https://example.com/proofs/field_visit_123.jpg",
		FieldValidatorID:    "validator_001",
	}

	_, err = makeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/approve", approveReq)
	require.NoError(t, err)

	investReq1 := dto.InvestLoanRequest{
		InvestorID: "investor_001",
		Amount:     15000.00,
	}

	investReq2 := dto.InvestLoanRequest{
		InvestorID: "investor_002",
		Amount:     10000.00,
	}

	_, err = makeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq1)
	require.NoError(t, err)

	_, err = makeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq2)
	require.NoError(t, err)

	disburseReq := dto.DisburseLoanRequest{
		SignedAgreementLink: "https://example.com/agreements/signed_agreement_123.pdf",
		FieldOfficerID:      "officer_001",
	}

	_, err = makeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/disburse", disburseReq)
	require.NoError(t, err)

	// Test 1: Try to approve a disbursed loan
	fmt.Println("1. Testing invalid transition: trying to approve a disbursed loan...")

	invalidApproveReq := dto.ApproveLoanRequest{
		FieldValidatorProof: "https://example.com/proofs/field_visit_456.jpg",
		FieldValidatorID:    "validator_002",
	}

	resp1, err := makeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/approve", invalidApproveReq)
	require.NoError(t, err)
	defer resp1.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp1.StatusCode)

	var errorResponse1 dto.ErrorResponse
	err = json.NewDecoder(resp1.Body).Decode(&errorResponse1)
	require.NoError(t, err)

	assert.Contains(t, errorResponse1.Message, "can only approve loans in proposed status")

	// Test 2: Try to invest in a disbursed loan
	fmt.Println("2. Testing invalid transition: trying to invest in a disbursed loan...")

	invalidInvestReq := dto.InvestLoanRequest{
		InvestorID: "investor_003",
		Amount:     5000.00,
	}

	resp2, err := makeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", invalidInvestReq)
	require.NoError(t, err)
	defer resp2.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp2.StatusCode)

	var errorResponse2 dto.ErrorResponse
	err = json.NewDecoder(resp2.Body).Decode(&errorResponse2)
	require.NoError(t, err)

	assert.Contains(t, errorResponse2.Message, "loan is not in approved status")

	// Test 3: Try to disburse an already disbursed loan
	fmt.Println("3. Testing invalid transition: trying to disburse an already disbursed loan...")

	invalidDisburseReq := dto.DisburseLoanRequest{
		SignedAgreementLink: "https://example.com/agreements/signed_agreement_456.pdf",
		FieldOfficerID:      "officer_002",
	}

	resp3, err := makeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/disburse", invalidDisburseReq)
	require.NoError(t, err)
	defer resp3.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp3.StatusCode)

	var errorResponse3 dto.ErrorResponse
	err = json.NewDecoder(resp3.Body).Decode(&errorResponse3)
	require.NoError(t, err)

	assert.Contains(t, errorResponse3.Message, "can only disburse fully invested loans")

	fmt.Println("=== Invalid State Transition Scenarios Completed Successfully ===")
}

// TestInvalidInvestment tests invalid investment scenarios
func TestInvalidInvestment(t *testing.T) {
	fmt.Println("=== Testing Invalid Investment Scenarios ===")

	setup := setupTestServer()
	defer setup.server.Close()

	baseURL := setup.server.URL

	// Create and approve a loan
	createReq := dto.CreateLoanRequest{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
		AgreementLetterLink: "https://example.com/agreement/user123.pdf",
	}

	createResp, err := makeRequest("POST", baseURL+"/api/v1/loans/", createReq)
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

	_, err = makeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/approve", approveReq)
	require.NoError(t, err)

	// Test 1: Try to invest more than principal amount
	fmt.Println("1. Testing invalid investment: trying to invest more than principal amount...")

	investReq1 := dto.InvestLoanRequest{
		InvestorID: "investor_001",
		Amount:     30000.00, // More than the 25000 principal
	}

	resp1, err := makeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq1)
	require.NoError(t, err)
	defer resp1.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp1.StatusCode)

	var errorResponse1 dto.ErrorResponse
	err = json.NewDecoder(resp1.Body).Decode(&errorResponse1)
	require.NoError(t, err)

	assert.Contains(t, errorResponse1.Message, "total investment amount would exceed loan principal")

	// Test 2: Try to invest negative amount
	fmt.Println("2. Testing invalid investment: trying to invest negative amount...")

	investReq2 := dto.InvestLoanRequest{
		InvestorID: "investor_002",
		Amount:     -1000.00,
	}

	resp2, err := makeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq2)
	require.NoError(t, err)
	defer resp2.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp2.StatusCode)

	var errorResponse2 dto.ErrorResponse
	err = json.NewDecoder(resp2.Body).Decode(&errorResponse2)
	require.NoError(t, err)

	assert.Contains(t, errorResponse2.Message, "Field validation for 'Amount' failed on the 'gt' tag")

	// Test 3: Try to invest zero amount
	fmt.Println("3. Testing invalid investment: trying to invest zero amount...")

	investReq3 := dto.InvestLoanRequest{
		InvestorID: "investor_003",
		Amount:     0.00,
	}

	resp3, err := makeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq3)
	require.NoError(t, err)
	defer resp3.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp3.StatusCode)

	var errorResponse3 dto.ErrorResponse
	err = json.NewDecoder(resp3.Body).Decode(&errorResponse3)
	require.NoError(t, err)

	// Fix: Use the correct error message for zero amount
	assert.Contains(t, errorResponse3.Message, "Field validation for 'Amount' failed on the 'required' tag")

	fmt.Println("=== Invalid Investment Scenarios Completed Successfully ===")
}

// TestInvalidDisbursement tests invalid disbursement scenarios
func TestInvalidDisbursement(t *testing.T) {
	fmt.Println("=== Testing Invalid Disbursement Scenarios ===")

	setup := setupTestServer()
	defer setup.server.Close()

	baseURL := setup.server.URL

	// Create and approve a loan
	createReq := dto.CreateLoanRequest{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
		AgreementLetterLink: "https://example.com/agreement/user123.pdf",
	}

	createResp, err := makeRequest("POST", baseURL+"/api/v1/loans/", createReq)
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

	_, err = makeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/approve", approveReq)
	require.NoError(t, err)

	// Test 1: Try to disburse without full investment
	fmt.Println("1. Testing invalid disbursement: trying to disburse without full investment...")

	disburseReq1 := dto.DisburseLoanRequest{
		SignedAgreementLink: "https://example.com/agreements/signed_agreement_123.pdf",
		FieldOfficerID:      "officer_001",
	}

	resp1, err := makeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/disburse", disburseReq1)
	require.NoError(t, err)
	defer resp1.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp1.StatusCode)

	var errorResponse1 dto.ErrorResponse
	err = json.NewDecoder(resp1.Body).Decode(&errorResponse1)
	require.NoError(t, err)

	assert.Contains(t, errorResponse1.Message, "can only disburse fully invested loans")

	// Test 2: Try to disburse with partial investment
	fmt.Println("2. Testing invalid disbursement: trying to disburse with partial investment...")

	// Add partial investment first
	investReq := dto.InvestLoanRequest{
		InvestorID: "investor_001",
		Amount:     15000.00,
	}

	_, err = makeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq)
	require.NoError(t, err)

	disburseReq2 := dto.DisburseLoanRequest{
		SignedAgreementLink: "https://example.com/agreements/signed_agreement_456.pdf",
		FieldOfficerID:      "officer_002",
	}

	resp2, err := makeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/disburse", disburseReq2)
	require.NoError(t, err)
	defer resp2.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp2.StatusCode)

	var errorResponse2 dto.ErrorResponse
	err = json.NewDecoder(resp2.Body).Decode(&errorResponse2)
	require.NoError(t, err)

	assert.Contains(t, errorResponse2.Message, "can only disburse fully invested loans")

	fmt.Println("=== Invalid Disbursement Scenarios Completed Successfully ===")
}
