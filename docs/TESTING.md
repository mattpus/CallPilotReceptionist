# Testing Documentation

## Overview

This document describes the testing strategy, test organization, and how to run tests for the Vapi AI Integration Backend.

## Testing Philosophy

### Goals
1. **High Coverage**: Target 80%+ code coverage
2. **Fast Feedback**: Unit tests should run in < 1 second
3. **Isolation**: Tests should not depend on external services
4. **Maintainability**: Tests should be easy to understand and modify
5. **Reliability**: Tests should be deterministic and not flaky

### Testing Pyramid

```
        /\
       /  \
      / E2E \       <- Few, slow, high-value
     /______\
    /        \
   / Integration \  <- Some, moderate speed
  /______________\
 /                \
/   Unit Tests     \ <- Many, fast, focused
/____________________\
```

## Test Organization

### Directory Structure

```
internal/
├── application/
│   └── services/
│       ├── auth_service.go
│       ├── auth_service_test.go
│       ├── call_service.go
│       ├── mocks_test.go          <- Shared mocks for testing
│       └── ...
├── domain/
│   └── entities/
│       ├── entities.go
│       └── entities_test.go       <- Entity validation tests
└── api/
    └── handlers/
        ├── auth_handler.go
        └── auth_handler_test.go   <- Handler integration tests
```

### Test Files

| File | Purpose | Status |
|------|---------|--------|
| `auth_service_test.go` | AuthService unit tests | ✅ Complete |
| `entities_test.go` | Domain entity validation tests | ✅ Complete |
| `mocks_test.go` | Shared test mocks (repositories, providers) | ✅ Complete |
| `business_service_test.go` | BusinessService unit tests | ⚠️ To be added |
| `call_service_test.go` | CallService unit tests | ⚠️ To be added |
| `interaction_service_test.go` | InteractionService unit tests | ⚠️ To be added |
| `analytics_service_test.go` | AnalyticsService unit tests | ⚠️ To be added |

## Running Tests

### All Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Specific Package Tests

```bash
# Test services only
go test ./internal/application/services/

# Test entities only
go test ./internal/domain/entities/

# Test with race detector
go test -race ./...

# Run tests in parallel
go test -parallel 4 ./...
```

### Watch Mode

Use `entr` for continuous testing:

```bash
# Install entr (macOS)
brew install entr

# Watch Go files and run tests on change
find . -name "*.go" | entr -c go test ./...
```

## Unit Tests

### Service Tests

**Location**: `internal/application/services/*_test.go`

**Approach**:
- Use mock repositories and providers
- Test business logic in isolation
- Table-driven tests for multiple scenarios
- Test both success and error paths

**Example**:

```go
func TestAuthService_Register(t *testing.T) {
    tests := []struct {
        name          string
        request       dto.RegisterRequest
        mockUserRepo  *mockUserRepository
        mockBizRepo   *mockBusinessRepository
        expectedError bool
    }{
        {
            name: "successful registration",
            request: dto.RegisterRequest{
                Email: "test@example.com",
                Password: "SecurePass123!",
                BusinessName: "Test Business",
            },
            // ... mocks ...
            expectedError: false,
        },
        // ... more test cases ...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            service := NewAuthService(tt.mockUserRepo, tt.mockBizRepo, ...)
            _, err := service.Register(context.Background(), tt.request)
            // ... assertions ...
        })
    }
}
```

### Entity Tests

**Location**: `internal/domain/entities/entities_test.go`

**Tests**:
- Field validation (required, format, length)
- Business rule enforcement
- State transitions
- Edge cases

**Example**:

```go
func TestBusiness_Validate(t *testing.T) {
    tests := []struct {
        name          string
        business      *Business
        expectedError bool
        errorContains string
    }{
        {
            name: "valid business",
            business: &Business{
                Name: "Test Business",
                Phone: "+1234567890",
            },
            expectedError: false,
        },
        {
            name: "missing name",
            business: &Business{
                Phone: "+1234567890",
            },
            expectedError: true,
            errorContains: "name is required",
        },
    }
    // ... test execution ...
}
```

## Integration Tests

### Handler Tests

**Location**: `internal/api/handlers/*_test.go`

**Approach**:
- Use `httptest.NewRecorder` for HTTP testing
- Mock services (not repositories)
- Test complete request/response cycle
- Verify middleware behavior

**Example**:

```go
func TestAuthHandler_Register(t *testing.T) {
    mockService := &mockAuthService{}
    handler := NewAuthHandler(mockService)
    
    reqBody := `{"email":"test@example.com","password":"pass123"}`
    req := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(reqBody))
    w := httptest.NewRecorder()
    
    handler.Register(w, req)
    
    assert.Equal(t, http.StatusCreated, w.Code)
    // ... more assertions ...
}
```

## Test Mocks

### Shared Mocks

**Location**: `internal/application/services/mocks_test.go`

**Mocks Available**:
- `mockVoiceProvider` - Voice AI provider mock
- `mockUserRepository` - User data access mock
- `mockBusinessRepository` - Business data access mock
- `mockCallRepository` - Call data access mock
- `mockTranscriptRepository` - Transcript data access mock
- `mockInteractionRepository` - Interaction data access mock
- `mockAppointmentRepository` - Appointment data access mock

