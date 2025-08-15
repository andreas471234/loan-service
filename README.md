# Loan Service API with FSM

A comprehensive REST API service for managing the full lifecycle of loans, from proposal to disbursement. Uses a Finite State Machine (FSM) for robust state management with forward-only state transitions and business rules to ensure accurate loan progression and data integrity.

## Features

- **Loan Lifecycle Management**: Complete workflow from proposal to disbursement
- **Finite State Machine (FSM)**: Robust state management with explicit transitions
- **Forward-Only State Transitions**: Enforced business rules for loan status changes (no rollback)
- **RESTful API**: Clean, intuitive API design
- **Database Integration**: SQLite database with GORM ORM
- **Validation**: Request validation and error handling
- **CORS Support**: Cross-origin resource sharing enabled

## Tech Stack

- **Go 1.21+**: Core programming language
- **Gin**: HTTP web framework
- **GORM**: Object-relational mapping
- **SQLite**: Lightweight database
- **UUID**: Unique identifier generation
- **FSM**: Finite State Machine for state management

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Git

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd loan-service
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the application:
```bash
go run .
```

The API will be available at `http://localhost:8080`

## API Endpoints

### Health Check
- `GET /health` - Service health status

### Loan Management

#### Get All Loans
- `GET /api/v1/loans` - Retrieve all loans
- Query Parameters:
  - `status` (optional): Filter by loan status
  - `borrower_id` (optional): Filter by borrower ID

#### Get Single Loan
- `GET /api/v1/loans/{id}` - Retrieve a specific loan by ID

#### Create Loan
- `POST /api/v1/loans` - Create a new loan
- Request Body:
```json
{
  "borrower_id": "string",
  "principal_amount": 10000.00,
  "rate": 5.5,
  "roi": 7.0,
  "agreement_letter_link": "https://example.com/agreement.pdf"
}
```

#### Update Loan
- `PUT /api/v1/loans/{id}` - Update an existing loan (only in proposed status)
- Request Body (all fields optional):
```json
{
  "principal_amount": 15000.00,
  "rate": 6.0,
  "roi": 8.0,
  "agreement_letter_link": "https://example.com/updated-agreement.pdf"
}
```

#### Delete Loan
- `DELETE /api/v1/loans/{id}` - Delete a loan (only in proposed status)

### Loan State Transitions (FSM)

#### Get Valid Transitions
- `GET /api/v1/loans/{id}/transitions` - Get valid state transitions for a loan
- Response:
```json
{
  "message": "Valid transitions retrieved successfully",
  "data": {
    "current_state": "proposed",
    "transitions": [
      {
        "from": "proposed",
        "to": "approved",
        "action": "approve"
      }
    ]
  }
}
```

#### Approve Loan
- `PUT /api/v1/loans/{id}/approve` - Approve a proposed loan

#### Invest in Loan
- `PUT /api/v1/loans/{id}/invest` - Mark an approved loan as invested

#### Disburse Loan
- `PUT /api/v1/loans/{id}/disburse` - Disburse an invested loan

## Loan States

1. **Proposed** - Initial state when loan is created
2. **Approved** - Loan has been approved for funding
3. **Invested** - Funds have been invested in the loan
4. **Disbursed** - Loan amount has been disbursed to borrower

## FSM State Transition Rules

The Finite State Machine enforces the following transitions:

- **Proposed** → **Approved** (via "approve" action)
- **Approved** → **Invested** (via "invest" action)
- **Invested** → **Disbursed** (via "disburse" action)

### Business Rules

- Loans can only move forward in the lifecycle (no rollback allowed)
- Only loans in **Proposed** status can be updated or deleted
- Only loans in **Proposed** status can be approved
- Only loans in **Approved** status can be invested
- Only loans in **Invested** status can be disbursed

## Loan Model Fields

The loan model contains only the essential information:

- **borrower_id**: Unique identifier for the borrower
- **principal_amount**: The loan amount requested
- **rate**: Interest rate that the borrower will pay
- **roi**: Return on investment for investors
- **agreement_letter_link**: Link to the generated agreement letter
- **status**: Current state of the loan
- **created_at**: Timestamp when loan was created
- **updated_at**: Timestamp when loan was last updated

## FSM Implementation

The service uses a Finite State Machine to manage loan states:

```go
type FSM struct {
    CurrentState LoanStatus
    Transitions  []StateTransition
}

type StateTransition struct {
    From   LoanStatus
    To     LoanStatus
    Action string
}
```

### FSM Methods

- `CanTransition(to LoanStatus) bool` - Check if transition is valid
- `Transition(to LoanStatus) error` - Perform state transition
- `GetCurrentState() LoanStatus` - Get current state
- `GetValidTransitions() []StateTransition` - Get valid transitions from current state

## Response Format

### Success Response
```json
{
  "message": "Operation completed successfully",
  "data": {
    // Response data
  }
}
```

### Error Response
```json
{
  "error": "Error type",
  "message": "Detailed error message"
}
```

## Environment Variables

Create a `.env` file in the root directory:

```env
PORT=8080
GIN_MODE=debug
```

## Database

The application uses SQLite as the database. The database file (`loan_service.db`) will be created automatically when you first run the application.

## Development

### Project Structure
```
loan-service/
├── main.go          # Application entry point
├── models.go        # Data models, FSM types and interfaces
├── handlers.go      # HTTP request handlers with FSM logic
├── go.mod           # Go module file
├── go.sum           # Go module checksums
└── README.md        # This file
```

### Running Tests
```bash
go test ./...
```

### Building
```bash
go build -o loan-service .
```

## Example Usage

### Create a Loan
```bash
curl -X POST http://localhost:8080/api/v1/loans \
  -H "Content-Type: application/json" \
  -d '{
    "borrower_id": "user123",
    "principal_amount": 25000.00,
    "rate": 4.5,
    "roi": 6.0,
    "agreement_letter_link": "https://example.com/agreement/user123.pdf"
  }'
```

### Get Valid Transitions
```bash
curl http://localhost:8080/api/v1/loans/{loan-id}/transitions
```

### Get All Loans
```bash
curl http://localhost:8080/api/v1/loans
```

### Approve a Loan
```bash
curl -X PUT http://localhost:8080/api/v1/loans/{loan-id}/approve
```

### Test the API with FSM
```bash
chmod +x test_api.sh
./test_api.sh
```

## License

This project is licensed under the MIT License.
