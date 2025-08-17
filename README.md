# Loan Service API

A modern, production-ready REST API service for managing the complete lifecycle of loans, from initial proposal to final disbursement. Built with Go, this service implements a robust Finite State Machine (FSM) for reliable state management and enforces strict business rules to ensure data integrity and accurate loan progression.

## Features

### Core Loan Management
- **Complete Loan Lifecycle**: End-to-end loan processing from proposal to disbursement
- **Finite State Machine (FSM)**: Robust state management with explicit, forward-only transitions
- **Multi-Investor Support**: Multiple investors can contribute to a single loan
- **Automatic State Transitions**: Smart progression based on business rules
- **Investment Tracking**: Real-time tracking of total investments and individual contributions

### Advanced Validation & Security
- **Custom Image Link Validation**: Validates field validator proof images with support for multiple formats
- **Comprehensive Input Validation**: All API endpoints with strict data validation
- **Business Rule Enforcement**: Strict state transition and investment limit validation
- **Data Integrity Protection**: Prevents invalid operations and maintains audit trails

### Auto-Generation Features
- **Smart Agreement Letter Links**: Automatically generated when loans become fully invested
- **Automatic Date Recording**: Approval and disbursement dates automatically captured
- **UUID Generation**: Secure, unique identifiers for all entities
- **Audit Trail**: Complete tracking of all loan activities and state changes

### API & Integration
- **RESTful API Design**: Clean, intuitive API with comprehensive documentation
- **CORS Support**: Cross-origin resource sharing for web applications
- **Request/Response Logging**: Detailed logging for debugging and monitoring
- **Error Recovery**: Graceful error handling with meaningful error messages
- **Health Check Endpoint**: Service health monitoring

### Database & Storage
- **Multi-Database Support**: SQLite for development, PostgreSQL for production
- **GORM Integration**: Powerful ORM with automatic migrations
- **Soft Deletes**: Safe deletion with data preservation
- **Optimized Queries**: Efficient database operations with proper indexing

### Testing & Quality Assurance
- **Comprehensive Test Suite**: 10 test files with extensive coverage
- **Unit Testing**: Individual component testing
- **Integration Testing**: End-to-end API testing
- **Test Utilities**: Reusable test helpers and utilities

### Development & Operations
- **Docker Support**: Containerized deployment with Docker and Docker Compose
- **Makefile Automation**: Streamlined development workflow
- **Code Formatting**: Automatic code formatting with `go fmt`
- **Linting**: Code quality checks with `go vet`
- **Race Detection**: Concurrency testing for thread safety

## Tech Stack

### Backend Framework
- **Go 1.21+**: High-performance, concurrent programming language
- **Gin**: Fast HTTP web framework with middleware support
- **GORM**: Feature-rich ORM for database operations

### Database
- **SQLite**: Lightweight, file-based database (development)
- **PostgreSQL**: Production-ready relational database (configurable)

### Validation & Security
- **go-playground/validator**: Powerful validation library
- **Custom Validators**: Image link validation and business rule enforcement
- **UUID**: Secure unique identifier generation

### Testing & Quality
- **Go Testing**: Built-in testing framework
- **Test Coverage**: Comprehensive coverage reporting
- **Integration Tests**: End-to-end API validation

### DevOps & Deployment
- **Docker**: Containerization for consistent deployment
- **Docker Compose**: Multi-service orchestration
- **Make**: Build automation and task management

## Project Structure

```
loan-service/
├── api/                    # API versioning and routing
│   └── v1/                # API v1 routes and handlers
├── cmd/                    # Application entry points
│   └── server/            # Main server application
├── internal/               # Private application code
│   ├── config/            # Configuration management
│   ├── database/          # Database connection and setup
│   ├── domain/            # Business entities and logic
│   ├── dto/               # Data Transfer Objects
│   ├── handler/           # HTTP request handlers
│   ├── middleware/        # HTTP middleware (CORS, logging, recovery)
│   ├── repository/        # Data access layer
│   ├── service/           # Business logic layer
│   └── testutils/         # Test utilities and helpers
├── tests/                 # Test files
│   └── integration/       # Integration tests
├── docker-compose.yml     # Docker Compose configuration
├── Dockerfile             # Docker image definition
├── Makefile               # Build and development tasks
├── go.mod                 # Go module file
├── go.sum                 # Go module checksums
├── env.example            # Environment variables template
├── E2E_TEST_CASES.md      # Integration test documentation
└── README.md              # This file
```