**Usage**:

```go
mockRepo := &mockUserRepository{
    getByEmailFunc: func(ctx context.Context, email string) (*entities.User, error) {
        return &entities.User{
            ID: uuid.New(),
            Email: email,
        }, nil
    },
}

service := NewAuthService(mockRepo, ...)
```

### Mock Best Practices

1. **Return Real Data**: Use actual entity structs, not nil
2. **Test Error Paths**: Mock repository errors to test error handling
3. **Verify Calls**: Track method calls if needed
4. **Keep Simple**: Don't add complexity to mocks
5. **One Mock Per Test**: Don't reuse mocks across tests

## Testing Patterns

### Table-Driven Tests

**Why**: Test multiple scenarios with less boilerplate

```go
tests := []struct {
    name     string
    input    string
    expected string
}{
    {"scenario 1", "input1", "output1"},
    {"scenario 2", "input2", "output2"},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        result := function(tt.input)
        assert.Equal(t, tt.expected, result)
    })
}
```

### Test Fixtures

**Why**: Reusable test data

```go
func newTestUser() *entities.User {
    return &entities.User{
        ID:           uuid.New(),
        Email:        "test@example.com",
        PasswordHash: "$2a$10$...",
        Role:         "owner",
    }
}
```

### Context Usage

**Always pass context**:

```go
func TestService_Method(t *testing.T) {
    ctx := context.Background()
    result, err := service.Method(ctx, ...)
    // ...
}
```

## Coverage Goals

### Target Coverage by Package

| Package | Target | Current |
|---------|--------|---------|
| `services` | 85% | ~60% |
| `entities` | 90% | 100% ✅ |
| `handlers` | 80% | 0% |
| `middleware` | 75% | 0% |
| `providers` | 70% | 0% |

### Viewing Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View in browser
go tool cover -html=coverage.out

# View summary
go tool cover -func=coverage.out
```

### Coverage Report Example

```
github.com/vapiAIIntegration/internal/application/services/auth_service.go:38:    Register        87.5%
github.com/vapiAIIntegration/internal/application/services/auth_service.go:116:   Login           91.2%
github.com/vapiAIIntegration/internal/application/services/auth_service.go:167:   GenerateTokens  100.0%
total:                                                                             82.3%
```

## Test Data Management

### Test Database

For integration tests requiring a database:

```bash
# Use test database
export DATABASE_URL="postgresql://user:pass@localhost:5432/test_db"

# Or use docker-compose
docker-compose -f docker-compose.test.yml up -d
```

### Cleanup

```go
func TestWithDB(t *testing.T) {
    db := setupTestDB(t)
    defer teardownTestDB(t, db)
    
    // ... test code ...
}

func setupTestDB(t *testing.T) *sql.DB {
    // Create test database, run migrations
}

func teardownTestDB(t *testing.T, db *sql.DB) {
    // Clean up test data, close connections
}
```

## Continuous Integration

### GitHub Actions Example

```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.21
      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...
      - name: Upload coverage
        uses: codecov/codecov-action@v2
```

## Best Practices

### DO
✅ Test public APIs, not private functions  
✅ Use table-driven tests for multiple scenarios  
✅ Mock external dependencies  
✅ Test error cases thoroughly  
✅ Keep tests focused and simple  
✅ Use descriptive test names  
✅ Clean up resources (defer statements)

### DON'T
❌ Test implementation details  
❌ Share state between tests  
❌ Use time.Sleep() for synchronization  
❌ Ignore test failures  
❌ Skip writing tests for "simple" code  
❌ Use production database for tests  
❌ Make tests dependent on execution order

## Debugging Tests

### Verbose Output

```bash
go test -v ./internal/application/services/
```

### Run Single Test

```bash
go test -v -run TestAuthService_Register ./internal/application/services/
```

### Debug with Delve

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug test
dlv test ./internal/application/services/ -- -test.run TestAuthService_Register
```

## Next Steps

### Priority 1: Complete Service Tests
- [ ] CallService tests (webhook handling, transcript fetch)
- [ ] BusinessService tests (CRUD operations)
- [ ] InteractionService tests (appointment extraction)
- [ ] AnalyticsService tests (statistics, trends)

### Priority 2: Add Handler Tests
- [ ] AuthHandler integration tests
- [ ] CallHandler integration tests
- [ ] BusinessHandler integration tests

### Priority 3: Infrastructure Tests
- [ ] VapiProvider tests (mock HTTP server)
- [ ] Repository integration tests (test database)
- [ ] Middleware tests (auth, logging, error handling)

### Priority 4: E2E Tests
- [ ] Complete registration → login → call flow
- [ ] Webhook processing flow
- [ ] Dashboard data retrieval

## Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Table-Driven Tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Go Test Comments](https://golang.org/wiki/TableDrivenTests)
- [Testify Framework](https://github.com/stretchr/testify) (optional)

---

**Last Updated**: 2024  
**Test Coverage**: ~65%  
**Target Coverage**: 80%+
