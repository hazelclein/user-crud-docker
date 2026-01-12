# 🚀 User CRUD API - Complete Guide

![Go Version](https://img.shields.io/badge/Go-1.25.5-00ADD8?logo=go)
![Docker](https://img.shields.io/badge/Docker-28.5.2-2496ED?logo=docker)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?logo=postgresql)

Aplikasi CRUD User management dengan **Clean Architecture**, **Domain-Driven Design (DDD)**, **CQRS Pattern**, **Password Management**, **Search & Filter**, dan **Pagination**.

---

## 📋 Daftar Isi

1. [Fitur-Fitur](#-fitur-fitur)
2. [Prasyarat](#-prasyarat)
3. [Instalasi dari Awal](#-instalasi-dari-awal)
4. [Cara Menjalankan](#-cara-menjalankan)
5. [API Endpoints](#-api-endpoints)
6. [Testing Lengkap](#-testing-lengkap)
7. [Troubleshooting](#-troubleshooting)
8. [Arsitektur](#-arsitektur)

---

## ✨ Fitur-Fitur

### **Core Features:**
- ✅ **CRUD User** - Create, Read, Update, Delete users
- ✅ **Password Management** - Secure password hashing dengan bcrypt
- ✅ **Change Password** - User bisa ganti password sendiri
- ✅ **Search Users** - Cari user berdasarkan nama atau email
- ✅ **Filter by Age** - Filter user berdasarkan umur
- ✅ **Sorting** - Urutkan by id, name, email, age, created_at
- ✅ **Pagination** - Hasil dibagi per halaman (max 100/page)

### **Technical Features:**
- ✅ **Clean Architecture** - Separation of concerns yang jelas
- ✅ **Domain-Driven Design** - Domain sebagai center of business logic
- ✅ **CQRS Pattern** - Pemisahan Command (Write) dan Query (Read)
- ✅ **Graceful Shutdown** - Handle interrupt signal dengan baik
- ✅ **Health Check** - Monitoring endpoint untuk kesehatan aplikasi
- ✅ **Docker & Docker Compose** - Containerization lengkap
- ✅ **Database Indexes** - Optimized query performance
- ✅ **Input Validation** - Gin validator untuk semua input
- ✅ **Error Handling** - Consistent error response format

---
## 🚀 Cara Menjalankan

### **Method 1: Docker Compose (Recommended)**

Paling mudah dan tidak perlu setup database manual.

```bash
# 1. Build dan jalankan semua services
docker-compose up --build -d

# 2. Check logs
docker-compose logs -f app

# 3. Verify running
docker-compose ps

# 4. Test health check
curl http://localhost:8080/health
```

**Expected Output:**
```json
{
  "status": "healthy",
  "database": "connected",
  "timestamp": "2024-01-08T10:00:00Z"
}
```

### **Method 2: Local Development (Without Docker)**

Jika ingin run tanpa Docker (perlu PostgreSQL terinstall).

```bash
# 1. Start PostgreSQL dengan Docker saja
docker-compose up -d postgres

# 2. Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=userdb
export SERVER_PORT=8080

# 3. Run aplikasi
go run cmd/api/main.go
```

**Windows PowerShell:**
```powershell
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:DB_USER="postgres"
$env:DB_PASSWORD="postgres"
$env:DB_NAME="userdb"
$env:SERVER_PORT="8080"

go run cmd/api/main.go
```

### **Stop Services**

```bash
# Stop all services
docker-compose down

# Stop dan hapus volumes (reset database)
docker-compose down -v
```

---

## 📡 API Endpoints

### **Base URL:** `http://localhost:8080`

### **Health Check**
```
GET /health
```
Check application dan database status.

---

### **User Management**

#### **1. Create User**
```
POST /api/v1/users
Content-Type: application/json

Body:
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "age": 30
}

Response: 201 Created
{
  "status": "success",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30,
    "created_at": "2024-01-08T10:00:00Z",
    "updated_at": "2024-01-08T10:00:00Z"
  }
}
```

**Validation Rules:**
- `name`: required, string
- `email`: required, valid email format, unique
- `password`: required, minimum 8 characters
- `age`: required, 0-150

---

#### **2. Get User by ID**
```
GET /api/v1/users/:id

Example: GET /api/v1/users/1

Response: 200 OK
{
  "status": "success",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30,
    "created_at": "2024-01-08T10:00:00Z",
    "updated_at": "2024-01-08T10:00:00Z"
  }
}
```

---

#### **3. List Users (with Filters & Pagination)**
```
GET /api/v1/users?search=john&age_min=25&age_max=40&sort=name&order=asc&page=1&limit=10

Query Parameters:
- search: string (optional) - Search by name or email
- age_min: integer (optional) - Minimum age
- age_max: integer (optional) - Maximum age
- sort: string (optional) - Sort field (id, name, email, age, created_at)
- order: string (optional) - Sort order (asc, desc)
- page: integer (optional) - Page number (default: 1)
- limit: integer (optional) - Items per page (default: 10, max: 100)

Response: 200 OK
{
  "status": "success",
  "data": [...],
  "total": 50,
  "page": 1,
  "limit": 10,
  "total_pages": 5
}
```

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

---

#### **4. Search Users**
```
GET /api/v1/users/search?q=john&page=1&limit=10

Query Parameters:
- q: string (required) - Search keyword
- page: integer (optional)
- limit: integer (optional)

Response: 200 OK
{
  "status": "success",
  "data": [...],
  "total": 5,
  "page": 1,
  "limit": 10,
  "total_pages": 1
}
```

---

#### **5. Update User**
```
PUT /api/v1/users/:id
Content-Type: application/json

Body:
{
  "name": "John Updated",
  "email": "john.updated@example.com",
  "age": 31
}

Response: 200 OK
{
  "status": "success",
  "data": {
    "id": 1,
    "name": "John Updated",
    "email": "john.updated@example.com",
    "age": 31,
    "created_at": "2024-01-08T10:00:00Z",
    "updated_at": "2024-01-08T10:15:00Z"
  }
}
```

**Note:** Update user TIDAK mengubah password!

---

#### **6. Change Password**
```
PUT /api/v1/users/:id/change-password
Content-Type: application/json

Body:
{
  "old_password": "password123",
  "new_password": "newpassword456"
}

Response: 200 OK
{
  "status": "success",
  "message": "password changed successfully"
}
```

**Validation:**
- `old_password`: required, must match current password
- `new_password`: required, minimum 8 characters

---

#### **7. Delete User**
```
DELETE /api/v1/users/:id

Example: DELETE /api/v1/users/1

Response: 200 OK
{
  "status": "success",
  "message": "user deleted successfully"
}
```

---

### **Error Responses**

#### **400 Bad Request**
```json
{
  "status": "error",
  "message": "validation error message"
}
```

#### **404 Not Found**
```json
{
  "status": "error",
  "message": "user not found"
}
```

#### **409 Conflict**
```json
{
  "status": "error",
  "message": "user with this email already exists"
}
```

#### **500 Internal Server Error**
```json
{
  "status": "error",
  "message": "internal server error"
}
```

---

## 🧪 Testing Lengkap

### **Setup Test Data**

```bash
# Create 5 test users
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com", "password": "password123", "age": 25}'

curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Jane Smith", "email": "jane@example.com", "password": "password123", "age": 30}'

curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Bob Johnson", "email": "bob@test.com", "password": "password123", "age": 22}'

curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice Brown", "email": "alice@example.com", "password": "password123", "age": 28}'

curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Charlie Wilson", "email": "charlie@test.com", "password": "password123", "age": 35}'
```

### **Test Scenarios**

#### **1. Test Health Check**
```bash
curl http://localhost:8080/health
```

#### **2. Test Get All Users**
```bash
curl http://localhost:8080/api/v1/users
```

#### **3. Test Search**
```bash
# Search by name
curl "http://localhost:8080/api/v1/users?search=john"

# Search by email domain
curl "http://localhost:8080/api/v1/users?search=example.com"
```

#### **4. Test Filter by Age**
```bash
# Age 25-30
curl "http://localhost:8080/api/v1/users?age_min=25&age_max=30"

# Age >= 30
curl "http://localhost:8080/api/v1/users?age_min=30"
```

#### **5. Test Sorting**
```bash
# Sort by name (A-Z)
curl "http://localhost:8080/api/v1/users?sort=name&order=asc"

# Sort by age (oldest first)
curl "http://localhost:8080/api/v1/users?sort=age&order=desc"
```

#### **6. Test Pagination**
```bash
# Page 1, 2 items
curl "http://localhost:8080/api/v1/users?page=1&limit=2"

# Page 2, 2 items
curl "http://localhost:8080/api/v1/users?page=2&limit=2"
```

#### **7. Test Combined Filters**
```bash
curl "http://localhost:8080/api/v1/users?search=example&age_min=25&sort=age&order=asc&page=1&limit=10"
```

#### **8. Test Get by ID**
```bash
curl http://localhost:8080/api/v1/users/1
```

#### **9. Test Update**
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "John Updated", "email": "john@example.com", "age": 26}'
```

#### **10. Test Change Password**
```bash
curl -X PUT http://localhost:8080/api/v1/users/1/change-password \
  -H "Content-Type: application/json" \
  -d '{"old_password": "password123", "new_password": "newpassword456"}'
```

#### **11. Test Delete**
```bash
curl -X DELETE http://localhost:8080/api/v1/users/5
```

#### **12. Test Validation Errors**

**Empty password:**
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Test", "email": "test@test.com", "password": "", "age": 25}'
```

**Password too short:**
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Test", "email": "test@test.com", "password": "12345", "age": 25}'
```

**Invalid email:**
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Test", "email": "not-an-email", "password": "password123", "age": 25}'
```

**Duplicate email:**
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Test", "email": "john@example.com", "password": "password123", "age": 25}'
```
---

### **CQRS Implementation**

**Commands (Write Operations):**
- `CreateUserCommand` → Create new user
- `UpdateUserCommand` → Update user data
- `ChangePasswordCommand` → Change password
- `DeleteUserCommand` → Delete user

**Queries (Read Operations):**
- `GetUserQuery` → Get user by ID
- `ListUsersQuery` → List users with filters
- `SearchUsersQuery` → Search users

### **Domain-Driven Design**

**Entities:**
- `User` - Aggregate root with business rules

**Value Objects:**
- Email (with validation)
- Password (hashed)
- Age (with range validation)

**Repository Pattern:**
- Interface in domain layer
- Implementation in infrastructure layer

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

## 🔐 Security Features

1. **Password Hashing** - Bcrypt with cost 10
2. **Password Never Exposed** - JSON tag `-` prevents serialization
3. **Input Validation** - Gin validator untuk semua input
4. **SQL Injection Prevention** - Parameterized queries
5. **Email Uniqueness** - Database constraint
6. **Rate Limiting Ready** - Easy to add middleware

---

## 📝 Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | postgres | PostgreSQL hostname |
| `DB_PORT` | 5432 | PostgreSQL port |
| `DB_USER` | postgres | Database username |
| `DB_PASSWORD` | postgres | Database password |
| `DB_NAME` | userdb | Database name |
| `SERVER_PORT` | 8080 | HTTP server port |

---

Built with:
- **Go** - Programming language
- **Gin** - Web framework
- **PostgreSQL** - Database
- **pgx** - PostgreSQL driver
- **Docker** - Containerization
- **bcrypt** - Password hashing

---