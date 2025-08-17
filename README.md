# Loan Service API

A modern REST API service for managing the complete lifecycle of loans, from initial proposal to final disbursement. Built with Go, this service implements a robust Finite State Machine (FSM) for reliable state management and enforces strict business rules.

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Git
- Docker (optional)

### Installation & Setup

1. **Clone and setup**

```bash
git clone https://github.com/andreas471234/loan-service.git
cd loan-service
make deps
cp env.example .env
```

2. **Run the application**

```bash
make run
```

The API will be available at `http://localhost:8080`

## Documentation

- **[Complete Documentation](docs/README.md)** - API reference, testing guide, and development setup
- **[E2E Test Cases](docs/E2E_TEST_CASES.md)** - Integration test scenarios
- **[Postman Setup](docs/POSTMAN_SETUP.md)** - API testing collection setup

## Development Commands

```bash
make run              # Run the application
make test             # Run all tests
make test-integration # Run integration tests only
make coverage         # Run tests with coverage report
make fmt              # Format code
make lint             # Run linter
make docker-up        # Start with Docker Compose
make help             # Show all available commands
```

## Key Features

- **Complete Loan Lifecycle** - From proposal to disbursement
- **Multi-Investor Support** - Multiple investors can contribute to a single loan
- **Finite State Machine** - Robust state management with forward-only transitions
- **Custom Validation** - Image link validation and business rule enforcement
- **Auto-Generation** - Agreement letter links generated when fully invested

## API Overview

- **Core Operations**: CRUD operations for loan management
- **State Transitions**: Approve, invest, and disburse loans
- **Health Check**: Service monitoring endpoint

For complete API documentation, see [docs/README.md](docs/README.md#api-documentation).

## Project Structure

```text
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
├── docs/               # Documentation
│   ├── E2E_TEST_CASES.md
│   ├── POSTMAN_SETUP.md
│   └── business_requirement.txt
├── assets/             # Static assets
│   └── postman/        # Postman collections
├── deployments/        # Deployment configuration
│   ├── Dockerfile
│   └── docker-compose.yml
├── tools/              # Development tools and generated files
├── data/               # Data files
├── Makefile            # Build automation
├── go.mod              # Go module file
└── env.example         # Environment template
```

## Quick Examples

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

### Approve a Loan

```bash
curl -X PUT http://localhost:8080/api/v1/loans/{loan-id}/approve \
  -H "Content-Type: application/json" \
  -d '{
    "field_validator_proof": "https://example.com/proof.jpg",
    "field_validator_id": "validator123"
  }'
```

## Testing

- **Unit Tests**: Individual component testing
- **Integration Tests**: End-to-end API validation
- **E2E Tests**: Complete workflow scenarios

For comprehensive testing guide, see [docs/README.md](docs/README.md#testing-guide).

## License

This project is licensed under the MIT License.
