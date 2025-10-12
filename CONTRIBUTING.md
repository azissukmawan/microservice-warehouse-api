# 🤝 Contributing to Warehouse & Point of Sales API

Thank you for your interest in contributing to this project! This repository participates in **Hacktoberfest 2025** 🎃

[![Hacktoberfest](https://img.shields.io/badge/hacktoberfest-2025-orange.svg)](https://hacktoberfest.com)

## 📋 Table of Contents

- [Getting Started](#-getting-started)
- [Hacktoberfest 2025](#-hacktoberfest-2025)
- [Types of Contributions](#-types-of-contributions)
- [Development Setup](#-development-setup)
- [Coding Guidelines](#-coding-guidelines)
- [Testing Requirements](#-testing-requirements)
- [Pull Request Process](#-pull-request-process)
- [Issue Guidelines](#-issue-guidelines)
- [Code Review Process](#-code-review-process)
- [Community Guidelines](#-community-guidelines)

## 🚀 Getting Started

### Prerequisites

Make sure you have:
- **Go** 1.24.3+
- **Docker** & **Docker Compose**
- **Git**
- **Text Editor** (VS Code recommended)

### Quick Setup

1. **Fork repository** on GitHub
2. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/micro-warehouse-api.git
   cd micro-warehouse-api
   ```

3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/ORIGINAL_OWNER/micro-warehouse-api.git
   ```

4. **Setup development environment** (see [Development Setup](#-development-setup))

### Labels for Hacktoberfest

Look for issues with labels:
- ![hacktoberfest](https://img.shields.io/badge/-hacktoberfest-orange) - Suitable for Hacktoberfest
- ![good first issue](https://img.shields.io/badge/-good%20first%20issue-green) - Good for beginners
- ![help wanted](https://img.shields.io/badge/-help%20wanted-blue) - Needs help
- ![documentation](https://img.shields.io/badge/-documentation-lightblue) - Documentation
- ![bug](https://img.shields.io/badge/-bug-red) - Bug fixes
- ![enhancement](https://img.shields.io/badge/-enhancement-purple) - New features

### Quality Standards for Hacktoberfest

To ensure your PR is accepted:
- ✅ **Meaningful contributions** - No spam or trivial changes
- ✅ **Follow coding standards** - According to project guidelines
- ✅ **Include tests** - If applicable
- ✅ **Update documentation** - If required
- ✅ **Descriptive commit messages** - Clear and informative

## 🔧 Types of Contributions

### 🐛 Bug Fixes

**What we need:**
- Fix existing bugs in any service
- Improve error handling
- Fix performance issues
- Resolve security vulnerabilities

**How to start:**
- Look for issues labeled `bug`
- Reproduce the bug locally
- Write tests that fail with the bug
- Fix the bug and ensure tests pass

### 📚 Documentation

**What we need:**
- Improve README sections
- Add code comments
- Create tutorials
- API documentation improvements
- Setup guides for different OS

**Good for beginners:** ✅

### ✨ New Features

**What we need:**
- Unit tests for existing code
- Integration tests
- Load testing
- Security testing
- Test utilities

**Examples:**
- Add logging middleware
- Implement caching
- Add validation
- Create health check endpoints
- Improve security

### 🧪 Testing

**What we need:**
- Unit tests for existing code
- Integration tests
- Load testing
- Test utilities and helpers

**Good for beginners:** ✅

### 🔒 Security

**What we need:**
- Security audit
- Vulnerability fixes
- Security middleware
- Input validation
- Authentication improvements

### 🚀 Performance

**What we need:**
- Performance optimizations
- Memory usage improvements
- Database query optimization
- Caching implementations

## 🛠️ Development Setup

### 1. Environment Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/micro-warehouse-api.git
cd micro-warehouse-api

# Setup environment variables
for service in api-gateway user-service product-service merchant-service transaction-service warehouse-service notification-service; do
  cp $service/env.example $service/.env
done
```

### 2. Start Infrastructure Services

```bash
# Start databases and supporting services
docker-compose up -d user_db product_db transaction_db merchant_db notification_db warehouse_db redis rabbitmq

# Wait for services to be healthy
docker-compose ps
```

### 3. Development Options

#### Option A: Docker Development (Recommended)
```bash
# Build and run all services
docker-compose up --build

# Or run specific service
docker-compose up --build user-service
```

#### Option B: Local Development
```bash
# Install dependencies for a service
cd user-service
go mod tidy

# Run service locally
go run main.go start

# With hot reload (install air first)
go install github.com/cosmtrek/air@latest
air
```

### 4. Verify Setup

```bash
# Check API Gateway health
curl http://localhost:8080/health

# Check individual services
curl http://localhost:8081/health  # User Service
curl http://localhost:8082/health  # Product Service
```

### Project Structure

Follow the established structure:
```
service-name/
├── main.go                 # Entry point
├── cmd/                   # CLI commands
├── app/                   # Application setup & DI
├── configs/               # Configuration
├── controller/            # HTTP handlers
├── database/              # Database connections
├── middleware/            # Custom middleware
├── model/                 # Database models
├── repository/            # Data access layer
├── usecase/              # Business logic
└── pkg/                  # Shared utilities
```

## 🧪 Testing Requirements

### Unit Tests

**Required for all new code:**
```go
func TestUserService_CreateUser(t *testing.T) {
    // Setup
    repo := &mockUserRepository{}
    service := NewUserService(repo)

    // Test cases
    tests := []struct {
        name    string
        req     CreateUserRequest
        want    *User
        wantErr bool
    }{
        {
            name: "valid user creation",
            req: CreateUserRequest{
                Name:  "John Doe",
                Email: "john@example.com",
            },
            want: &User{
                ID:    1,
                Name:  "John Doe",
                Email: "john@example.com",
            },
            wantErr: false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := service.CreateUser(context.Background(), tt.req)
            if (err != nil) != tt.wantErr {
                t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            // Assert results...
        })
    }
}
```

### Running Tests

```bash
# Run tests for specific service
cd user-service
go test ./...

# Run with coverage
go test -cover ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -run TestUserService_CreateUser

# Run integration tests
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```


### Integration Tests Example

```go
func TestUserAPI_CreateUser_Integration(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)

    // Setup test server
    app := setupTestApp(t, db)

    // Test data
    reqBody := `{
        "name": "John Doe",
        "email": "john@example.com",
        "password": "password123"
    }`

    // Make request
    req := httptest.NewRequest("POST", "/api/v1/users", strings.NewReader(reqBody))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req)
    require.NoError(t, err)
    require.Equal(t, 201, resp.StatusCode)

    // Verify response
    var user User
    err = json.NewDecoder(resp.Body).Decode(&user)
    require.NoError(t, err)
    assert.Equal(t, "John Doe", user.Name)
    assert.Equal(t, "john@example.com", user.Email)
}
```

## 📝 Issue Guidelines

### Creating Good Issues

#### Bug Reports
```markdown
**Describe the bug**
A clear and concise description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Start the service with '...'
2. Make request to '....'
3. See error

**Expected behavior**
A clear description of what you expected to happen.

**Screenshots/Logs**
If applicable, add screenshots or log outputs.

**Environment:**
- OS: [e.g. Windows, macOS, Linux]
- Go version: [e.g. 1.24.3]
- Docker version: [e.g. 20.10.8]

**Additional context**
Add any other context about the problem here.
```

#### Feature Requests
```markdown
**Is your feature request related to a problem?**
A clear description of what the problem is.

**Describe the solution you'd like**
A clear description of what you want to happen.

**Describe alternatives you've considered**
Other solutions you've considered.

**Additional context**
Add any other context or screenshots about the feature request.
```

### Issue Labels

We use these labels to categorize issues:

- ![bug](https://img.shields.io/badge/-bug-red) - Something isn't working
- ![enhancement](https://img.shields.io/badge/-enhancement-purple) - New feature or request
- ![documentation](https://img.shields.io/badge/-documentation-lightblue) - Documentation improvements
- ![good first issue](https://img.shields.io/badge/-good%20first%20issue-green) - Good for newcomers
- ![help wanted](https://img.shields.io/badge/-help%20wanted-blue) - Extra attention is needed
- ![hacktoberfest](https://img.shields.io/badge/-hacktoberfest-orange) - Hacktoberfest eligible
- ![priority: high](https://img.shields.io/badge/-priority:%20high-red) - High priority
- ![priority: low](https://img.shields.io/badge/-priority:%20low-green) - Low priority


## 📞 Contact

Need help or have questions?

- **Issues**: [Create an issue](https://github.com/your-username/micro-warehouse-api/issues/new)

---

**Happy Contributing! 🚀**

> Remember: Every contribution, no matter how small, makes a difference. We appreciate your time and effort in making this project better for everyone.
