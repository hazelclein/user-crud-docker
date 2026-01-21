# 🚀 User CRUD API - Complete Documentation

[![Go Version](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker)](https://www.docker.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?logo=postgresql)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7-DC382D?logo=redis)](https://redis.io/)
[![Jaeger](https://img.shields.io/badge/Jaeger-Tracing-66CFE3)](https://www.jaegertracing.io/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

RESTful API untuk manajemen user dengan arsitektur **Clean Architecture**, **Domain-Driven Design (DDD)**, dan **CQRS Pattern**. Dibangun dengan Go, PostgreSQL, Redis, dan dilengkapi dengan distributed tracing menggunakan Jaeger.

---

## 📋 Table of Contents

- [Features](#-features)
- [Tech Stack](#-tech-stack)
- [Architecture](#-architecture)
- [Prerequisites](#-prerequisites)
- [Installation](#-installation)
- [Configuration](#-configuration)
- [API Documentation](#-api-documentation)
- [Examples](#-examples)
- [Development](#-development)
- [Testing](#-testing)
- [Monitoring](#-monitoring)
- [Troubleshooting](#-troubleshooting)
- [Contributing](#-contributing)
- [License](#-license)

---

## ✨ Features

### **Core Features**
- ✅ **CRUD Operations** - Create, Read, Update, Delete users
- ✅ **Password Management** - Secure bcrypt password hashing
- ✅ **Change Password** - Users can change their own password
- ✅ **Advanced Search** - Search users by name or email
- ✅ **Flexible Filtering** - Filter by age range, sort by multiple fields
- ✅ **Pagination** - Efficient data loading with customizable page size
- ✅ **Input Validation** - Comprehensive request validation
- ✅ **Error Handling** - Consistent error response format

### **Technical Features**
- ✅ **Clean Architecture** - Clear separation of concerns
- ✅ **Domain-Driven Design** - Business logic centered in domain layer
- ✅ **CQRS Pattern** - Separate Command (Write) and Query (Read) operations
- ✅ **Redis Caching** - Fast data access with Redis cache
- ✅ **Distributed Tracing** - Request tracing with Jaeger
- ✅ **Rate Limiting** - Protect API from abuse
- ✅ **Circuit Breaker** - Fault tolerance and resilience
- ✅ **Health Check** - Monitor application and dependencies status
- ✅ **Graceful Shutdown** - Clean shutdown handling
- ✅ **Docker Support** - Full containerization with Docker Compose
- ✅ **Database Migrations** - Automatic schema management
- ✅ **Database Indexes** - Optimized query performance

---

## 🛠 Tech Stack

| Technology | Purpose | Version |
|------------|---------|---------|
| **Go** | Backend language | 1.23+ |
| **Gin** | Web framework | Latest |
| **PostgreSQL** | Primary database | 16 |
| **Redis** | Caching layer | 7 |
| **pgx** | PostgreSQL driver | v5 |
| **Jaeger** | Distributed tracing | Latest |
| **Docker** | Containerization | Latest |
| **Docker Compose** | Orchestration | Latest |

---

## 🏗 Architecture

```
user-crud/
├── cmd/
│   └── api/
│       └── main.go                    # Application entry point
│
├── internal/
│   ├── application/                   # Application layer (Use Cases)
│   │   ├── command/                   # Write operations (CQRS)
│   │   │   ├── create_user.go
│   │   │   ├── update_user.go
│   │   │   ├── delete_user.go
│   │   │   └── change_password.go
│   │   └── query/                     # Read operations (CQRS)
│   │       ├── get_user.go
│   │       └── list_users.go
│   │
│   ├── domain/                        # Domain layer (Business Logic)
│   │   ├── user.go                    # User entity
│   │   └── repository.go              # Repository interface
│   │
│   ├── infrastructure/                # Infrastructure layer
│   │   ├── http/
│   │   │   ├── handler/               # HTTP handlers
│   │   │   ├── middleware/            # Middlewares
│   │   │   └── router/                # Route definitions
│   │   ├── persistence/               # Database implementation
│   │   │   ├── database.go
│   │   │   └── postgres_user_repository.go
│   │   ├── cache/                     # Redis implementation
│   │   └── tracing/                   # Jaeger tracing
│   │
│   └── config/                        # Configuration
│       └── config.go
│
├── migrations/                        # Database migrations
│   └── 001_create_users_table.sql
│
├── docker-compose.yml                 # Docker orchestration
├── Dockerfile                         # App container definition
├── go.mod                             # Go dependencies
└── README.md                          # This file
```

### **Architecture Principles**

#### **1. Clean Architecture Layers**
```
Presentation → Application → Domain ← Infrastructure
     ↓              ↓           ↓            ↓
   HTTP        Use Cases    Business      Database
  Handlers                   Logic        Redis/etc
```

#### **2. CQRS Pattern**
- **Commands** (Write): Create, Update, Delete, ChangePassword
- **Queries** (Read): Get, List, Search

#### **3. Dependency Rule**
- Inner layers don't know about outer layers
- Dependencies point inward
- Domain layer is independent

---

## 📦 Prerequisites

Before you begin, ensure you have:

- **Docker** (20.10+) - [Install Docker](https://docs.docker.com/get-docker/)
- **Docker Compose** (2.0+) - Usually comes with Docker Desktop
- **Git** - [Install Git](https://git-scm.com/downloads)
- **curl** or **Postman** - For API testing

**For local development without Docker:**
- **Go** (1.23+) - [Install Go](https://go.dev/doc/install)
- **PostgreSQL** (16+) - [Install PostgreSQL](https://www.postgresql.org/download/)
- **Redis** (7+) - [Install Redis](https://redis.io/download)

---

## 🚀 Installation

### **Method 1: Using Docker Compose (Recommended)**

This is the easiest way to get started. All dependencies are automatically set up.

```bash
# 1. Clone the repository
git clone https://github.com/hazelclein/user-crud-docker.git
cd user-crud-docker

# 2. Build and start all services
docker-compose up --build -d

# 3. Check if all containers are running
docker-compose ps

# 4. View application logs
docker-compose logs -f app

# 5. Test the API
curl http://localhost:8080/health
```

**Expected output:**
```json
{
  "status": "healthy",
  "database": "connected",
  "timestamp": "2026-01-21T10:00:00Z"
}
```

### **Method 2: Local Development**

For development without Docker:

```bash
# 1. Start PostgreSQL and Redis with Docker
docker-compose up -d postgres redis jaeger

# 2. Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=userdb
export SERVER_PORT=8080
export REDIS_HOST=localhost
export REDIS_PORT=6379
export JAEGER_ENDPOINT=http://localhost:14268/api/traces

# 3. Install Go dependencies
go mod download

# 4. Run the application
go run cmd/api/main.go
```

### **Stopping Services**

```bash
# Stop all services
docker-compose down

# Stop and remove volumes (database data will be lost)
docker-compose down -v

# Stop and remove orphan containers
docker-compose down --remove-orphans
```

---

## ⚙️ Configuration

### **Environment Variables**

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `postgres` | PostgreSQL hostname |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | Database username |
| `DB_PASSWORD` | `postgres` | Database password |
| `DB_NAME` | `userdb` | Database name |
| `SERVER_PORT` | `8080` | HTTP server port |
| `REDIS_HOST` | `redis` | Redis hostname |
| `REDIS_PORT` | `6379` | Redis port |
| `JAEGER_ENDPOINT` | `http://jaeger:14268/api/traces` | Jaeger collector endpoint |

### **Docker Compose Configuration**

The `docker-compose.yml` file defines four services:

1. **postgres** - PostgreSQL database
2. **redis** - Redis cache
3. **jaeger** - Distributed tracing
4. **app** - Your Go application

### **Ports Mapping**

| Service | Internal Port | External Port | Purpose |
|---------|---------------|---------------|---------|
| App | 8080 | 8080 | REST API |
| PostgreSQL | 5432 | 5432 | Database |
| Redis | 6379 | 6379 | Cache |
| Jaeger UI | 16686 | 16686 | Tracing UI |
| Jaeger Collector | 14268 | 14268 | Trace collector |

---

## 📡 API Documentation

### **Base URL**
```
http://localhost:8080
```

### **Response Format**

#### Success Response
```json
{
  "status": "success",
  "data": { ... }
}
```

#### Error Response
```json
{
  "status": "error",
  "message": "Error description"
}
```

#### Paginated Response
```json
{
  "status": "success",
  "data": [ ... ],
  "total": 100,
  "page": 1,
  "limit": 10,
  "total_pages": 10
}
```

---

### **Endpoints**

#### **1. Health Check**

Check application and database status.

```http
GET /health
```

**Response:** `200 OK`
```json
{
  "status": "healthy",
  "database": "connected",
  "timestamp": "2026-01-21T10:00:00Z"
}
```

---

#### **2. Create User**

Create a new user with encrypted password.

```http
POST /api/v1/users
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "age": 30
}
```

**Validation Rules:**
- `name`: required, string, 2-100 characters
- `email`: required, valid email format, unique
- `password`: required, minimum 8 characters
- `age`: required, integer, 0-150

**Response:** `201 Created`
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30,
    "created_at": "2026-01-21T10:00:00Z",
    "updated_at": "2026-01-21T10:00:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request` - Validation error
- `409 Conflict` - Email already exists
- `500 Internal Server Error` - Server error

---

#### **3. Get User by ID**

Retrieve a specific user by their ID.

```http
GET /api/v1/users/:id
```

**Path Parameters:**
- `id` (integer, required) - User ID

**Example:**
```http
GET /api/v1/users/1
```

**Response:** `200 OK`
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30,
    "created_at": "2026-01-21T10:00:00Z",
    "updated_at": "2026-01-21T10:00:00Z"
  }
}
```

**Error Responses:**
- `404 Not Found` - User not found

---

#### **4. List Users**

Get a list of users with optional filtering, sorting, and pagination.

```http
GET /api/v1/users?search=john&age_min=25&age_max=40&sort=name&order=asc&page=1&limit=10
```

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `search` | string | - | Search by name or email (case-insensitive) |
| `age_min` | integer | - | Minimum age filter |
| `age_max` | integer | - | Maximum age filter |
| `sort` | string | `id` | Sort field: `id`, `name`, `email`, `age`, `created_at` |
| `order` | string | `asc` | Sort order: `asc` or `desc` |
| `page` | integer | `1` | Page number (starts from 1) |
| `limit` | integer | `10` | Items per page (max: 100) |

**Examples:**

```bash
# Get all users
GET /api/v1/users

# Search by name
GET /api/v1/users?search=john

# Filter by age range
GET /api/v1/users?age_min=25&age_max=35

# Sort by age (descending)
GET /api/v1/users?sort=age&order=desc

# Pagination
GET /api/v1/users?page=2&limit=20

# Combined filters
GET /api/v1/users?search=example&age_min=25&sort=name&order=asc&page=1&limit=10
```

**Response:** `200 OK`
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "age": 30,
      "created_at": "2026-01-21T10:00:00Z",
      "updated_at": "2026-01-21T10:00:00Z"
    }
  ],
  "total": 50,
  "page": 1,
  "limit": 10,
  "total_pages": 5
}
```

---

#### **5. Search Users**

Dedicated search endpoint for finding users by keyword.

```http
GET /api/v1/users/search?q=john&page=1&limit=10
```

**Query Parameters:**
- `q` (string, required) - Search keyword
- `page` (integer, optional) - Page number
- `limit` (integer, optional) - Items per page

**Response:** `200 OK`
```json
{
  "status": "success",
  "data": [ ... ],
  "total": 5,
  "page": 1,
  "limit": 10,
  "total_pages": 1
}
```

---

#### **6. Update User**

Update user information (excluding password).

```http
PUT /api/v1/users/:id
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "John Updated",
  "email": "john.updated@example.com",
  "age": 31
}
```

**Note:** Password cannot be changed via this endpoint. Use Change Password endpoint instead.

**Response:** `200 OK`
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "name": "John Updated",
    "email": "john.updated@example.com",
    "age": 31,
    "created_at": "2026-01-21T10:00:00Z",
    "updated_at": "2026-01-21T10:15:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request` - Validation error
- `404 Not Found` - User not found
- `409 Conflict` - Email already exists

---

#### **7. Change Password**

Change user's password with validation.

```http
PUT /api/v1/users/:id/change-password
Content-Type: application/json
```

**Request Body:**
```json
{
  "old_password": "password123",
  "new_password": "newpassword456"
}
```

**Validation Rules:**
- `old_password`: required, must match current password
- `new_password`: required, minimum 8 characters, must be different from old password

**Response:** `200 OK`
```json
{
  "status": "success",
  "message": "password changed successfully"
}
```

**Error Responses:**
- `400 Bad Request` - Validation error or incorrect old password
- `404 Not Found` - User not found

---

#### **8. Delete User**

Permanently delete a user.

```http
DELETE /api/v1/users/:id
```

**Response:** `200 OK`
```json
{
  "status": "success",
  "message": "user deleted successfully"
}
```

**Error Responses:**
- `404 Not Found` - User not found

---

## 💡 Examples

### **Using cURL**

#### Create a User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Smith",
    "email": "alice@example.com",
    "password": "securepass123",
    "age": 28
  }'
```

#### Get All Users
```bash
curl http://localhost:8080/api/v1/users
```

#### Search Users
```bash
curl "http://localhost:8080/api/v1/users?search=alice&sort=name&order=asc"
```

#### Update User
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Johnson",
    "email": "alice.j@example.com",
    "age": 29
  }'
```

#### Change Password
```bash
curl -X PUT http://localhost:8080/api/v1/users/1/change-password \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "securepass123",
    "new_password": "newsecurepass456"
  }'
```

#### Delete User
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

---

### **Using JavaScript (Fetch API)**

```javascript
// Create User
const createUser = async () => {
  const response = await fetch('http://localhost:8080/api/v1/users', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      name: 'Bob Wilson',
      email: 'bob@example.com',
      password: 'password123',
      age: 35
    })
  });
  
  const data = await response.json();
  console.log(data);
};

// Get Users with Filters
const getUsers = async () => {
  const params = new URLSearchParams({
    search: 'example',
    age_min: 25,
    age_max: 40,
    sort: 'name',
    order: 'asc',
    page: 1,
    limit: 10
  });
  
  const response = await fetch(`http://localhost:8080/api/v1/users?${params}`);
  const data = await response.json();
  console.log(data);
};

// Update User
const updateUser = async (id) => {
  const response = await fetch(`http://localhost:8080/api/v1/users/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      name: 'Bob Wilson Updated',
      email: 'bob.updated@example.com',
      age: 36
    })
  });
  
  const data = await response.json();
  console.log(data);
};
```

---

### **Using Python (requests)**

```python
import requests

BASE_URL = "http://localhost:8080/api/v1"

# Create User
def create_user():
    response = requests.post(f"{BASE_URL}/users", json={
        "name": "Charlie Brown",
        "email": "charlie@example.com",
        "password": "password123",
        "age": 42
    })
    print(response.json())

# Get Users with Filters
def get_users():
    params = {
        "search": "example",
        "age_min": 25,
        "age_max": 50,
        "sort": "age",
        "order": "desc",
        "page": 1,
        "limit": 20
    }
    response = requests.get(f"{BASE_URL}/users", params=params)
    print(response.json())

# Update User
def update_user(user_id):
    response = requests.put(f"{BASE_URL}/users/{user_id}", json={
        "name": "Charlie Brown Updated",
        "email": "charlie.updated@example.com",
        "age": 43
    })
    print(response.json())

# Change Password
def change_password(user_id):
    response = requests.put(f"{BASE_URL}/users/{user_id}/change-password", json={
        "old_password": "password123",
        "new_password": "newpassword456"
    })
    print(response.json())

# Delete User
def delete_user(user_id):
    response = requests.delete(f"{BASE_URL}/users/{user_id}")
    print(response.json())
```

---

## 🧪 Testing

### **Manual Testing Script**

Create a test script to populate sample data:

```bash
#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"

echo "Creating test users..."

# Create 5 test users
curl -X POST $BASE_URL/users -H "Content-Type: application/json" -d '{"name":"John Doe","email":"john@example.com","password":"password123","age":25}'
curl -X POST $BASE_URL/users -H "Content-Type: application/json" -d '{"name":"Jane Smith","email":"jane@example.com","password":"password123","age":30}'
curl -X POST $BASE_URL/users -H "Content-Type: application/json" -d '{"name":"Bob Johnson","email":"bob@test.com","password":"password123","age":22}'
curl -X POST $BASE_URL/users -H "Content-Type: application/json" -d '{"name":"Alice Brown","email":"alice@example.com","password":"password123","age":28}'
curl -X POST $BASE_URL/users -H "Content-Type: application/json" -d '{"name":"Charlie Wilson","email":"charlie@test.com","password":"password123","age":35}'

echo ""
echo "Test users created successfully!"
echo ""
echo "Testing search..."
curl "$BASE_URL/users?search=example"

echo ""
echo "Testing filter by age..."
curl "$BASE_URL/users?age_min=25&age_max=30"
```

Save as `test.sh` and run:
```bash
chmod +x test.sh
./test.sh
```

### **Test Scenarios**

#### 1. **Validation Tests**

```bash
# Empty password (should fail)
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","email":"test@test.com","password":"","age":25}'

# Password too short (should fail)
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","email":"test@test.com","password":"12345","age":25}'

# Invalid email (should fail)
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","email":"not-an-email","password":"password123","age":25}'

# Age out of range (should fail)
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","email":"test@test.com","password":"password123","age":200}'
```

#### 2. **Duplicate Email Test**

```bash
# Create first user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"duplicate@test.com","password":"password123","age":25}'

# Try to create user with same email (should fail with 409 Conflict)
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Another User","email":"duplicate@test.com","password":"password123","age":30}'
```

#### 3. **Pagination Test**

```bash
# Get page 1 (2 items per page)
curl "http://localhost:8080/api/v1/users?page=1&limit=2"

# Get page 2
curl "http://localhost:8080/api/v1/users?page=2&limit=2"

# Get page 3
curl "http://localhost:8080/api/v1/users?page=3&limit=2"
```

#### 4. **Sorting Test**

```bash
# Sort by name (A-Z)
curl "http://localhost:8080/api/v1/users?sort=name&order=asc"

# Sort by age (oldest first)
curl "http://localhost:8080/api/v1/users?sort=age&order=desc"

# Sort by creation date (newest first)
curl "http://localhost:8080/api/v1/users?sort=created_at&order=desc"
```

---

## 📊 Monitoring

### **Jaeger Tracing UI**

Access the Jaeger UI to view distributed traces:

```
http://localhost:16686
```

**Features:**
- View request traces
- Analyze performance bottlenecks
- Debug issues in distributed systems
- Monitor service dependencies

### **Application Logs**

View real-time logs:

```bash
# All services
docker-compose logs -f

# Only app logs
docker-compose logs -f app

# Only database logs
docker-compose logs -f postgres

# Last 100 lines
docker-compose logs --tail=100 app
```

### **Container Health**

Check container status:

```bash
# List all containers
docker-compose ps

# Check resource usage
docker stats

# Inspect specific container
docker inspect user_crud_app
```

---

## 🔧 Development

### **Project Structure Explained**

```
internal/
├── application/          # Application Logic (Use Cases)
│   ├── command/          # Write operations
│   │   └── *.go          # CreateUser, UpdateUser, etc.
│   └── query/            # Read operations
│       └── *.go          # GetUser, ListUsers, etc.
│
├── domain/               # Business Logic (Core)
│   ├── user.go           # User entity with business rules
│   └── repository.go     # Repository interface (contract)
│
├── infrastructure/       # External Dependencies
│   ├── http/             # HTTP layer
│   │   ├── handler/      # HTTP request handlers
│   │   ├── middleware/   # Rate limiting, tracing, etc.
│   │   └── router/       # Route definitions
│   ├── persistence/      # Database layer
│   │   ├── database.go   # Database connection
│   │   └── *_repository.go  # Repository implementations
│   ├── cache/            # Redis layer
│   └── tracing/          # Jaeger configuration
│
└── config/               # Configuration management
    └── config.go         # Environment variables
```

### **Adding New Features**

#### Example: Add "GetUserByEmail" endpoint

**1. Add to Domain Repository Interface**
```go
// internal/domain/repository.go
type UserRepository interface {
    // ... existing methods
    GetByEmail(ctx context.Context, email string) (*User, error)
}
```

**2. Implement in Repository**
```go
// internal/infrastructure/persistence/postgres_user_repository.go
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    // Implementation here
}
```

**3. Create Query Handler**
```go
// internal/application/query/get_user_by_email.go
type GetUserByEmailHandler struct {
    repo domain.UserRepository
}

func (h *GetUserByEmailHandler) Handle(ctx context.Context, email string) (*domain.User, error) {
    return h.repo.GetByEmail(ctx, email)
}
```

**4. Add HTTP Handler**
```go
// internal/infrastructure/http/handler/handler.go
func (h *Handler) GetUserByEmail(c *gin.Context) {
    email := c.Param("email")
    user, err := h.getUserByEmailHandler.Handle(c.Request.Context(), email)
    // Handle response
}
```

**5. Register Route**
```go
// internal/infrastructure/http/router/router.go
users.GET("/email/:email", h.GetUserByEmail)
```

### **Database Migrations**

Add new migrations in `migrations/` folder:

```sql
-- migrations/002_add_phone_column.sql
ALTER TABLE users ADD COLUMN phone VARCHAR(20);
```

### **Hot Reload for Development**

Use Air for automatic reload on code changes:

```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

---

## 🐛 Troubleshooting

### **Common Issues**

#### 1. **Database Connection Failed**

**Error:**
```
Failed to connect to database: dial tcp 127.0.0.1:5432: connect: connection refused
```

**Solution:**
- Check if PostgreSQL container is running: `docker-compose ps`
- Verify `DB_HOST` is set to `postgres` (not `localhost`)
- Restart containers: `docker-compose restart`

#### 2. **Port Already in Use**

**Error:**
```
Bind for 0.0.0.0:8080 failed: port is already allocated
```

**Solution:**
```bash
# Find process using port 8080
lsof -i :8080  # Mac/Linux
netstat -ano | findstr :8080  # Windows

# Kill the process or change SERVER_PORT in docker-compose.yml
```

#### 3. **Migration Errors**

**Error:**
```
Migrations failed: table already exists
```

**Solution:**
```bash
# Reset database
docker-compose down -v
docker-compose up -d
```

#### 4. **Orphan Containers**

**Error:**
```
Found orphan containers (wizardly_payne) for this project
```

**Solution:**
```bash
# Remove orphan containers
docker-compose down --remove-orphans
docker-compose up -d
```

#### 5. **Redis Connection Failed**

**Solution:**
```bash
# Check Redis container
docker-compose logs redis

# Restart Redis
docker-compose restart redis
```

---

## 🔒 Security Features

### **Password Security**
- ✅ Bcrypt hashing with cost factor 10
- ✅ Passwords never exposed in API responses
- ✅ Minimum password length: 8 characters
- ✅ Old password verification for password change

### **Input Validation**
- ✅ Email format validation
- ✅ Age range validation (0-150)
- ✅ SQL injection prevention (parameterized queries)
- ✅ Request body validation with Gin validator

### **Database Security**
- ✅ Unique email constraint
- ✅ Indexed columns for performance
- ✅ Connection pooling with pgx
- ✅ Prepared statements

---

## 📈 Performance Optimization

### **Database Indexes**

```sql
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_name ON users(name);
CREATE INDEX idx_users_age ON users(age);
CREATE INDEX idx_users_created_at ON users(created_at);
```

### **Redis Caching**
- User data cached for fast retrieval
- Automatic cache invalidation on updates
- TTL-based expiration

### **Connection Pooling**
- PostgreSQL connection pool (pgxpool)
- Redis connection pool
- Configurable pool size

---

## 🚀 Deployment

### **Production Checklist**

- [ ] Change default passwords
- [ ] Set strong `DB_PASSWORD`
- [ ] Enable SSL for PostgreSQL (`sslmode=require`)
- [ ] Set up proper logging
- [ ] Configure monitoring alerts
- [ ] Set up backup strategy
- [ ] Use secrets management (not environment variables)
- [ ] Enable rate limiting
- [ ] Set up reverse proxy (Nginx)
- [ ] Configure CORS properly
- [ ] Enable HTTPS
- [ ] Set up CI/CD pipeline

### **Docker Production Build**

```dockerfile
# Multi-stage build for smaller image
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

### **Environment-Specific Configs**

```yaml
# docker-compose.prod.yml
version: '3.8'
services:
  app:
    restart: always
    environment:
      - DB_PASSWORD=${DB_PASSWORD}  # Use secrets
      - ENVIRONMENT=production
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
```

---

## 📚 API Best Practices

### **RESTful Principles**
- ✅ Use HTTP methods correctly (GET, POST, PUT, DELETE)
- ✅ Use plural nouns for resources (`/users`, not `/user`)
- ✅ Use proper HTTP status codes
- ✅ Version your API (`/api/v1`)

### **Response Codes**
| Code | Meaning | Usage |
|------|---------|-------|
| 200 | OK | Successful GET, PUT, DELETE |
| 201 | Created | Successful POST |
| 400 | Bad Request | Validation error |
| 404 | Not Found | Resource not found |
| 409 | Conflict | Duplicate resource (email) |
| 500 | Internal Server Error | Server error |

### **Pagination Best Practices**
- Default page size: 10
- Maximum page size: 100
- Always include total count
- Use `page` and `limit` parameters

---

## 🤝 Contributing

We welcome contributions! Please follow these guidelines:

### **How to Contribute**

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### **Code Style**

- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` for formatting
- Add comments for exported functions
- Write unit tests for new features

### **Commit Messages**

```
feat: add user profile endpoint
fix: resolve database connection issue
docs: update API documentation
test: add unit tests for user service
refactor: improve error handling
```

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 👥 Authors

- **Your Name** - *Initial work* - [hazelclein](https://github.com/hazelclein)

---

## 🙏 Acknowledgments

- [Gin Web Framework](https://gin-gonic.com/)
- [pgx PostgreSQL Driver](https://github.com/jackc/pgx)
- [Jaeger Tracing](https://www.jaegertracing.io/)
- [Docker](https://www.docker.com/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

---

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/hazelclein/user-crud-docker/issues)
- **Discussions**: [GitHub Discussions](https://github.com/hazelclein/user-crud-docker/discussions)
- **Email**: your.email@example.com

---

## 🗺 Roadmap

### **Planned Features**

- [ ] JWT Authentication
- [ ] Role-based Access Control (RBAC)
- [ ] User profile pictures
- [ ] Email verification
- [ ] Password reset functionality
- [ ] Audit logging
- [ ] GraphQL API
- [ ] WebSocket support
- [ ] Swagger/OpenAPI documentation
- [ ] Unit and integration tests
- [ ] CI/CD pipeline
- [ ] Kubernetes deployment manifests

---

## 📊 Database Schema

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    age INTEGER NOT NULL CHECK (age >= 0 AND age <= 150),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_name ON users(name);
CREATE INDEX idx_users_age ON users(age);
CREATE INDEX idx_users_created_at ON users(created_at);
```

---

## 🎯 Quick Start Guide

**For the impatient developer:**

```bash
# 1. Clone and enter
git clone https://github.com/hazelclein/user-crud-docker.git && cd user-crud-docker

# 2. Start everything
docker-compose up -d

# 3. Create a user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"password123","age":25}'

# 4. Get all users
curl http://localhost:8080/api/v1/users

# Done! 🎉
```

---

**Built with ❤️ using Go, PostgreSQL, Redis, and Docker**