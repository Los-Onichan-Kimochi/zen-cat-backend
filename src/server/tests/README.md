# Test Suite Documentation

## Overview

This test suite follows the **GIVEN-WHEN-THEN** pattern from the marcadores247-backend-v2 project, providing comprehensive test coverage for both API and BLL (Business Logic Layer) components of the zen-cat-backend application.

## Structure

The test suite is organized into the following directories:

```
tests/
├── api/                    # API endpoint tests
│   ├── api_wrapper.go     # Main API test wrapper
│   ├── user/              # User API tests
│   ├── community/         # Community API tests
│   ├── plan/              # Plan API tests
│   ├── service/           # Service API tests
│   ├── login/             # Login API tests
│   ├── reservation/       # Reservation API tests
│   └── ...                # Other API modules
├── bll/
│   └── controller/        # BLL controller tests
│       ├── controller_wrapper.go  # Main BLL test wrapper
│       ├── user/          # User controller tests
│       ├── community/     # Community controller tests
│       ├── session/       # Session controller tests
│       └── ...            # Other controller modules
├── utils/                 # Test utilities
│   └── string.go          # String generation utilities
└── setup.go               # Test database setup utilities
```

## Test Pattern: GIVEN-WHEN-THEN

All tests follow the GIVEN-WHEN-THEN pattern for clarity and consistency:

```go
func TestExampleFunction(t *testing.T) {
    /*
        GIVEN: Initial conditions or setup
        WHEN:  The action being tested
        THEN:  Expected results or outcomes
    */
    // GIVEN
    // Set up test data and dependencies
    
    // WHEN
    // Execute the function or API call being tested
    
    // THEN
    // Assert the expected results
}
```

## Test Wrappers

### API Wrapper
The API wrapper (`api/api_wrapper.go`) provides:
- Database reset functionality for each test
- API server initialization
- Echo server setup for HTTP testing

### BLL Controller Wrapper
The BLL controller wrapper (`bll/controller/controller_wrapper.go`) provides:
- Individual controller test wrappers for each module
- Database reset functionality
- Logger mock initialization

## Test Categories

### 1. API Tests
API tests verify HTTP endpoints and test:
- **Success scenarios**: Valid requests with expected responses
- **Error scenarios**: Invalid requests, not found, unauthorized, etc.
- **Validation**: Request body validation, UUID validation
- **HTTP status codes**: 200, 201, 400, 404, 422, etc.

Example API test modules:
- `user/get_user_test.go` - Tests GET /user/{userId}
- `user/create_user_test.go` - Tests POST /user/
- `user/fetch_users_test.go` - Tests GET /user/

### 2. BLL Controller Tests
BLL tests verify business logic and test:
- **Data manipulation**: CRUD operations
- **Business rules**: Validation, constraints
- **Error handling**: Invalid data, missing dependencies
- **Data integrity**: Database state verification

Example BLL test modules:
- `user/get_user_test.go` - Tests User.GetUser()
- `user/create_user_test.go` - Tests User.CreateUser()
- `session/create_session_test.go` - Tests Session.CreateSession()

## Test Utilities

### String Utilities (`utils/string.go`)
- `GenerateRandomEmail()` - Generates unique test emails
- `GenerateRandomString(length)` - Generates random strings
- `GenerateMatchingString(substr)` - Generates strings containing substrings

### Helper Functions
Each test file includes common helpers:
- `strPtr(s string) *string` - Creates string pointers for optional fields

## Running Tests

### Run All Tests
```bash
go test ./src/server/tests/...
```

### Run API Tests Only
```bash
go test ./src/server/tests/api/...
```

### Run BLL Tests Only
```bash
go test ./src/server/tests/bll/...
```

### Run Specific Module Tests
```bash
# User API tests
go test ./src/server/tests/api/user/...

# User BLL tests
go test ./src/server/tests/bll/controller/user/...
```

### Run with Verbose Output
```bash
go test -v ./src/server/tests/...
```

## Test Data Management

### Database Setup
- Each test uses `ClearPostgresqlDatabase()` to ensure clean state
- Test data is created within each test for isolation
- Database transactions ensure no side effects between tests

### Test Environment
- Tests require a local PostgreSQL database
- Environment variables should point to test database
- Database should be different from development/production

## Coverage Areas

The test suite covers the following modules:

### API Endpoints
- ✅ User management (CRUD operations)
- ✅ Community management
- ✅ Plan management
- ✅ Service management
- ✅ Authentication (Login)
- ✅ Reservation management
- 🔄 Session management (in progress)
- 🔄 Professional management (in progress)
- 🔄 Local management (in progress)

### BLL Controllers
- ✅ User controller
- ✅ Community controller
- ✅ Session controller
- 🔄 Other controllers (in progress)

## Best Practices

1. **Test Isolation**: Each test should be independent and not rely on other tests
2. **Clear Naming**: Test function names should clearly describe what they test
3. **Comprehensive Coverage**: Test both success and failure scenarios
4. **Real Data**: Use realistic test data that mimics production scenarios
5. **Cleanup**: Always clean up test data and reset database state
6. **Assertions**: Use specific assertions rather than generic ones

## Example Test Structure

```go
func TestCreateUserSuccessfully(t *testing.T) {
    /*
        GIVEN: A valid user creation request
        WHEN:  POST /user/ is called with valid user data
        THEN:  A HTTP_201_CREATED status should be returned with the created user
    */
    // GIVEN
    server, db := apiTest.NewApiServerTestWrapper(t)
    
    createUserRequest := schemas.CreateUserRequest{
        Name:          "John",
        FirstLastName: "Doe",
        Email:         utilsTest.GenerateRandomEmail(),
        // ... other fields
    }
    
    // WHEN
    requestBody, _ := json.Marshal(createUserRequest)
    req := httptest.NewRequest(http.MethodPost, "/user/", bytes.NewBuffer(requestBody))
    req.Header.Set("Content-Type", "application/json")
    
    rec := httptest.NewRecorder()
    server.Echo.ServeHTTP(rec, req)
    
    // THEN
    assert.Equal(t, http.StatusCreated, rec.Code)
    
    var response schemas.User
    err := json.NewDecoder(rec.Body).Decode(&response)
    assert.NoError(t, err)
    assert.Equal(t, createUserRequest.Name, response.Name)
    // ... additional assertions
}
```

This test structure ensures maintainable, readable, and comprehensive test coverage for the zen-cat-backend application. 