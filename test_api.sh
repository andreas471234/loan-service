#!/bin/bash

# Test script for Loan Service API
BASE_URL="http://localhost:8080"

echo "=== Loan Service API Test ==="
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
    "amount": 25000.00,
    "interest_rate": 4.5,
    "term": 36,
    "purpose": "Vehicle purchase",
    "description": "New car loan"
  }')

echo "$CREATE_RESPONSE" | python3 -m json.tool
echo

# Extract loan ID from response
LOAN_ID=$(echo "$CREATE_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['id'])")
echo "Created loan ID: $LOAN_ID"
echo

# Test getting all loans
echo "3. Getting all loans..."
curl -s "$BASE_URL/api/v1/loans/" | python3 -m json.tool
echo
echo

# Test getting specific loan
echo "4. Getting specific loan..."
curl -s "$BASE_URL/api/v1/loans/$LOAN_ID" | python3 -m json.tool
echo
echo

# Test approving the loan
echo "5. Approving the loan..."
curl -s -X PUT "$BASE_URL/api/v1/loans/$LOAN_ID/approve" | python3 -m json.tool
echo
echo

# Test investing in the loan
echo "6. Investing in the loan..."
curl -s -X PUT "$BASE_URL/api/v1/loans/$LOAN_ID/invest" | python3 -m json.tool
echo
echo

# Test disbursing the loan
echo "7. Disbursing the loan..."
curl -s -X PUT "$BASE_URL/api/v1/loans/$LOAN_ID/disburse" | python3 -m json.tool
echo
echo

# Test getting loans by status
echo "8. Getting loans by status (disbursed)..."
curl -s "$BASE_URL/api/v1/loans/?status=disbursed" | python3 -m json.tool
echo
echo

# Test creating another loan
echo "9. Creating another loan..."
CREATE_RESPONSE2=$(curl -s -X POST "$BASE_URL/api/v1/loans/" \
  -H "Content-Type: application/json" \
  -d '{
    "borrower_id": "user456",
    "amount": 15000.00,
    "interest_rate": 6.0,
    "term": 24,
    "purpose": "Home improvement",
    "description": "Kitchen renovation"
  }')

echo "$CREATE_RESPONSE2" | python3 -m json.tool
echo

# Extract second loan ID
LOAN_ID2=$(echo "$CREATE_RESPONSE2" | python3 -c "import sys, json; print(json.load(sys.stdin)['data']['id'])")
echo "Created second loan ID: $LOAN_ID2"
echo

# Test getting all loans again
echo "10. Getting all loans (should show 2 loans)..."
curl -s "$BASE_URL/api/v1/loans/" | python3 -m json.tool
echo
echo

echo "=== Test completed ===" 