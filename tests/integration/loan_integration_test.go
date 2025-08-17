package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"loan-service/internal/dto"
	"loan-service/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoanLifecycleIntegration(t *testing.T) {
	setup := testutils.SetupTestServer()
	defer setup.Server.Close()

	baseURL := setup.Server.URL

	// 1. Create a loan
	t.Log("1. Creating a new loan...")
	createReq := dto.CreateLoanRequest{
		BorrowerID:      "user123",
		PrincipalAmount: 25000.00,
		Rate:            4.5,
		ROI:             6.0,
	}

	createResp, err := testutils.MakeRequest("POST", baseURL+"/api/v1/loans/", createReq)
	require.NoError(t, err)
	defer createResp.Body.Close()

	assert.Equal(t, http.StatusCreated, createResp.StatusCode)

	var createResponse dto.SuccessResponse
	err = json.NewDecoder(createResp.Body).Decode(&createResponse)
	require.NoError(t, err)

	loanData := createResponse.Data.(map[string]interface{})
	loanID := loanData["id"].(string)

	t.Logf("Created loan ID: %s", loanID)
	assert.Equal(t, "proposed", loanData["status"])

	// 2. Approve the loan
	t.Log("2. Approving the loan...")
	approveReq := dto.ApproveLoanRequest{
		FieldValidatorProof: "https://example.com/proofs/field_visit_123.jpg",
		FieldValidatorID:    "validator_001",
	}

	approveResp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/approve", approveReq)
	require.NoError(t, err)
	defer approveResp.Body.Close()

	assert.Equal(t, http.StatusOK, approveResp.StatusCode)

	var approveResponse dto.SuccessResponse
	err = json.NewDecoder(approveResp.Body).Decode(&approveResponse)
	require.NoError(t, err)

	approvedLoan := approveResponse.Data.(map[string]interface{})
	assert.Equal(t, "approved", approvedLoan["status"])

	// 3. Invest in the loan
	t.Log("3. Investing in the loan...")
	investReq := dto.InvestLoanRequest{
		InvestorID:          "investor_001",
		Amount:              25000.00,
		AgreementLetterLink: "https://example.com/agreement.pdf",
	}

	investResp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq)
	require.NoError(t, err)
	defer investResp.Body.Close()

	assert.Equal(t, http.StatusOK, investResp.StatusCode)

	var investResponse dto.SuccessResponse
	err = json.NewDecoder(investResp.Body).Decode(&investResponse)
	require.NoError(t, err)

	investedLoan := investResponse.Data.(map[string]interface{})
	assert.Equal(t, "invested", investedLoan["status"])
	assert.Equal(t, 25000.0, investedLoan["total_invested"])

	// 4. Disburse the loan
	t.Log("4. Disbursing the loan...")
	disburseReq := dto.DisburseLoanRequest{
		SignedAgreementLink: "https://example.com/signed-agreement.pdf",
		FieldOfficerID:      "officer_001",
	}

	disburseResp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/disburse", disburseReq)
	require.NoError(t, err)
	defer disburseResp.Body.Close()

	assert.Equal(t, http.StatusOK, disburseResp.StatusCode)

	var disburseResponse dto.SuccessResponse
	err = json.NewDecoder(disburseResp.Body).Decode(&disburseResponse)
	require.NoError(t, err)

	disbursedLoan := disburseResponse.Data.(map[string]interface{})
	assert.Equal(t, "disbursed", disbursedLoan["status"])
}

func TestHealthEndpointIntegration(t *testing.T) {
	setup := testutils.SetupTestServer()
	defer setup.Server.Close()

	resp, err := testutils.MakeRequest("GET", setup.Server.URL+"/health", nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "loan-service", response["service"])
}

// ========== FAILURE TEST CASES ==========

