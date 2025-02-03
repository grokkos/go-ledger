# go-ledger

A simple ledger API implementation that handles basic financial transactions and balance tracking.

## 📌 Quick Links
- [Installation](#installation)
- [API Endpoints](#-api-endpoints)
- [Example API Usage](#-example-api-usage)
- [Testing](#-testing)
- [Implementation Details](#-implementation-details)
- [Future Improvements](#-future-improvements)

---

## Project Overview

This project implements a RESTful API for a basic ledger system that can:

- Record money movements (deposits and withdrawals)
- View current balance
- View transaction history

## 🏗️ Design Decisions

### 1️⃣ Architecture
- Clean architecture principles with clear separation of concerns
- Layered approach: handlers → services → models
- Dependency injection for better testability
- Interface-based design for flexibility and maintainability

### 2️⃣ Project Structure
```plaintext
go-ledger/
├── cmd/                    # Application entry points
│   └── api/
│       └── main.go        # Main application file
├── internal/              # Private application code
│   ├── handler/           # HTTP request handlers
│   ├── model/             # Data models
│   ├── service/           # Business logic
│   └── server/            # Server configuration
└── pkg/                   # Public packages
    └── response/          # Response handling utilities
```

### 3️⃣ Technical Choices
- **Language**: Go 1.21+
- **Framework**: Gorilla Mux for routing
- **Storage**: In-memory data structures (as per requirements)
- **Concurrency**: Mutex for thread-safe operations
- **Error Handling**: Custom error types
- **Testing**: Unit and integration tests

## 🚀 Getting Started

### Prerequisites
- Go 1.21 or higher
- Git

### Installation
1. Clone the repository:
```bash
git clone https://github.com/grokkos/go-ledger.git
cd go-ledger
```

2. Install dependencies:
```bash
go mod download
```

3. Run the application:
```bash
go run cmd/api/main.go
```
The server will start on port `8080`.

## 🔗 API Endpoints

### 📌 Record Transaction
```http
POST /api/v1/transactions
Content-Type: application/json
{
    "type": "deposit"|"withdrawal",
    "amount": float
}
```

### 📌 Get Transactions
```http
GET /api/v1/transactions
```

### 📌 Get Balance
```http
GET /api/v1/balance
```

## 📊 Example API Usage

### 1️⃣ Make a deposit:
```bash
curl -X POST http://localhost:8080/api/v1/transactions \
-H "Content-Type: application/json" \
-d '{ "type": "deposit", "amount": 100.00 }'
```

### 2️⃣ Make a withdrawal:
```bash
curl -X POST http://localhost:8080/api/v1/transactions \
-H "Content-Type: application/json" \
-d '{ "type": "withdrawal", "amount": 50.00 }'
```

### 3️⃣ Check balance:
```bash
curl http://localhost:8080/api/v1/balance
```

### 4️⃣ View transaction history:
```bash
curl http://localhost:8080/api/v1/transactions
```

## 🛠️ Implementation Details

### Concurrency Handling
- Uses mutex locks to ensure thread-safe operations
- Read/Write mutex for optimized concurrent access

### Error Handling
- Custom error types for business logic errors
- HTTP status codes mapped to error types
- Consistent error response format

### Input Validation
- Transaction type validation
- Amount validation (must be positive)
- Balance checks for withdrawals

### Server Features
- Graceful shutdown
- Configurable timeouts
- API versioning

## 📌 Assumptions

1. **Data Persistence**: Stored in-memory, lost on restart
2. **Currency**: Assumes single currency system, using `float64`
3. **Transactions**: Immediate processing, no scheduled transactions
4. **Security**: No authentication/authorization as per requirements


## ✅ Testing

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test ./... -cover
```

Run integration tests:
```bash
go test ./internal/integration -v
```

### Test Structure
1. **Unit Tests**
    - `handler/handler_test.go`: Tests HTTP request handling
    - `service/service_test.go`: Tests business logic
    - `server/server_test.go`: Tests server configuration
    - `response/response_test.go`: Tests response utilities
2. **Integration Tests**
    - `internal/integration/integration_test.go`, testing complete application flow
    -  Example test flow
    - 
      1. Check initial balance (should be 0)
      2. Make a deposit
      3. Verify balance after deposit
      4. Make a withdrawal
      5. Verify final balance
      6. Check transaction history
      7. Test invalid withdrawal (insufficient funds)}`

### Test Coverage Goals
- Ensure all endpoints function correctly
- Verify business logic and constraints
- Test error handling and edge cases

## 🔮 Future Improvements

### 🔹 Technical Enhancements
- Add persistent storage
- Implement decimal handling for amounts
- Add rate limiting
- Implement authentication/authorization
- API documentation with Swagger

### 🔹 Operational Features
- Metrics collection
- Enhanced logging
- Health check endpoints
- Configuration management
- Docker support

