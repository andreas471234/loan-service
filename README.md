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
- **Clean Architecture**: Well-structured codebase following Go best practices

## Tech Stack

- **Go 1.21+**: Core programming language
- **Gin**: HTTP web framework
- **GORM**: Object-relational mapping
- **SQLite**: Lightweight database
- **UUID**: Unique identifier generation
- **FSM**: Finite State Machine for state management

## Project Structure

```
loan-service/
├── api/                    # API versioning
│   └── v1/                # API v1 routes
├── cmd/                    # Application entry points
│   └── server/            # Main server application
├── internal/               # Private application code
│   ├── config/            # Configuration management
│   ├── database/          # Database connection and setup
│   ├── domain/            # Business logic and models
│   ├── dto/               # Data Transfer Objects
│   ├── handler/           # HTTP request handlers
│   ├── middleware/        # HTTP middleware
│   ├── repository/        # Data access layer
│   └── service/           # Business logic layer
├── build/                 # Build artifacts (generated)
├── docker-compose.yml     # Docker Compose configuration
├── Dockerfile             # Docker image definition
├── Makefile               # Build and development tasks
├── go.mod                 # Go module file
├── go.sum                 # Go module checksums
├── env.example            # Environment variables example
└── README.md              # This file
```

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
make deps
```

3. Set up environment variables:
```bash
cp env.example .env
# Edit .env file with your configuration
```

4. Run the application:
```bash
make run
```

The API will be available at `http://localhost:8080`

## Development Commands

```bash
# Build the application
make build

# Run tests
make test

# Run tests with coverage
make coverage

# Format code
make fmt

# Run linter
make lint

# Run all checks (format, lint, test)
make check

# Clean build artifacts
make clean

# Run with race detection
make run-race

# Docker commands
make docker-build
make docker-run
make docker-up
make docker-down

# Show all available commands
make help
```

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
- Request Body:
```json
{
  "field_validator_proof": "https://example.com/proof.jpg",
  "field_validator_id": "validator123"
}
```

#### Invest in Loan
- `PUT /api/v1/loans/{id}/invest` - Add investment to an approved loan
- Request Body:
```json
{
  "investor_id": "investor123",
  "amount": 5000.00
}
```

#### Disburse Loan
- `PUT /api/v1/loans/{id}/disburse` - Disburse a fully invested loan
- Request Body:
```json
{
  "signed_agreement_link": "https://example.com/signed-agreement.pdf",
  "field_officer_id": "officer123"
}
```

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
- Loans automatically transition to **Invested** when fully funded

## Architecture

The application follows Clean Architecture principles:

- **Domain Layer** (`internal/domain/`): Business entities and logic
- **Repository Layer** (`internal/repository/`): Data access abstraction
- **Service Layer** (`internal/service/`): Business logic orchestration
- **Handler Layer** (`internal/handler/`): HTTP request handling
- **DTO Layer** (`internal/dto/`): Data transfer objects for API contracts

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

Copy `env.example` to `.env` and configure:

```env
# Environment
ENVIRONMENT=development

# Server Configuration
PORT=8080
SERVER_READ_TIMEOUT=10
SERVER_WRITE_TIMEOUT=10
SERVER_IDLE_TIMEOUT=120

# Database Configuration
DB_DRIVER=sqlite
DB_NAME=loan_service.db
```

## Database

The application uses SQLite as the database. The database file will be created automatically when you first run the application.

## Testing

### Run Tests
```bash
make test
```

### Run Tests with Coverage
```bash
make coverage
```

### Functional Tests
```bash
go test -v ./functional_test.go ./test_utils.go
```

### API Tests
```bash
chmod +x test_api.sh
./test_api.sh
```

## Docker

### Build and Run with Docker
```bash
make docker-build
make docker-run
```

### Using Docker Compose
```bash
make docker-up
make docker-down
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
curl -X PUT http://localhost:8080/api/v1/loans/{loan-id}/approve \
  -H "Content-Type: application/json" \
  -d '{
    "field_validator_proof": "https://example.com/proof.jpg",
    "field_validator_id": "validator123"
  }'
```

## License

This project is licensed under the MIT License.
