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

// ========== HAPPY PATH TESTS ==========

func TestCompleteLoanLifecycle(t *testing.T) {
	setup := testutils.SetupTestServer()
	defer setup.Server.Close()

	baseURL := setup.Server.URL

	t.Log("=== Testing Complete Loan Lifecycle ===")

	// Step 1: Create a new loan application
	t.Log("Step 1: Creating loan application...")
	createReq := dto.CreateLoanRequest{
		BorrowerID:      "borrower_001",
		PrincipalAmount: 50000.00,
		Rate:            5.5,
		ROI:             7.2,
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

	t.Logf("✓ Loan created successfully with ID: %s", loanID)
	assert.Equal(t, "proposed", loanData["status"])
	assert.Equal(t, 50000.0, loanData["principal_amount"])

	// Step 2: Field validation and loan approval
	t.Log("Step 2: Processing field validation and approval...")
	approveReq := dto.ApproveLoanRequest{
		FieldValidatorProof: "https://example.com/field-validation/proof_123.jpg",
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

	// Verify approval details
	approvalDetails := approvedLoan["approval_details"].(map[string]interface{})
	assert.Equal(t, "https://example.com/field-validation/proof_123.jpg", approvalDetails["field_validator_proof"])
	assert.Equal(t, "validator_001", approvalDetails["field_validator_id"])
	assert.NotEmpty(t, approvalDetails["approval_date"])

	t.Log("✓ Loan approved with field validation proof and approval date recorded")

	// Step 3: Investment process
	t.Log("Step 3: Processing investment...")
	investReq := dto.InvestLoanRequest{
		InvestorID: "investor_001",
		Amount:     50000.00,
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
	assert.Equal(t, 50000.0, investedLoan["total_invested"])

	// Verify auto-generated agreement letter link
	assert.NotEmpty(t, investedLoan["agreement_letter_link"])
	agreementLink := investedLoan["agreement_letter_link"].(string)
	assert.Contains(t, agreementLink, "https://example.com/agreements/loan_")
	assert.Contains(t, agreementLink, "_agreement.pdf")

	t.Log("✓ Loan fully invested with auto-generated agreement letter link")

	// Step 4: Loan disbursement
	t.Log("Step 4: Processing loan disbursement...")
	disburseReq := dto.DisburseLoanRequest{
		SignedAgreementLink: "https://example.com/signed-agreements/loan_001_signed.pdf",
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

	// Verify disbursement details
	disbursementDetails := disbursedLoan["disbursement_details"].(map[string]interface{})
	assert.Equal(t, "https://example.com/signed-agreements/loan_001_signed.pdf", disbursementDetails["signed_agreement_link"])
	assert.Equal(t, "officer_001", disbursementDetails["field_officer_id"])
	assert.NotEmpty(t, disbursementDetails["disbursement_date"])

	t.Log("✓ Loan successfully disbursed with signed agreement and disbursement date recorded")
	t.Log("=== Complete loan lifecycle test passed ===")
}

func TestMultipleInvestorScenario(t *testing.T) {
	setup := testutils.SetupTestServer()
	defer setup.Server.Close()

	baseURL := setup.Server.URL

	t.Log("=== Testing Multiple Investor Scenario ===")

	// Step 1: Create loan application
	t.Log("Step 1: Creating loan application...")
	createReq := dto.CreateLoanRequest{
		BorrowerID:      "borrower_002",
		PrincipalAmount: 100000.00,
		Rate:            6.0,
		ROI:             8.5,
	}

	createResp, err := testutils.MakeRequest("POST", baseURL+"/api/v1/loans/", createReq)
	require.NoError(t, err)
	defer createResp.Body.Close()

	var createResponse dto.SuccessResponse
	err = json.NewDecoder(createResp.Body).Decode(&createResponse)
	require.NoError(t, err)

	loanData := createResponse.Data.(map[string]interface{})
	loanID := loanData["id"].(string)

	// Step 2: Approve loan
	t.Log("Step 2: Approving loan...")
	approveReq := dto.ApproveLoanRequest{
		FieldValidatorProof: "https://example.com/field-validation/proof_456.png",
		FieldValidatorID:    "validator_002",
	}

	approveResp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/approve", approveReq)
	require.NoError(t, err)
	defer approveResp.Body.Close()

	// Step 3: Multiple investments
	t.Log("Step 3: Processing multiple investments...")

	// First investor: 40% of loan
	investReq1 := dto.InvestLoanRequest{
		InvestorID: "investor_001",
		Amount:     40000.00,
	}

	investResp1, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq1)
	require.NoError(t, err)
	defer investResp1.Body.Close()

	assert.Equal(t, http.StatusOK, investResp1.StatusCode)

	var investResponse1 dto.SuccessResponse
	err = json.NewDecoder(investResp1.Body).Decode(&investResponse1)
	require.NoError(t, err)

	partialLoan := investResponse1.Data.(map[string]interface{})
	assert.Equal(t, "approved", partialLoan["status"]) // Still approved, not fully invested
	assert.Equal(t, 40000.0, partialLoan["total_invested"])

	t.Log("✓ First investment completed (40% of loan)")

	// Second investor: 35% of loan
	investReq2 := dto.InvestLoanRequest{
		InvestorID: "investor_002",
		Amount:     35000.00,
	}

	investResp2, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq2)
	require.NoError(t, err)
	defer investResp2.Body.Close()

	assert.Equal(t, http.StatusOK, investResp2.StatusCode)

	var investResponse2 dto.SuccessResponse
	err = json.NewDecoder(investResp2.Body).Decode(&investResponse2)
	require.NoError(t, err)

	partialLoan2 := investResponse2.Data.(map[string]interface{})
	assert.Equal(t, "approved", partialLoan2["status"]) // Still approved, not fully invested
	assert.Equal(t, 75000.0, partialLoan2["total_invested"])

	t.Log("✓ Second investment completed (75% of loan total)")

	// Third investor: 25% of loan (completes the investment)
	investReq3 := dto.InvestLoanRequest{
		InvestorID: "investor_003",
		Amount:     25000.00,
	}

	investResp3, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq3)
	require.NoError(t, err)
	defer investResp3.Body.Close()

	assert.Equal(t, http.StatusOK, investResp3.StatusCode)

	var investResponse3 dto.SuccessResponse
	err = json.NewDecoder(investResp3.Body).Decode(&investResponse3)
	require.NoError(t, err)

	fullyInvestedLoan := investResponse3.Data.(map[string]interface{})
	assert.Equal(t, "invested", fullyInvestedLoan["status"]) // Now fully invested
	assert.Equal(t, 100000.0, fullyInvestedLoan["total_invested"])

	// Verify auto-generated agreement letter link
	assert.NotEmpty(t, fullyInvestedLoan["agreement_letter_link"])

	t.Log("✓ Third investment completed - loan now fully invested with auto-generated agreement")
	t.Log("=== Multiple investor scenario test passed ===")
}

func TestHealthCheckEndpoint(t *testing.T) {
	setup := testutils.SetupTestServer()
	defer setup.Server.Close()

	t.Log("=== Testing Health Check Endpoint ===")

	resp, err := testutils.MakeRequest("GET", setup.Server.URL+"/health", nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "loan-service", response["service"])

	t.Log("✓ Health check endpoint responding correctly")
}

// ========== VALIDATION TESTS ==========

func TestInputValidationScenarios(t *testing.T) {
	setup := testutils.SetupTestServer()
	defer setup.Server.Close()

	baseURL := setup.Server.URL

	t.Log("=== Testing Input Validation Scenarios ===")

	t.Run("Loan Creation Validation", func(t *testing.T) {
		t.Log("Testing loan creation validation...")

		t.Run("Missing required fields", func(t *testing.T) {
			createReq := dto.CreateLoanRequest{
				PrincipalAmount: 25000.00,
				Rate:            4.5,
				ROI:             6.0,
				// Missing BorrowerID
			}

			resp, err := testutils.MakeRequest("POST", baseURL+"/api/v1/loans/", createReq)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			t.Log("✓ Correctly rejected loan creation with missing borrower ID")
		})

		t.Run("Invalid principal amount", func(t *testing.T) {
			createReq := dto.CreateLoanRequest{
				BorrowerID:      "borrower_001",
				PrincipalAmount: 0, // Invalid: must be greater than 0
				Rate:            4.5,
				ROI:             6.0,
			}

			resp, err := testutils.MakeRequest("POST", baseURL+"/api/v1/loans/", createReq)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			t.Log("✓ Correctly rejected loan creation with invalid principal amount")
		})

		t.Run("Negative interest rate", func(t *testing.T) {
			createReq := dto.CreateLoanRequest{
				BorrowerID:      "borrower_001",
				PrincipalAmount: 25000.00,
				Rate:            -1.0, // Invalid: must be greater than 0
				ROI:             6.0,
			}

			resp, err := testutils.MakeRequest("POST", baseURL+"/api/v1/loans/", createReq)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			t.Log("✓ Correctly rejected loan creation with negative interest rate")
		})
	})

	t.Run("Loan Approval Validation", func(t *testing.T) {
		t.Log("Testing loan approval validation...")

		// Create a loan for approval tests
		createReq := dto.CreateLoanRequest{
			BorrowerID:      "borrower_001",
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

		t.Run("Missing field validator proof", func(t *testing.T) {
			approveReq := dto.ApproveLoanRequest{
				FieldValidatorID: "validator_001",
				// Missing FieldValidatorProof
			}

			resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/approve", approveReq)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			t.Log("✓ Correctly rejected approval with missing field validator proof")
		})

		t.Run("Invalid image link format", func(t *testing.T) {
			approveReq := dto.ApproveLoanRequest{
				FieldValidatorProof: "https://example.com/document.pdf", // Invalid: not an image
				FieldValidatorID:    "validator_001",
			}

			resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/approve", approveReq)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			var errorResponse dto.ErrorResponse
			err = json.NewDecoder(resp.Body).Decode(&errorResponse)
			require.NoError(t, err)
			assert.Contains(t, errorResponse.Message, "image_link")

			t.Log("✓ Correctly rejected approval with invalid image link format")
		})

		t.Run("Valid image link and approval date verification", func(t *testing.T) {
			approveReq := dto.ApproveLoanRequest{
				FieldValidatorProof: "https://example.com/field-validation/proof_789.jpg", // Valid image
				FieldValidatorID:    "validator_003",
			}

			resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/approve", approveReq)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var approveResponse dto.SuccessResponse
			err = json.NewDecoder(resp.Body).Decode(&approveResponse)
			require.NoError(t, err)

			approvedLoan := approveResponse.Data.(map[string]interface{})
			approvalDetails := approvedLoan["approval_details"].(map[string]interface{})
			assert.NotEmpty(t, approvalDetails["approval_date"])
			assert.Equal(t, "https://example.com/field-validation/proof_789.jpg", approvalDetails["field_validator_proof"])

			t.Log("✓ Correctly approved loan with valid image link and recorded approval date")
		})
	})

	t.Run("Investment Validation", func(t *testing.T) {
		t.Log("Testing investment validation...")

		// Create and approve a loan for investment tests
		createReq := dto.CreateLoanRequest{
			BorrowerID:      "borrower_002",
			PrincipalAmount: 30000.00,
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
			FieldValidatorProof: "https://example.com/field-validation/proof_999.png",
			FieldValidatorID:    "validator_004",
		}

		approveResp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/approve", approveReq)
		require.NoError(t, err)
		defer approveResp.Body.Close()

		t.Run("Missing investor ID", func(t *testing.T) {
			investReq := dto.InvestLoanRequest{
				Amount: 30000.00,
				// Missing InvestorID
			}

			resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			t.Log("✓ Correctly rejected investment with missing investor ID")
		})

		t.Run("Invalid investment amount", func(t *testing.T) {
			investReq := dto.InvestLoanRequest{
				InvestorID: "investor_001",
				Amount:     0, // Invalid: must be greater than 0
			}

			resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			t.Log("✓ Correctly rejected investment with invalid amount")
		})

		t.Run("Investment exceeding principal amount", func(t *testing.T) {
			investReq := dto.InvestLoanRequest{
				InvestorID: "investor_001",
				Amount:     35000.00, // Exceeds principal amount of 30000
			}

			resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			t.Log("✓ Correctly rejected investment exceeding principal amount")
		})
	})

	t.Run("Disbursement Validation", func(t *testing.T) {
		t.Log("Testing disbursement validation...")

		// Create, approve, and invest in a loan for disbursement tests
		createReq := dto.CreateLoanRequest{
			BorrowerID:      "borrower_003",
			PrincipalAmount: 20000.00,
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
			FieldValidatorProof: "https://example.com/field-validation/proof_888.jpg",
			FieldValidatorID:    "validator_005",
		}

		approveResp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/approve", approveReq)
		require.NoError(t, err)
		defer approveResp.Body.Close()

		// Invest in the loan
		investReq := dto.InvestLoanRequest{
			InvestorID: "investor_001",
			Amount:     20000.00,
		}

		investResp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq)
		require.NoError(t, err)
		defer investResp.Body.Close()

		t.Run("Missing signed agreement link", func(t *testing.T) {
			disburseReq := dto.DisburseLoanRequest{
				FieldOfficerID: "officer_001",
				// Missing SignedAgreementLink
			}

			resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/disburse", disburseReq)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			t.Log("✓ Correctly rejected disbursement with missing signed agreement link")
		})

		t.Run("Missing field officer ID", func(t *testing.T) {
			disburseReq := dto.DisburseLoanRequest{
				SignedAgreementLink: "https://example.com/signed-agreements/loan_002_signed.pdf",
				// Missing FieldOfficerID
			}

			resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/disburse", disburseReq)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			t.Log("✓ Correctly rejected disbursement with missing field officer ID")
		})
	})

	t.Log("=== Input validation scenarios test completed ===")
}

// ========== BUSINESS RULE TESTS ==========

func TestBusinessRuleEnforcement(t *testing.T) {
	setup := testutils.SetupTestServer()
	defer setup.Server.Close()

	baseURL := setup.Server.URL

	t.Log("=== Testing Business Rule Enforcement ===")

	t.Run("State Transition Rules", func(t *testing.T) {
		t.Log("Testing state transition business rules...")

		t.Run("Cannot invest in proposed loan", func(t *testing.T) {
			// Create a loan without approving
			createReq := dto.CreateLoanRequest{
				BorrowerID:      "borrower_004",
				PrincipalAmount: 15000.00,
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

			// Try to invest in proposed loan
			investReq := dto.InvestLoanRequest{
				InvestorID: "investor_001",
				Amount:     15000.00,
			}

			resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			t.Log("✓ Correctly enforced rule: cannot invest in proposed loan")
		})

		t.Run("Cannot disburse partially invested loan", func(t *testing.T) {
			// Create and approve a loan
			createReq := dto.CreateLoanRequest{
				BorrowerID:      "borrower_005",
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
				FieldValidatorProof: "https://example.com/field-validation/proof_777.png",
				FieldValidatorID:    "validator_006",
			}

			approveResp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/approve", approveReq)
			require.NoError(t, err)
			defer approveResp.Body.Close()

			// Invest partially
			investReq := dto.InvestLoanRequest{
				InvestorID: "investor_001",
				Amount:     15000.00, // Only 60% of principal
			}

			investResp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq)
			require.NoError(t, err)
			defer investResp.Body.Close()

			// Try to disburse partially invested loan
			disburseReq := dto.DisburseLoanRequest{
				SignedAgreementLink: "https://example.com/signed-agreements/loan_003_signed.pdf",
				FieldOfficerID:      "officer_001",
			}

			resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/disburse", disburseReq)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			t.Log("✓ Correctly enforced rule: cannot disburse partially invested loan")
		})

		t.Run("Cannot approve already approved loan", func(t *testing.T) {
			// Create a loan
			createReq := dto.CreateLoanRequest{
				BorrowerID:      "borrower_006",
				PrincipalAmount: 10000.00,
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

			// First approval
			approveReq := dto.ApproveLoanRequest{
				FieldValidatorProof: "https://example.com/field-validation/proof_666.jpg",
				FieldValidatorID:    "validator_007",
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
			t.Log("✓ Correctly enforced rule: cannot approve already approved loan")
		})
	})

	t.Run("Investment Limit Rules", func(t *testing.T) {
		t.Log("Testing investment limit business rules...")

		// Create and approve a loan
		createReq := dto.CreateLoanRequest{
			BorrowerID:      "borrower_007",
			PrincipalAmount: 50000.00,
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
			FieldValidatorProof: "https://example.com/field-validation/proof_555.png",
			FieldValidatorID:    "validator_008",
		}

		approveResp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/approve", approveReq)
		require.NoError(t, err)
		defer approveResp.Body.Close()

		// First investment: 30000
		investReq1 := dto.InvestLoanRequest{
			InvestorID: "investor_001",
			Amount:     30000.00,
		}

		investResp1, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq1)
		require.NoError(t, err)
		defer investResp1.Body.Close()

		assert.Equal(t, http.StatusOK, investResp1.StatusCode)

		// Second investment: 25000 (total would be 55000, exceeding 50000)
		investReq2 := dto.InvestLoanRequest{
			InvestorID: "investor_002",
			Amount:     25000.00,
		}

		investResp2, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq2)
		require.NoError(t, err)
		defer investResp2.Body.Close()

		assert.Equal(t, http.StatusBadRequest, investResp2.StatusCode)
		t.Log("✓ Correctly enforced rule: total investment cannot exceed principal amount")
	})

	t.Log("=== Business rule enforcement test completed ===")
}

// ========== ERROR HANDLING TESTS ==========

func TestErrorHandlingScenarios(t *testing.T) {
	setup := testutils.SetupTestServer()
	defer setup.Server.Close()

	baseURL := setup.Server.URL

	t.Log("=== Testing Error Handling Scenarios ===")

	t.Run("Non-existent resource handling", func(t *testing.T) {
		t.Log("Testing non-existent resource error handling...")

		t.Run("Approve non-existent loan", func(t *testing.T) {
			approveReq := dto.ApproveLoanRequest{
				FieldValidatorProof: "https://example.com/field-validation/proof_444.jpg",
				FieldValidatorID:    "validator_009",
			}

			resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/non-existent-id/approve", approveReq)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			t.Log("✓ Correctly handled attempt to approve non-existent loan")
		})

		t.Run("Invest in non-existent loan", func(t *testing.T) {
			investReq := dto.InvestLoanRequest{
				InvestorID: "investor_001",
				Amount:     10000.00,
			}

			resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/non-existent-id/invest", investReq)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			t.Log("✓ Correctly handled attempt to invest in non-existent loan")
		})

		t.Run("Disburse non-existent loan", func(t *testing.T) {
			disburseReq := dto.DisburseLoanRequest{
				SignedAgreementLink: "https://example.com/signed-agreements/loan_004_signed.pdf",
				FieldOfficerID:      "officer_001",
			}

			resp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/non-existent-id/disburse", disburseReq)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			t.Log("✓ Correctly handled attempt to disburse non-existent loan")
		})
	})

	t.Run("Duplicate operation handling", func(t *testing.T) {
		t.Log("Testing duplicate operation error handling...")

		// Create, approve, invest, and disburse a loan
		createReq := dto.CreateLoanRequest{
			BorrowerID:      "borrower_008",
			PrincipalAmount: 20000.00,
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

		// Approve
		approveReq := dto.ApproveLoanRequest{
			FieldValidatorProof: "https://example.com/field-validation/proof_333.jpg",
			FieldValidatorID:    "validator_010",
		}

		approveResp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/approve", approveReq)
		require.NoError(t, err)
		defer approveResp.Body.Close()

		// Invest
		investReq := dto.InvestLoanRequest{
			InvestorID: "investor_001",
			Amount:     20000.00,
		}

		investResp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/invest", investReq)
		require.NoError(t, err)
		defer investResp.Body.Close()

		// Disburse
		disburseReq := dto.DisburseLoanRequest{
			SignedAgreementLink: "https://example.com/signed-agreements/loan_005_signed.pdf",
			FieldOfficerID:      "officer_001",
		}

		disburseResp, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/disburse", disburseReq)
		require.NoError(t, err)
		defer disburseResp.Body.Close()

		assert.Equal(t, http.StatusOK, disburseResp.StatusCode)

		// Try to disburse again
		disburseResp2, err := testutils.MakeRequest("PUT", baseURL+"/api/v1/loans/"+loanID+"/disburse", disburseReq)
		require.NoError(t, err)
		defer disburseResp2.Body.Close()

		assert.Equal(t, http.StatusBadRequest, disburseResp2.StatusCode)
		t.Log("✓ Correctly handled attempt to disburse already disbursed loan")
	})

	t.Log("=== Error handling scenarios test completed ===")
}
