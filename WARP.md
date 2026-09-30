# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Project Overview

Trusture is a blockchain-based decentralized framework for secure, transparent, and auditable NGO transactions. It provides tamper-evident, auditable donation management for NGOs, donors, regulators, and auditors using blockchain technology, smart contracts, and cryptographic anchoring on the Polygon network.

## Development Commands

### Backend (Go)

```bash
# Run the demo application (existing functionality)
go run cmd/main.go

# Run the API server (production-ready HTTP API)
go run cmd/api/main.go

# Build the demo binary
go build -o trusture cmd/main.go

# Build the API server binary
go build -o trusture-api cmd/api/main.go

# Run with verbose output
go run -v cmd/api/main.go

# Format code
go fmt ./...

# Vet code for potential issues
go vet ./...

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Install dependencies
go mod tidy
go mod download

# Generate Swagger documentation
swag init -g cmd/api/main.go -o docs/
```

### Frontend (React + Vite)

```bash
# Navigate to frontend directory first
cd Frontend

# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Lint code
npm run lint
```

### Full Stack Development

```bash
# Start backend (from root directory)
go run cmd/main.go

# In another terminal, start frontend
cd Frontend && npm run dev
```

## Architecture Overview

### High-Level System Components

The Trusture framework consists of several key architectural layers:

1. **HTTP API Server**: Production-ready REST API with JWT authentication
2. **Blockchain Layer**: Custom blockchain implementation for donations and expenditures
3. **Smart Contracts**: Polygon integration with cryptographic anchoring
4. **Entity Management**: NGOs, Donors, and Auditors with KYC verification
5. **Transparency Engine**: Real-time scoring and rating system
6. **Database Layer**: PostgreSQL with GORM ORM for data persistence
7. **Frontend Dashboards**: React-based interfaces for different user roles

### Backend Architecture (Go)

The Go backend is organized into the following package structure:

```
cmd/
├── api/main.go       # HTTP API server entry point
└── main.go           # Demo application entry point
pkg/
├── auth/             # JWT authentication and middleware
├── blockchain/       # Custom blockchain implementation
├── config/           # Environment configuration management
├── crypto/           # Cryptographic primitives (ZK proofs, multisig)
├── database/         # Database models and repositories (GORM)
├── entities/         # Core business entities (NGO, Donor, Auditor)
├── middleware/       # HTTP middleware (logging, CORS, etc.)
├── platform/         # Main platform orchestrator
├── polygon/          # Polygon blockchain integration
├── server/           # HTTP server and API handlers
└── transactions/     # Transaction processing (donations, expenditures)
docs/                 # API documentation (Swagger)
```

#### Key Components:

- **HTTP API Server** (`cmd/api/`, `pkg/server/`): Production REST API with authentication and middleware
- **Authentication System** (`pkg/auth/`): JWT-based authentication with bcrypt password hashing
- **Database Layer** (`pkg/database/`): PostgreSQL integration with GORM ORM and repositories
- **Configuration Management** (`pkg/config/`): Environment-based configuration with .env support
- **Middleware Stack** (`pkg/middleware/`): HTTP middleware for logging, CORS, security, and rate limiting
- **Platform Layer** (`pkg/platform/`): Central orchestrator that manages all system interactions
- **Blockchain Engine** (`pkg/blockchain/`): Custom permissioned blockchain for transaction recording
- **Entity Management** (`pkg/entities/`): Business logic for NGOs, donors, and auditors
- **Cryptographic Layer** (`pkg/crypto/`): Zero-knowledge proofs and multi-signature wallets
- **Polygon Integration** (`pkg/polygon/`): Anchoring and verification on public blockchain
- **Transaction Processing** (`pkg/transactions/`): Donation and expenditure workflow management

### Frontend Architecture (React)

The React frontend provides role-based dashboards:

- **Donor Dashboard**: NGO discovery, donation tracking, tax benefits
- **NGO Dashboard**: Campaign management, fund tracking, transparency metrics
- **Wallet Integration**: Web3 wallet connectivity for blockchain transactions

### Core Business Logic

#### Donation Flow:
1. Donor and NGO undergo KYC verification
2. Donation initiated with platform fee calculation
3. Zero-knowledge proof generated for privacy
4. Transaction recorded on permissioned blockchain
5. Block hash anchored to Polygon for immutability
6. E-bill generated with cryptographic signature

#### Expenditure Flow:
1. NGO submits expenditure with invoice details
2. Auditor validates compliance and documentation
3. Multi-signature approval process
4. Block recorded with GSTIN verification
5. Anchored to Polygon with expenditure proof

#### Rating System:
- Dynamic scoring based on donation utilization
- Transparency metrics from blockchain verification
- Documentation quality assessment
- KYC and certificate bonuses

## Key Technologies

- **Backend**: Go 1.24.4
- **Frontend**: React 19.1.1, Vite, Tailwind CSS
- **Blockchain**: Custom Go implementation + Polygon integration
- **Cryptography**: Zero-knowledge proofs, Multi-signature wallets
- **Storage**: Off-chain IPFS/S3 with cryptographic hashing

## Development Workflow

### Setting Up Development Environment

1. **Prerequisites**:
   - Go 1.24.4+
   - Node.js 16+
   - Git

2. **Clone and Setup**:
   ```bash
   git clone <repository>
   cd trusture
   go mod tidy
   cd Frontend && npm install
   ```

3. **Start Development**:
   ```bash
   # Terminal 1: Backend
   go run cmd/main.go
   
   # Terminal 2: Frontend
   cd Frontend && npm run dev
   ```

### Adding New Features

When adding new functionality:

1. **Entities**: Add new business logic to `pkg/entities/`
2. **Transactions**: Extend transaction types in `pkg/transactions/`
3. **Blockchain**: Modify blockchain logic in `pkg/blockchain/`
4. **Platform**: Update orchestration in `pkg/platform/`
5. **Frontend**: Add UI components in `Frontend/src/components/`

### Testing Strategy

The current codebase includes:
- Comprehensive demo scenarios in `cmd/main.go`
- Mock data structures for testing
- Simulated blockchain operations
- Integration testing with Polygon testnet

### Polygon Integration

The system integrates with Polygon Mumbai Testnet for:
- Block hash anchoring
- Transaction verification
- Gas cost estimation
- Network statistics

Configuration in `pkg/polygon/integration.go` includes:
- Provider URL configuration
- Gas price and limit settings
- Contract deployment simulation
- Anchor verification system

### Security Considerations

- KYC verification for all participants
- Multi-signature wallet requirements
- Zero-knowledge proofs for privacy
- Cryptographic anchoring on public blockchain
- GSTIN and invoice validation
- Immutable transaction records

## Common Development Tasks

### Adding a New Entity Type
1. Create entity struct in `pkg/entities/`
2. Implement verification and validation methods
3. Add to platform orchestrator in `pkg/platform/`
4. Update frontend components if UI is needed

### Extending Transaction Types
1. Add transaction struct in `pkg/transactions/`
2. Implement validation and processing logic
3. Update blockchain recording in `pkg/blockchain/`
4. Add Polygon anchoring if required

### Modifying Rating Algorithm
1. Update `CalculateRating()` in `pkg/entities/ngo.go`
2. Adjust transparency scoring logic
3. Update dashboard display in frontend components

### Adding New Blockchain Features
1. Extend blockchain functionality in `pkg/blockchain/`
2. Update block validation logic
3. Modify consensus mechanisms if needed
4. Update Polygon integration for new data types

The codebase is designed for extensibility and follows clean architecture principles with clear separation of concerns between blockchain, business logic, and presentation layers.