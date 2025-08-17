# Postman Collection Setup Guide

This guide explains how to set up and use the Postman collection for the Loan Service API.

## Files Included

1. **`loan-service-api.postman_collection.json`** - Complete API collection with all endpoints
2. **`loan-service-api.postman_environment.json`** - Environment variables for different configurations
3. **`POSTMAN_SETUP.md`** - This setup guide

## Quick Setup

### 1. Import Collection and Environment

1. Open Postman
2. Click **Import** button
3. Import both files:
   - `loan-service-api.postman_collection.json`
   - `loan-service-api.postman_environment.json`

### 2. Select Environment

1. In the top-right corner of Postman, select **"Loan Service API Environment"**
2. Verify the `base_url` is set to `http://localhost:8080` (or your server URL)

### 3. Start the Loan Service

Before testing, ensure the loan service is running:

```bash
# Option 1: Run locally
make run

# Option 2: Run with Docker
make docker-up

# Option 3: Run Docker container directly
docker run -p 8080:8080 loan-service
```

## Collection Structure

### 1. Health Check
- **GET /health** - Verify service is running

### 2. Loan Management
Basic CRUD operations for loan management:
- **GET /api/v1/loans/** - Get all loans
- **GET /api/v1/loans/{id}** - Get specific loan
- **POST /api/v1/loans/** - Create new loan
- **PUT /api/v1/loans/{id}** - Update loan
- **DELETE /api/v1/loans/{id}** - Delete loan

### 3. Loan State Transitions
Business operations for loan lifecycle:
- **PUT /api/v1/loans/{id}/approve** - Approve loan
- **PUT /api/v1/loans/{id}/invest** - Invest in loan
- **PUT /api/v1/loans/{id}/disburse** - Disburse loan
- **GET /api/v1/loans/{id}/transitions** - Get available transitions

### 4. Complete Loan Lifecycle
Step-by-step workflow demonstration:
1. Create Loan
2. Approve Loan
3. Invest in Loan (First Investment)
4. Invest in Loan (Second Investment)
5. Disburse Loan
6. Verify Final State

### 5. Validation Examples
Examples of validation errors and business rule violations:
- Missing required fields
- Invalid amounts
- Invalid image links
- Exceeding investment limits

## Environment Variables

The environment includes these variables:

| Variable | Description | Default Value |
|----------|-------------|---------------|
| `base_url` | API base URL | `http://localhost:8080` |
| `loan_id` | Current loan ID (auto-set) | Empty |
| `borrower_id` | Test borrower ID | `test_borrower_001` |
| `investor_id` | Test investor ID | `test_investor_001` |
| `field_validator_id` | Test validator ID | `test_validator_001` |
| `field_officer_id` | Test officer ID | `test_officer_001` |
| `principal_amount` | Test principal amount | `10000.00` |
| `rate` | Test interest rate | `5.0` |
| `roi` | Test ROI | `7.0` |
| `investment_amount` | Test investment amount | `5000.00` |
| `field_validator_proof` | Test validation proof URL | `https://example.com/images/validation_proof.jpg` |
| `signed_agreement_link` | Test agreement link | `https://example.com/documents/signed_agreement.pdf` |

## Usage Examples

### Testing the Complete Lifecycle

1. **Start with Health Check**
   - Run the "Health Check" request to verify the service is running

2. **Create a Loan**
   - Use "Create Loan" from the "Loan Management" folder
   - The response will contain a loan ID

3. **Follow the Complete Lifecycle**
   - Use the "Complete Loan Lifecycle" folder
   - Each request will automatically set the `loan_id` variable
   - Run requests in sequence: 1 → 2 → 3 → 4 → 5 → 6

### Testing Individual Operations

1. **Create a loan** using "Create Loan"
2. **Copy the loan ID** from the response
3. **Set the loan_id variable** in the environment
4. **Test state transitions** in any order (following business rules)

### Testing Validation

Use the "Validation Examples" folder to test:
- Input validation
- Business rule enforcement
- Error handling

## API Features Demonstrated

### Auto-Generation Features
- **Agreement Letter Links**: Automatically generated when loans are fully invested
- **Approval Dates**: Automatically recorded when loans are approved

### Validation Features
- **Image Link Validation**: Ensures `field_validator_proof` contains valid image URLs
- **Business Rules**: Enforces investment limits and state transitions
- **Required Fields**: Validates all required parameters

### State Management
- **Finite State Machine**: Proper loan state transitions
- **Business Logic**: Enforces rules like "cannot invest in proposed loans"

## Troubleshooting

### Common Issues

1. **Connection Refused**
   - Ensure the loan service is running
   - Check the `base_url` variable
   - Verify the port (default: 8080)

2. **Validation Errors**
   - Check request body format
   - Ensure all required fields are provided
   - Verify data types (numbers vs strings)

3. **Business Rule Violations**
   - Follow the loan lifecycle: proposed → approved → invested → disbursed
   - Check loan status before attempting operations
   - Ensure investment amounts don't exceed principal

### Environment Setup

To create additional environments:

1. **Development Environment**
   ```
   base_url: http://localhost:8080
   ```

2. **Docker Environment**
   ```
   base_url: http://localhost:8080
   ```

3. **Production Environment**
   ```
   base_url: https://your-production-domain.com
   ```

## Advanced Usage

### Using Variables in Requests

You can use environment variables in request bodies:

```json
{
  "borrower_id": "{{borrower_id}}",
  "principal_amount": {{principal_amount}},
  "rate": {{rate}},
  "roi": {{roi}}
}
```

### Automated Testing

The collection includes test scripts that:
- Automatically set the `loan_id` variable after creating a loan
- Log important information to the console
- Can be extended for automated testing workflows

### Collection Runner

Use Postman's Collection Runner to:
- Run all requests in sequence
- Test complete workflows
- Generate test reports
- Automate API testing

## Support

For issues with the API:
1. Check the service logs
2. Verify the database connection
3. Review the API documentation in the README.md
4. Run the integration tests: `make test-integration`