func TestCreateLoanFailureCases(t *testing.T) {
	setup := testutils.SetupTestServer()
	defer setup.Server.Close()

	baseURL := setup.Server.URL

	t.Run("Create loan with missing required fields", func(t *testing.T) {
		// Missing borrower_id
		createReq := dto.CreateLoanRequest{
			PrincipalAmount: 25000.00,
			Rate:            4.5,
			ROI:             6.0,
		}

		resp, err := testutils.MakeRequest("POST", baseURL+"/api/v1/loans/", createReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create loan with invalid principal amount", func(t *testing.T) {
		createReq := dto.CreateLoanRequest{
			BorrowerID:      "user123",
			PrincipalAmount: 0, // Invalid: must be greater than 0
			Rate:            4.5,
			ROI:             6.0,
		}

		resp, err := testutils.MakeRequest("POST", baseURL+"/api/v1/loans/", createReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create loan with negative rate", func(t *testing.T) {
		createReq := dto.CreateLoanRequest{
			BorrowerID:      "user123",
			PrincipalAmount: 25000.00,
			Rate:            -1.0, // Invalid: must be greater than 0
			ROI:             6.0,
		}

		resp, err := testutils.MakeRequest("POST", baseURL+"/api/v1/loans/", createReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestApproveLoanFailureCases(t *testing.T) {
	setup := testutils.SetupTestServer()
	defer setup.Server.Close()

	baseURL := setup.Server.URL

	// Create a loan first
	createReq := dto.CreateLoanRequest{
		BorrowerID:      "user123",
		PrincipalAmount: 25000.00,
		Rate:            4.5,
		ROI:             6.0,
	}

	createResp, err := testutils.MakeRequest("POST", baseURL+"/api/v1/loans/", createReq)
	require.NoError(t, err)
	defer createResp.Body.Close()

	var createResponse dto.SuccessResponse
	err = json.NewDecoder(createResp.Body).Decode(&createResponse)
	require.NoError(t, err)

	loanData := createResponse.Data.(map[string]interface{})
	loanID := loanData["id"].(string)

	t.Run("Approve loan with missing field validator proof", func(t *testing.T) {
		approveReq := dto.ApproveLoanRequest{
			FieldValidatorID: "validator_001",
			// Missing FieldValidatorProof
		}

		resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/approve", approveReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Approve loan with missing field validator ID", func(t *testing.T) {
		approveReq := dto.ApproveLoanRequest{
			FieldValidatorProof: "https://example.com/proofs/field_visit_123.jpg",
			// Missing FieldValidatorID
		}

		resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/approve", approveReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Approve non-existent loan", func(t *testing.T) {
		approveReq := dto.ApproveLoanRequest{
			FieldValidatorProof: "https://example.com/proofs/field_visit_123.jpg",
			FieldValidatorID:    "validator_001",
		}

		resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/non-existent-id/approve", approveReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("Approve already approved loan", func(t *testing.T) {
		// First approval
		approveReq := dto.ApproveLoanRequest{
			FieldValidatorProof: "https://example.com/proofs/field_visit_123.jpg",
			FieldValidatorID:    "validator_001",
		}

		resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/approve", approveReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Try to approve again
		resp2, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/approve", approveReq)
		require.NoError(t, err)
		defer resp2.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp2.StatusCode)
	})
}

func TestInvestLoanFailureCases(t *testing.T) {
	setup := testutils.SetupTestServer()
	defer setup.Server.Close()

	baseURL := setup.Server.URL

	// Create and approve a loan first
	createReq := dto.CreateLoanRequest{
		BorrowerID:      "user123",
		PrincipalAmount: 25000.00,
		Rate:            4.5,
		ROI:             6.0,
	}

	createResp, err := testutils.MakeRequest("POST", baseURL+"/api/v1/loans/", createReq)
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

	approveResp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/approve", approveReq)
	require.NoError(t, err)
	defer approveResp.Body.Close()

	t.Run("Invest in loan with missing investor ID", func(t *testing.T) {
		investReq := dto.InvestLoanRequest{
			Amount:              25000.00,
			// Missing InvestorID
		}

		resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Invest in loan with invalid amount", func(t *testing.T) {
		investReq := dto.InvestLoanRequest{
			InvestorID:          "investor_001",
			Amount:              0, // Invalid: must be greater than 0
		}

		resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Invest in non-existent loan", func(t *testing.T) {
		investReq := dto.InvestLoanRequest{
			InvestorID: "investor_001",
			Amount:     25000.00,
		}

		resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/non-existent-id/invest", investReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Invest in proposed loan (not approved)", func(t *testing.T) {
		// Create another loan without approving
		createReq2 := dto.CreateLoanRequest{
			BorrowerID:          "user456",
			PrincipalAmount:     10000.00,
			Rate:                4.5,
			ROI:                 6.0,
			
		}

		createResp2, err := testutils.MakeRequest("POST", baseURL+"/api/v1/loans/", createReq2)
		require.NoError(t, err)
		defer createResp2.Body.Close()

		var createResponse2 dto.SuccessResponse
		err = json.NewDecoder(createResp2.Body).Decode(&createResponse2)
		require.NoError(t, err)

		loanData2 := createResponse2.Data.(map[string]interface{})
		loanID2 := loanData2["id"].(string)

		// Try to invest in proposed loan
		investReq := dto.InvestLoanRequest{
			InvestorID: "investor_001",
			Amount:     10000.00,
		}

		resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID2+"/invest", investReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Invest amount exceeding loan principal", func(t *testing.T) {
		investReq := dto.InvestLoanRequest{
			InvestorID: "investor_001",
			Amount:     30000.00, // Exceeds principal amount of 25000
		}

		resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestDisburseLoanFailureCases(t *testing.T) {
	setup := testutils.SetupTestServer()
	defer setup.Server.Close()

	baseURL := setup.Server.URL

	// Create, approve, and invest in a loan first
	createReq := dto.CreateLoanRequest{
		BorrowerID:          "user123",
		PrincipalAmount:     25000.00,
		Rate:                4.5,
		ROI:                 6.0,
		
	}

	createResp, err := testutils.MakeRequest("POST", baseURL+"/api/v1/loans/", createReq)
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

	approveResp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/approve", approveReq)
	require.NoError(t, err)
	defer approveResp.Body.Close()

	// Invest in the loan
	investReq := dto.InvestLoanRequest{
		InvestorID: "investor_001",
		Amount:     25000.00,
	}

	investResp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq)
	require.NoError(t, err)
	defer investResp.Body.Close()

	t.Run("Disburse loan with missing signed agreement link", func(t *testing.T) {
		disburseReq := dto.DisburseLoanRequest{
			FieldOfficerID: "officer_001",
			// Missing SignedAgreementLink
		}

		resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/disburse", disburseReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Disburse loan with missing field officer ID", func(t *testing.T) {
		disburseReq := dto.DisburseLoanRequest{
			SignedAgreementLink: "https://example.com/signed-agreement.pdf",
			// Missing FieldOfficerID
		}

		resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/disburse", disburseReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Disburse non-existent loan", func(t *testing.T) {
		disburseReq := dto.DisburseLoanRequest{
			SignedAgreementLink: "https://example.com/signed-agreement.pdf",
			FieldOfficerID:      "officer_001",
		}

		resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/non-existent-id/disburse", disburseReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("Disburse loan that is not fully invested", func(t *testing.T) {
		// Create another loan and approve it
		createReq2 := dto.CreateLoanRequest{
			BorrowerID:          "user456",
			PrincipalAmount:     10000.00,
			Rate:                4.5,
			ROI:                 6.0,
			
		}

		createResp2, err := testutils.MakeRequest("POST", baseURL+"/api/v1/loans/", createReq2)
		require.NoError(t, err)
		defer createResp2.Body.Close()

		var createResponse2 dto.SuccessResponse
		err = json.NewDecoder(createResp2.Body).Decode(&createResponse2)
		require.NoError(t, err)

		loanData2 := createResponse2.Data.(map[string]interface{})
		loanID2 := loanData2["id"].(string)

		// Approve the loan
		approveReq2 := dto.ApproveLoanRequest{
			FieldValidatorProof: "https://example.com/proofs/field_visit_456.jpg",
			FieldValidatorID:    "validator_002",
		}

		approveResp2, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID2+"/approve", approveReq2)
		require.NoError(t, err)
		defer approveResp2.Body.Close()

		// Invest partially (not fully)
		investReq2 := dto.InvestLoanRequest{
			InvestorID: "investor_001",
			Amount:     5000.00, // Only half of the principal
		}

		investResp2, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID2+"/invest", investReq2)
		require.NoError(t, err)
		defer investResp2.Body.Close()

		// Try to disburse partially invested loan
		disburseReq := dto.DisburseLoanRequest{
			SignedAgreementLink: "https://example.com/signed-agreement.pdf",
			FieldOfficerID:      "officer_001",
		}

		resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID2+"/disburse", disburseReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Disburse already disbursed loan", func(t *testing.T) {
		// Disburse the loan first
		disburseReq := dto.DisburseLoanRequest{
			SignedAgreementLink: "https://example.com/signed-agreement.pdf",
			FieldOfficerID:      "officer_001",
		}

		resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/disburse", disburseReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Try to disburse again
		resp2, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/disburse", disburseReq)
		require.NoError(t, err)
		defer resp2.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp2.StatusCode)
	})
}

func TestStateTransitionFailureCases(t *testing.T) {
	setup := testutils.SetupTestServer()
	defer setup.Server.Close()

	baseURL := setup.Server.URL

	t.Run("Cannot approve already approved loan", func(t *testing.T) {
		// Create a loan
		createReq := dto.CreateLoanRequest{
			BorrowerID:          "user123",
			PrincipalAmount:     25000.00,
			Rate:                4.5,
			ROI:                 6.0,
			
		}

		createResp, err := testutils.MakeRequest("POST", baseURL+"/api/v1/loans/", createReq)
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

		approveResp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/approve", approveReq)
		require.NoError(t, err)
		defer approveResp.Body.Close()

		assert.Equal(t, http.StatusOK, approveResp.StatusCode)

		// Try to approve again
		approveResp2, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/approve", approveReq)
		require.NoError(t, err)
		defer approveResp2.Body.Close()

		assert.Equal(t, http.StatusBadRequest, approveResp2.StatusCode)
	})

	t.Run("Cannot invest in proposed loan", func(t *testing.T) {
		// Create a loan
		createReq := dto.CreateLoanRequest{
			BorrowerID:          "user456",
			PrincipalAmount:     10000.00,
			Rate:                4.5,
			ROI:                 6.0,
			
		}

		createResp, err := testutils.MakeRequest("POST", baseURL+"/api/v1/loans/", createReq)
		require.NoError(t, err)
		defer createResp.Body.Close()

		var createResponse dto.SuccessResponse
		err = json.NewDecoder(createResp.Body).Decode(&createResponse)
		require.NoError(t, err)

		loanData := createResponse.Data.(map[string]interface{})
		loanID := loanData["id"].(string)

		// Try to invest in proposed loan
		investReq := dto.InvestLoanRequest{
			InvestorID: "investor_001",
			Amount:     10000.00,
		}

		investResp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq)
		require.NoError(t, err)
		defer investResp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, investResp.StatusCode)
	})

	t.Run("Cannot disburse approved but not invested loan", func(t *testing.T) {
		// Create a loan
		createReq := dto.CreateLoanRequest{
			BorrowerID:          "user789",
			PrincipalAmount:     10000.00,
			Rate:                4.5,
			ROI:                 6.0,
			
		}

		createResp, err := testutils.MakeRequest("POST", baseURL+"/api/v1/loans/", createReq)
		require.NoError(t, err)
		defer createResp.Body.Close()

		var createResponse dto.SuccessResponse
		err = json.NewDecoder(createResp.Body).Decode(&createResponse)
		require.NoError(t, err)

		loanData := createResponse.Data.(map[string]interface{})
		loanID := loanData["id"].(string)

		// Approve the loan
		approveReq := dto.ApproveLoanRequest{
			FieldValidatorProof: "https://example.com/proofs/field_visit_789.jpg",
			FieldValidatorID:    "validator_003",
		}

		approveResp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/approve", approveReq)
		require.NoError(t, err)
		defer approveResp.Body.Close()

		// Try to disburse approved but not invested loan
		disburseReq := dto.DisburseLoanRequest{
			SignedAgreementLink: "https://example.com/signed-agreement.pdf",
			FieldOfficerID:      "officer_001",
		}

		disburseResp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/disburse", disburseReq)
		require.NoError(t, err)
		defer disburseResp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, disburseResp.StatusCode)
	})
}

func TestMultipleInvestmentFailureCases(t *testing.T) {
	setup := testutils.SetupTestServer()
	defer setup.Server.Close()

	baseURL := setup.Server.URL

	// Create and approve a loan
	createReq := dto.CreateLoanRequest{
		BorrowerID:          "user123",
		PrincipalAmount:     10000.00,
		Rate:                4.5,
		ROI:                 6.0,
		
	}

	createResp, err := testutils.MakeRequest("POST", baseURL+"/api/v1/loans/", createReq)
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

	approveResp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/approve", approveReq)
	require.NoError(t, err)
	defer approveResp.Body.Close()

	t.Run("Multiple investments exceeding principal amount", func(t *testing.T) {
		// First investment: 6000
		investReq1 := dto.InvestLoanRequest{
			InvestorID: "investor_001",
			Amount:     6000.00,
		}

		resp1, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq1)
		require.NoError(t, err)
		defer resp1.Body.Close()

		assert.Equal(t, http.StatusOK, resp1.StatusCode)

		// Second investment: 5000 (total would be 11000, exceeding 10000)
		investReq2 := dto.InvestLoanRequest{
			InvestorID: "investor_002",
			Amount:     5000.00,
		}

		resp2, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq2)
		require.NoError(t, err)
		defer resp2.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp2.StatusCode)
	})

	t.Run("Valid multiple investments", func(t *testing.T) {
		// Create another loan for this test
		createReq2 := dto.CreateLoanRequest{
			BorrowerID:          "user456",
			PrincipalAmount:     10000.00,
			Rate:                4.5,
			ROI:                 6.0,
			
		}

		createResp2, err := testutils.MakeRequest("POST", baseURL+"/api/v1/loans/", createReq2)
		require.NoError(t, err)
		defer createResp2.Body.Close()

		var createResponse2 dto.SuccessResponse
		err = json.NewDecoder(createResp2.Body).Decode(&createResponse2)
		require.NoError(t, err)

		loanData2 := createResponse2.Data.(map[string]interface{})
		loanID2 := loanData2["id"].(string)

		// Approve the loan
		approveReq2 := dto.ApproveLoanRequest{
			FieldValidatorProof: "https://example.com/proofs/field_visit_456.jpg",
			FieldValidatorID:    "validator_002",
		}

		approveResp2, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID2+"/approve", approveReq2)
		require.NoError(t, err)
		defer approveResp2.Body.Close()

		// First investment: 4000
		investReq1 := dto.InvestLoanRequest{
			InvestorID: "investor_001",
			Amount:     4000.00,
		}

		resp1, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID2+"/invest", investReq1)
		require.NoError(t, err)
		defer resp1.Body.Close()

		assert.Equal(t, http.StatusOK, resp1.StatusCode)

		// Second investment: 6000 (total: 10000, exactly the principal)
		investReq2 := dto.InvestLoanRequest{
			InvestorID:          "investor_002",
			Amount:              6000.00,
			AgreementLetterLink: "https://example.com/agreement.pdf",
		}

		resp2, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID2+"/invest", investReq2)
		require.NoError(t, err)
		defer resp2.Body.Close()

		assert.Equal(t, http.StatusOK, resp2.StatusCode)

		// Verify the loan is now in invested status
		var investResponse dto.SuccessResponse
		err = json.NewDecoder(resp2.Body).Decode(&investResponse)
		require.NoError(t, err)

		// Convert the data to LoanResponse
		investResponseData, err := json.Marshal(investResponse.Data)
		require.NoError(t, err)
		
		var investedLoan dto.LoanResponse
		err = json.Unmarshal(investResponseData, &investedLoan)
		require.NoError(t, err)
		
		assert.Equal(t, "invested", string(investedLoan.Status))
		assert.Equal(t, 10000.0, investedLoan.TotalInvested)
	})
}
