# Loan Service API

A comprehensive REST API service for managing the full lifecycle of loans, from proposal to disbursement. Supports multiple loan states including Proposed, Approved, Invested, and Disbursed, with business rules and workflows to ensure accurate state transitions and data integrity.

## Features

- **Loan Lifecycle Management**: Complete workflow from proposal to disbursement
- **State Transitions**: Enforced business rules for loan status changes
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
  "amount": 10000.00,
  "interest_rate": 5.5,
  "term": 12,
  "purpose": "Home improvement",
  "description": "Optional description"
}
```

#### Update Loan
- `PUT /api/v1/loans/{id}` - Update an existing loan (only in proposed status)
- Request Body (all fields optional):
```json
{
  "amount": 15000.00,
  "interest_rate": 6.0,
  "term": 24,
  "purpose": "Business expansion",
  "description": "Updated description"
}
```

#### Delete Loan
- `DELETE /api/v1/loans/{id}` - Delete a loan (only in proposed status)

### Loan State Transitions

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
5. **Rejected** - Loan has been rejected (not implemented in current version)

## State Transition Rules

- Only loans in **Proposed** status can be updated or deleted
- Only loans in **Proposed** status can be approved
- Only loans in **Approved** status can be invested
- Only loans in **Invested** status can be disbursed

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
├── models.go        # Data models and types
├── handlers.go      # HTTP request handlers
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
    "amount": 25000.00,
    "interest_rate": 4.5,
    "term": 36,
    "purpose": "Vehicle purchase",
    "description": "New car loan"
  }'
```

### Get All Loans
```bash
curl http://localhost:8080/api/v1/loans
```

### Approve a Loan
```bash
curl -X PUT http://localhost:8080/api/v1/loans/{loan-id}/approve
```

## License

This project is licensed under the MIT License.
