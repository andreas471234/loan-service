# Documentation

This directory contains comprehensive documentation for the Loan Service API project.

## Quick Navigation

- **[API Documentation](#api-documentation)** - Complete API reference and examples
- **[Testing Guide](#testing-guide)** - How to run tests and understand test coverage
- **[Development Setup](#development-setup)** - Getting started with development
- **[Deployment Guide](#deployment-guide)** - How to deploy the application

## API Documentation

### Core Features
- **Complete Loan Lifecycle Management** - From proposal to disbursement
- **Finite State Machine (FSM)** - Robust state management with forward-only transitions
- **Multi-Investor Support** - Multiple investors can contribute to a single loan
- **Custom Validation** - Image link validation and business rule enforcement
- **Auto-Generation** - Agreement letter links generated when fully invested

### API Endpoints

#### Core Loan Operations
- `GET /api/v1/loans` - Get all loans
- `GET /api/v1/loans/{id}` - Get specific loan
- `POST /api/v1/loans` - Create new loan
- `PUT /api/v1/loans/{id}` - Update loan (proposed status only)
- `DELETE /api/v1/loans/{id}` - Delete loan (proposed status only)

#### Loan State Transitions
- `GET /api/v1/loans/{id}/transitions` - Get valid state transitions
- `PUT /api/v1/loans/{id}/approve` - Approve loan
- `PUT /api/v1/loans/{id}/invest` - Invest in loan
- `PUT /api/v1/loans/{id}/disburse` - Disburse loan

#### Health Check
- `GET /health` - Service health status

### Loan Workflow

#### State Progression
1. **Proposed** → Initial state when loan is created
2. **Approved** → Loan has been approved for funding
3. **Invested** → Funds have been invested in the loan
4. **Disbursed** → Loan amount has been disbursed to borrower

#### Business Rules
- Loans can only move forward in the lifecycle (no rollback)
- Only loans in **Proposed** status can be updated or deleted
- Only loans in **Approved** status can be invested
- Only loans in **Invested** status can be disbursed
- Total investment cannot exceed loan principal amount
- Agreement letter links are auto-generated when fully invested

## Testing Guide

### Running Tests
```bash
make test              # Run all tests
make test-unit         # Run unit tests only
make test-integration  # Run integration tests only
make coverage          # Generate coverage report
```

### Test Documentation
- **[E2E_TEST_CASES.md](E2E_TEST_CASES.md)** - Comprehensive integration test scenarios
- **Coverage Reports** - Generated in `tools/coverage.html`

### Test Categories
1. **Unit Tests** - Individual component testing
2. **Integration Tests** - End-to-end API validation
3. **E2E Tests** - Complete workflow scenarios

## Development Setup

### Prerequisites
- Go 1.21 or higher
- Git
- Docker (optional)

### Quick Start
```bash
git clone <repository-url>
cd loan-service
make deps
cp env.example .env
make run
```

### Development Commands
```bash
make run              # Run the application
make fmt              # Format code
make lint             # Run linter
make docker-up        # Start with Docker Compose
make help             # Show all available commands
```

### Environment Configuration
```env
ENVIRONMENT=development
PORT=8080
DB_DRIVER=sqlite
DB_NAME=loan_service.db
```

## Deployment Guide

### Docker Deployment
```bash
make docker-build      # Build Docker image
make docker-run        # Run Docker container
make docker-up         # Start with Docker Compose
```

### Manual Docker Commands
```bash
docker build -t loan-service -f deployments/Dockerfile .
docker run -p 8080:8080 loan-service
docker-compose -f deployments/docker-compose.yml up -d
```

## API Testing

### Postman Collection
Ready-to-use Postman collection available in `assets/postman/`:
- `loan-service-api.postman_collection.json` - Complete API collection
- `loan-service-api.postman_environment.json` - Environment variables

### Setup Instructions
1. Import both files into Postman
2. Select "Loan Service API Environment"
3. Set base URL to `http://localhost:8080`
4. Start testing API endpoints

For detailed Postman setup, see **[POSTMAN_SETUP.md](POSTMAN_SETUP.md)**.

## Project Structure

```
loan-service/
├── api/v1/              # API routes and handlers
├── cmd/server/          # Application entry point
├── internal/            # Application code
│   ├── config/         # Configuration management
│   ├── database/       # Database connection and setup
│   ├── domain/         # Business entities and FSM
│   ├── dto/            # Data Transfer Objects
│   ├── handler/        # HTTP request handlers
│   ├── middleware/     # HTTP middleware
│   ├── repository/     # Data access layer
│   ├── service/        # Business logic layer
│   └── testutils/      # Test utilities
├── tests/              # Test files
├── docs/               # Documentation (this directory)
├── assets/             # Static assets (Postman collections)
├── deployments/        # Deployment configuration
├── tools/              # Development tools and generated files
└── data/               # Data files
```

## Additional Resources

- **[business_requirement.txt](business_requirement.txt)** - Original business requirements
- **[Main README](../README.md)** - Project overview and quick start
- **[Makefile](../Makefile)** - Available development commands