## Quick Start

### Prerequisites
- Go 1.21 or higher
- Git
- Docker (optional, for containerized deployment)

### Installation

1. **Clone the repository**
```bash
git clone <repository-url>
cd loan-service
```

2. **Install dependencies**
```bash
make deps
```

3. **Set up environment**
```bash
cp env.example .env
# Edit .env file with your configuration
```

4. **Run the application**
```bash
make run
```

The API will be available at `http://localhost:8080`

## Development Commands

### Build & Run
```bash
make build          # Build the application
make run            # Run the application
make run-race       # Run with race detection
```

### Testing
```bash
make test           # Run all tests
make test-unit      # Run unit tests only
make test-integration # Run integration tests only
make coverage       # Run tests with coverage report
```

### Code Quality
```bash
make fmt            # Format code
make lint           # Run linter
make check          # Run format, lint, and test
```

### Docker
```bash
make docker-build   # Build Docker image
make docker-run     # Run Docker container
make docker-up      # Start with Docker Compose
make docker-down    # Stop Docker Compose services
```

### Utilities
```bash
make clean          # Clean build artifacts
make deps           # Install dependencies
make help           # Show all available commands
```

## API Endpoints

### Health Check
- `GET /health` - Service health status

### Loan Management

#### Get All Loans
- `GET /api/v1/loans` - Retrieve all loans
- **Query Parameters**:
  - `status` (optional): Filter by loan status
  - `borrower_id` (optional): Filter by borrower ID

#### Get Single Loan
- `GET /api/v1/loans/{id}` - Retrieve a specific loan by ID

#### Create Loan
- `POST /api/v1/loans` - Create a new loan
- **Request Body**:
```json
{
  "borrower_id": "string",
  "principal_amount": 10000.00,
  "rate": 5.5,
  "roi": 7.0
}
```

#### Update Loan
- `PUT /api/v1/loans/{id}` - Update an existing loan (only in proposed status)
- **Request Body** (all fields optional):
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
- **Response**:
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
- **Request Body**:
```json
{
  "field_validator_proof": "https://example.com/proof.jpg",
  "field_validator_id": "validator123"
}
```
- **Note**: `field_validator_proof` must be a valid image link (supports .jpg, .jpeg, .png, .gif, .bmp, .webp, .svg)

#### Invest in Loan
- `PUT /api/v1/loans/{id}/invest` - Add investment to an approved loan
- **Request Body**:
```json
{
  "investor_id": "investor123",
  "amount": 5000.00
}
```
- **Note**: Agreement letter link is automatically generated when loan becomes fully invested

#### Disburse Loan
- `PUT /api/v1/loans/{id}/disburse` - Disburse a fully invested loan
- **Request Body**:
```json
{
  "signed_agreement_link": "https://example.com/signed-agreement.pdf",
  "field_officer_id": "officer123"
}
```

## Loan States & Workflow

### State Progression
1. **Proposed** → Initial state when loan is created
2. **Approved** → Loan has been approved for funding
3. **Invested** → Funds have been invested in the loan
4. **Disbursed** → Loan amount has been disbursed to borrower

### FSM State Transition Rules
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
- Agreement letter links are auto-generated when loans become fully invested
- Approval dates are automatically recorded when loans are approved
- Multiple investors can contribute to the same loan
- Total investment cannot exceed loan principal amount

## Architecture

The application follows Clean Architecture principles with clear separation of concerns:

