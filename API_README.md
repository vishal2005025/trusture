# Trusture API Server

A complete REST API server for the Trusture blockchain-based NGO donation auditing framework.

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12+ (or use in-memory mode for development)
- Git

### 1. Clone and Setup

```bash
git clone <repository-url>
cd trusture
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Environment Configuration

Copy the example environment file:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```bash
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
GIN_MODE=debug

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=trusture
DB_PASSWORD=trusture123
DB_NAME=trusture_db
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRY_HOURS=24

# Blockchain Configuration
POLYGON_RPC=https://polygon-mumbai.g.alchemy.com/v2/demo
POLYGON_PRIVATE_KEY=1111111111111111111111111111111111111111111111111111111111111111
POLYGON_GAS_LIMIT=300000
POLYGON_GAS_PRICE_GWEI=30
```

### 4. Database Setup

#### Option A: PostgreSQL (Recommended)

1. Install PostgreSQL
2. Create database and user:

```sql
CREATE DATABASE trusture_db;
CREATE USER trusture WITH PASSWORD 'trusture123';
GRANT ALL PRIVILEGES ON DATABASE trusture_db TO trusture;
```

#### Option B: In-Memory Development (No DB required)

Set `DB_HOST=memory` in your `.env` file.

### 5. Run the API Server

```bash
# Run the demo (existing functionality)
go run cmd/main.go

# Run the API server
go run cmd/api/main.go
```

The API server will start on `http://localhost:8080`

### 6. Access API Documentation

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **API Documentation**: http://localhost:8080/docs
- **Health Check**: http://localhost:8080/health
- **API Status**: http://localhost:8080/api/v1/status

## 📡 API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh JWT token
- `POST /api/v1/auth/logout` - User logout

### Public Endpoints (No authentication required)
- `GET /api/v1/stats` - Platform statistics
- `GET /api/v1/ngos` - List verified NGOs
- `GET /api/v1/ngos/{id}` - Get NGO profile
- `GET /api/v1/ngos/{id}/rating` - Get NGO rating
- `GET /api/v1/status` - System status
- `GET /api/v1/verify/{hash}` - Verify blockchain data

### NGO Endpoints (Requires NGO authentication)
- `GET /api/v1/ngos/profile` - Get NGO profile
- `GET /api/v1/ngos/dashboard` - Get NGO dashboard
- `POST /api/v1/ngos/expenditures` - Create expenditure
- `GET /api/v1/ngos/expenditures` - List expenditures
- `GET /api/v1/ngos/donations` - List received donations

### Donor Endpoints (Requires Donor authentication)
- `GET /api/v1/donors/profile` - Get donor profile
- `GET /api/v1/donors/dashboard` - Get donor dashboard
- `POST /api/v1/donors/donations` - Create donation
- `GET /api/v1/donors/donations` - List donations
- `GET /api/v1/donors/tax-benefits` - Get tax benefits

### Auditor Endpoints (Requires Auditor authentication)
- `GET /api/v1/auditors/profile` - Get auditor profile
- `GET /api/v1/auditors/dashboard` - Get auditor dashboard
- `GET /api/v1/auditors/pending-expenditures` - Get pending audits
- `POST /api/v1/auditors/audit/{expenditure_id}` - Audit expenditure

## 🔐 Authentication

The API uses JWT (JSON Web Tokens) for authentication. Include the token in the Authorization header:

```bash
curl -H "Authorization: Bearer <your-jwt-token>" \
     http://localhost:8080/api/v1/ngos/profile
```

## 📝 Example API Usage

### 1. Register a new NGO

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "ngo@example.com",
    "password": "password123",
    "user_type": "ngo",
    "name": "Example NGO",
    "registration_number": "REG123456",
    "category": "Education"
  }'
```

### 2. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "ngo@example.com",
    "password": "password123"
  }'
```

### 3. Get NGO Profile (with JWT token)

```bash
curl -H "Authorization: Bearer <jwt-token>" \
     http://localhost:8080/api/v1/ngos/profile
```

### 4. Get Platform Statistics

```bash
curl http://localhost:8080/api/v1/stats
```

## 🏗️ Architecture

### Backend Components

- **HTTP Server**: Gin web framework with middleware
- **Authentication**: JWT-based with bcrypt password hashing
- **Database**: GORM with PostgreSQL (or in-memory for development)
- **Blockchain**: Custom blockchain + Polygon integration
- **Cryptography**: Zero-knowledge proofs and multi-signature support
- **Logging**: Structured logging with Logrus
- **Documentation**: Swagger/OpenAPI 3.0

### Middleware Stack

1. **Error Recovery**: Panic recovery and error handling
2. **Request Logging**: Comprehensive request/response logging
3. **CORS**: Cross-origin resource sharing
4. **Security Headers**: XSS protection, content type sniffing prevention
5. **Rate Limiting**: IP-based rate limiting (100 req/min default)
6. **Authentication**: JWT validation for protected routes
7. **Content-Type Validation**: JSON content-type enforcement

## 🔧 Development

### Project Structure

```
├── cmd/
│   ├── api/main.go          # API server entry point
│   └── main.go              # Demo application
├── pkg/
│   ├── auth/                # JWT authentication
│   ├── config/              # Configuration management
│   ├── database/            # Database models and repositories
│   ├── middleware/          # HTTP middleware
│   └── server/              # HTTP server and handlers
├── docs/                    # API documentation
├── Frontend/                # React frontend (separate)
├── .env.example             # Environment variables template
└── README.md
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/auth
```

### Code Generation

Generate Swagger documentation:

```bash
# Install swag CLI
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
swag init -g cmd/api/main.go -o docs/
```

## 🐳 Docker Support

### Build Docker Image

```bash
docker build -t trusture-api .
```

### Run with Docker Compose

```bash
docker-compose up -d
```

## 📊 Monitoring

### Health Checks

- **Basic Health**: `GET /health`
- **System Status**: `GET /api/v1/status`
- **Platform Stats**: `GET /api/v1/stats`

### Logging

Logs are structured JSON format by default. Set `LOG_FORMAT=text` for human-readable logs.

Log levels: `debug`, `info`, `warn`, `error`

## 🔒 Security Features

- **JWT Authentication** with configurable expiration
- **Password Hashing** using bcrypt with high cost
- **Rate Limiting** to prevent abuse
- **CORS Protection** with configurable origins
- **Security Headers** (HSTS, XSS Protection, etc.)
- **Input Validation** with Gin binding validation
- **SQL Injection Protection** via GORM ORM

## 🚀 Production Deployment

### Environment Variables for Production

```bash
# Production settings
ENVIRONMENT=production
GIN_MODE=release
LOG_LEVEL=info
LOG_FORMAT=json

# Strong JWT secret
JWT_SECRET=<generate-strong-random-secret>

# Production database
DB_HOST=<production-db-host>
DB_PASSWORD=<secure-password>
DB_SSLMODE=require

# Production Polygon settings
POLYGON_RPC=<mainnet-rpc-url>
POLYGON_PRIVATE_KEY=<secure-private-key>
```

### Recommended Production Setup

1. **Reverse Proxy**: Use nginx or similar
2. **SSL/TLS**: Terminate SSL at load balancer/proxy
3. **Database**: Use managed PostgreSQL service
4. **Monitoring**: Set up log aggregation and monitoring
5. **Secrets**: Use environment variables or secret management service

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

MIT License - see LICENSE file for details.

## 📞 Support

- Documentation: `/swagger/index.html`
- Issues: GitHub Issues
- Email: support@trusture.io