#!/bin/bash

# Test script for Loan Service API with FSM and State Rules
BASE_URL="http://localhost:8080"

echo "=== Loan Service API Test with FSM and State Rules ==="
echo

# Test health endpoint
echo "1. Testing health endpoint..."
curl -s "$BASE_URL/health" | python3 -m json.tool
echo
echo

# Test creating a loan
echo "2. Creating a new loan..."
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/loans/" \
  -H "Content-Type: application/json" \
  -d '{
    "borrower_id": "user123",
    "principal_amount": 25000.00,
    "rate": 4.5,
    "roi": 6.0,
    "agreement_letter_link": "https://example.com/agreement/user123.pdf"
  }')

echo "$CREATE_RESPONSE" | python3 -m json.tool
echo

# Extract loan ID from response
LOAN_ID=$(echo "$CREATE_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['id'])")
echo "Created loan ID: $LOAN_ID"
echo

# Test getting valid transitions for proposed loan
echo "3. Getting valid transitions for proposed loan..."
curl -s "$BASE_URL/api/v1/loans/$LOAN_ID/transitions" | python3 -m json.tool
echo
echo

# Test getting all loans
echo "4. Getting all loans..."
curl -s "$BASE_URL/api/v1/loans/" | python3 -m json.tool
echo
echo

# Test getting specific loan
echo "5. Getting specific loan..."
curl -s "$BASE_URL/api/v1/loans/$LOAN_ID" | python3 -m json.tool
echo
echo

# Test approving the loan with required details
echo "6. Approving the loan with field validator details..."
curl -s -X PUT "$BASE_URL/api/v1/loans/$LOAN_ID/approve" \
  -H "Content-Type: application/json" \
  -d '{
    "field_validator_proof": "https://example.com/proofs/field_visit_123.jpg",
    "field_validator_id": "validator_001"
  }' | python3 -m json.tool
echo
echo

# Test getting valid transitions for approved loan
echo "7. Getting valid transitions for approved loan..."
curl -s "$BASE_URL/api/v1/loans/$LOAN_ID/transitions" | python3 -m json.tool
echo
echo

# Test adding first investment
echo "8. Adding first investment (15,000)..."
curl -s -X PUT "$BASE_URL/api/v1/loans/$LOAN_ID/invest" \
  -H "Content-Type: application/json" \
  -d '{
    "investor_id": "investor_001",
    "amount": 15000.00
  }' | python3 -m json.tool
echo
echo

# Test adding second investment to complete the loan
echo "9. Adding second investment (10,000) to complete the loan..."
curl -s -X PUT "$BASE_URL/api/v1/loans/$LOAN_ID/invest" \
  -H "Content-Type: application/json" \
  -d '{
    "investor_id": "investor_002",
    "amount": 10000.00
  }' | python3 -m json.tool
echo
echo

# Test getting valid transitions for invested loan
echo "10. Getting valid transitions for invested loan..."
curl -s "$BASE_URL/api/v1/loans/$LOAN_ID/transitions" | python3 -m json.tool
echo
echo

# Test disbursing the loan with required details
echo "11. Disbursing the loan with field officer details..."
curl -s -X PUT "$BASE_URL/api/v1/loans/$LOAN_ID/disburse" \
  -H "Content-Type: application/json" \
  -d '{
    "signed_agreement_link": "https://example.com/agreements/signed_agreement_123.pdf",
    "field_officer_id": "officer_001"
  }' | python3 -m json.tool
echo
echo

# Test getting valid transitions for disbursed loan
echo "12. Getting valid transitions for disbursed loan..."
curl -s "$BASE_URL/api/v1/loans/$LOAN_ID/transitions" | python3 -m json.tool
echo
echo

# Test getting loans by status
echo "13. Getting loans by status (disbursed)..."
curl -s "$BASE_URL/api/v1/loans/?status=disbursed" | python3 -m json.tool
echo
echo

# Test creating another loan
echo "14. Creating another loan..."
CREATE_RESPONSE2=$(curl -s -X POST "$BASE_URL/api/v1/loans/" \
  -H "Content-Type: application/json" \
  -d '{
    "borrower_id": "user456",
    "principal_amount": 15000.00,
    "rate": 5.0,
    "roi": 7.5,
    "agreement_letter_link": "https://example.com/agreement/user456.pdf"
  }')

echo "$CREATE_RESPONSE2" | python3 -m json.tool
echo

# Extract second loan ID
LOAN_ID2=$(echo "$CREATE_RESPONSE2" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['id'])")
echo "Created second loan ID: $LOAN_ID2"
echo

# Test updating the second loan
echo "15. Updating the second loan..."
curl -s -X PUT "$BASE_URL/api/v1/loans/$LOAN_ID2" \
  -H "Content-Type: application/json" \
  -d '{
    "principal_amount": 18000.00,
    "rate": 5.5,
    "roi": 8.0,
    "agreement_letter_link": "https://example.com/agreement/user456_updated.pdf"
  }' | python3 -m json.tool
echo
echo

# Test getting loans by borrower ID
echo "16. Getting loans by borrower ID (user123)..."
curl -s "$BASE_URL/api/v1/loans/?borrower_id=user123" | python3 -m json.tool
echo
echo

# Test invalid state transition (trying to approve a disbursed loan)
echo "17. Testing invalid state transition (trying to approve a disbursed loan)..."
curl -s -X PUT "$BASE_URL/api/v1/loans/$LOAN_ID/approve" \
  -H "Content-Type: application/json" \
  -d '{
    "field_validator_proof": "https://example.com/proofs/field_visit_456.jpg",
    "field_validator_id": "validator_002"
  }' | python3 -m json.tool
echo
echo

# Test invalid investment (trying to invest more than principal)
echo "18. Testing invalid investment (trying to invest more than principal)..."
curl -s -X PUT "$BASE_URL/api/v1/loans/$LOAN_ID2/invest" \
  -H "Content-Type: application/json" \
  -d '{
    "investor_id": "investor_003",
    "amount": 20000.00
  }' | python3 -m json.tool
echo
echo

# Test invalid disbursement (trying to disburse without full investment)
echo "19. Testing invalid disbursement (trying to disburse without full investment)..."
curl -s -X PUT "$BASE_URL/api/v1/loans/$LOAN_ID2/disburse" \
  -H "Content-Type: application/json" \
  -d '{
    "signed_agreement_link": "https://example.com/agreements/signed_agreement_456.pdf",
    "field_officer_id": "officer_002"
  }' | python3 -m json.tool
echo
echo

echo "=== FSM and State Rules Test completed ==="