### Layer Structure
- **Domain Layer** (`internal/domain/`): Business entities, logic, and FSM implementation
- **Repository Layer** (`internal/repository/`): Data access abstraction and database operations
- **Service Layer** (`internal/service/`): Business logic orchestration and validation
- **Handler Layer** (`internal/handler/`): HTTP request handling and response formatting
- **DTO Layer** (`internal/dto/`): Data transfer objects for API contracts
- **Middleware Layer** (`internal/middleware/`): HTTP middleware (CORS, logging, recovery)
- **Config Layer** (`internal/config/`): Configuration management and environment handling
- **Database Layer** (`internal/database/`): Database connection and setup

### Key Components
- **27 Go source files** with comprehensive functionality
- **10 test files** ensuring code quality and reliability
- **Custom validation** with image link validation
- **Middleware stack** for logging, CORS, and error recovery
- **FSM implementation** for robust state management

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

## Environment Configuration

Copy `env.example` to `.env` and configure:

```env
# Environment
ENVIRONMENT=development

# Server Configuration
PORT=8080
SERVER_READ_TIMEOUT=10
SERVER_WRITE_TIMEOUT=10
SERVER_IDLE_TIMEOUT=120

# Database Configuration (SQLite - Default)
DB_DRIVER=sqlite
DB_NAME=loan_service.db

# Database Configuration (PostgreSQL - Optional)
# DB_DRIVER=postgres
# DB_HOST=localhost
# DB_PORT=5432
# DB_USER=postgres
# DB_PASSWORD=password
# DB_NAME=loan_service
# DB_SSLMODE=disable
```

## Database

The application uses SQLite as the default database for development. The database file is created automatically on first run. PostgreSQL is supported for production deployments and can be configured by updating the environment variables.

### Features
- **Automatic Migration**: Database schema created automatically
- **Soft Deletes**: Safe deletion with data preservation
- **Audit Fields**: Created, updated, and deleted timestamps
- **Foreign Key Relationships**: Proper data integrity constraints

## Testing

### Test Coverage
- **Unit Tests**: Individual component testing
- **Integration Tests**: End-to-end API validation
- **Test Utilities**: Reusable test helpers

### Running Tests
```bash
make test              # Run all tests
make test-unit         # Run unit tests only
make test-integration  # Run integration tests only
make coverage          # Generate coverage report
```

### Test Documentation
See `E2E_TEST_CASES.md` for comprehensive integration test documentation.

## Docker Deployment

### Build and Run
```bash
make docker-build      # Build Docker image
make docker-run        # Run Docker container
```

### Docker Compose
```bash
make docker-up         # Start services
make docker-down       # Stop services
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
    "roi": 6.0
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

### Invest in a Loan
```bash
curl -X PUT http://localhost:8080/api/v1/loans/{loan-id}/invest \
  -H "Content-Type: application/json" \
  -d '{
    "investor_id": "investor123",
    "amount": 25000.00
  }'
```

### Disburse a Loan
```bash
curl -X PUT http://localhost:8080/api/v1/loans/{loan-id}/disburse \
  -H "Content-Type: application/json" \
  -d '{
    "signed_agreement_link": "https://example.com/signed-agreement.pdf",
    "field_officer_id": "officer123"
  }'
```

## Key Features Summary

### Core Functionality
- Complete loan lifecycle management
- Multi-investor support with partial investments
- Automatic state transitions and business rule enforcement
- Comprehensive validation and error handling

### Advanced Features
- Custom image link validation for field validator proofs
- Auto-generation of agreement letter links
- Automatic date recording for approvals and disbursements
- Robust FSM implementation with forward-only transitions

### Developer Experience
- Comprehensive test suite with 10 test files
- Docker support for easy deployment
- Makefile automation for common tasks
- Detailed API documentation and examples

### Production Ready
- CORS support for web applications
- Request logging and error recovery
- Multi-database support (SQLite/PostgreSQL)
- Health check endpoint for monitoring

## License

This project is licensed under the MIT License.
