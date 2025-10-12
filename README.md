# 🏪 Warehouse & Point of Sales - Golang Microservices API

[![Go Version](https://img.shields.io/badge/go-1.24.3-blue.svg)](https://golang.org)
[![Docker](https://img.shields.io/badge/docker-compose-blue.svg)](https://docker.com)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Hacktoberfest](https://img.shields.io/badge/hacktoberfest-2025-orange.svg)](https://hacktoberfest.com)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

A modern microservices system for warehouse management and point of sales built with Go, using microservices architecture with API Gateway pattern.

## 📋 Table of Contents

- [Description](#-description)
- [System Architecture](#-system-architecture)
- [Key Features](#-key-features)
- [Tech Stack](#-tech-stack)
- [Services Overview](#-services-overview)
- [Prerequisites](#-prerequisites)
- [Installation](#-installation)
- [Configuration](#-configuration)
- [Usage](#-usage)
- [API Documentation](#-api-documentation)
- [Database Schema](#-database-schema)
- [Development](#-development)
- [Deployment](#-deployment)
- [Monitoring](#-monitoring)
- [Contributing](#-contributing)
- [License](#-license)

## 📝 Description

This system is a complete implementation of warehouse management system and point of sales designed with microservices architecture. The system enables product management, inventory, transactions, merchants, notifications, and user management in a scalable and maintainable environment.

### System Advantages:
- **Microservices Architecture**: Each service is independent and can be scaled separately
- **API Gateway**: Centralized routing and authentication
- **Event-Driven**: Uses RabbitMQ for asynchronous communication
- **Caching**: Redis for performance optimization
- **Database per Service**: Each service has its own separate database
- **Containerized**: Full Docker support with docker-compose
- **Security**: JWT authentication and rate limiting

## 🏗️ Arsitektur Sistem

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Client Apps   │────│   API Gateway   │────│   Load Balancer │
│  (Web/Mobile)   │    │    (Port 8080)  │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                               │
                               ├─── Authentication & Rate Limiting
                               │
        ┌──────────────────────┼──────────────────────┐
        │                      │                      │
┌───────▼────┐    ┌────────────▼──┐    ┌─────────────▼───┐
│User Service│    │Product Service│    │Merchant Service │
│ (Port 8081)│    │  (Port 8082)  │    │  (Port 8084)    │
└────────────┘    └───────────────┘    └─────────────────┘
        │                  │                      │
        │                  │                      │
┌───────▼────┐    ┌────────▼──────┐    ┌─────────▼───────┐
│Transaction │    │Warehouse      │    │Notification     │
│Service     │    │Service        │    │Service          │
│(Port 8085) │    │(Port 8083)    │    │(Port 8086)      │
└────────────┘    └───────────────┘    └─────────────────┘
        │                  │                      │
        └──────────────────┼──────────────────────┘
                           │
    ┌──────────────────────┼──────────────────────┐
    │                      │                      │
┌───▼────┐        ┌────────▼─────┐        ┌──────▼──┐
│ Redis  │        │  PostgreSQL  │        │RabbitMQ │
│        │        │   (Multiple  │        │         │
│        │        │   Databases) │        │         │
└────────┘        └──────────────┘        └─────────┘
```

## ✨ Key Features

### 🔐 User Management
- User registration and authentication
- Role-based access control (RBAC)
- Profile management with photo upload
- JWT-based authentication

### 📦 Product Management
- CRUD operations for products
- Product categories
- Product image management
- Bulk operations

### 🏪 Merchant Management
- Merchant registration and management
- Merchant-specific products
- Inventory tracking per merchant
- Sales analytics

### 🏭 Warehouse Management
- Multiple warehouse support
- Inventory tracking
- Stock movement history
- Low stock alerts

### 💳 Transaction Processing
- Point of sales functionality
- Order management
- Payment integration
- Transaction history

### 📧 Notification System
- Email notifications
- Event-driven notifications
- Notification templates
- Delivery tracking

## 🛠️ Tech Stack

### Backend
- **Language**: Go 1.24.3
- **Framework**: Fiber v2 (Fast HTTP framework)
- **Architecture**: Microservices

### Database
- **Primary DB**: PostgreSQL (One per service)
- **Cache**: Redis
- **Message Queue**: RabbitMQ

### DevOps & Deployment
- **Containerization**: Docker & Docker Compose
- **Process Management**: Docker Healthchecks
- **Environment Management**: Environment variables

### Libraries & Tools
- **HTTP Router**: Fiber v2
- **ORM**: GORM
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Validation**: go-playground/validator/v10
- **Configuration**: Viper
- **Logging**: Zerolog
- **CLI**: Cobra
- **Storage**: Supabase Storage

## 🔧 Services Overview

| Service | Port | Database | Description |
|---------|------|----------|-------------|
| API Gateway | 8080 | Redis | Entry point, routing, auth, rate limiting |
| User Service | 8081 | PostgreSQL | User management, authentication |
| Product Service | 8082 | PostgreSQL | Product catalog, categories |
| Warehouse Service | 8083 | PostgreSQL | Inventory, warehouse management |
| Merchant Service | 8084 | PostgreSQL | Merchant operations, merchant products |
| Transaction Service | 8085 | PostgreSQL | Orders, payments, transactions |
| Notification Service | 8086 | PostgreSQL | Email, notifications |

### Infrastructure Services
- **PostgreSQL**: 6 separate databases (ports 5432-5437)
- **Redis**: Cache and session storage (port 6379)
- **RabbitMQ**: Message queue (ports 5672, 15672)

## 📋 Prerequisites

Make sure you have the following software installed:

- **Docker** (version 20.0+)
- **Docker Compose** (version 2.0+)
- **Go** 1.24.3+ (for development)
- **Git**

### Verify Prerequisites
```bash
docker --version
docker-compose --version
go version
git --version
```

## 🚀 Installation

### 1. Clone Repository
```bash
git clone https://github.com/your-username/micro-warehouse-api.git
cd micro-warehouse-api
```

### 2. Setup Environment Variables
Each service requires environment files. Copy the example files and customize them:

```bash
# API Gateway
cp api-gateway/env.example api-gateway/.env

# User Service
cp user-service/env.example user-service/.env

# Product Service
cp product-service/env.example product-service/.env

# Merchant Service
cp merchant-service/env.example merchant-service/.env

# Transaction Service
cp transaction-service/env.example transaction-service/.env

# Warehouse Service
cp warehouse-service/env.example warehouse-service/.env

# Notification Service
cp notification-service/env.example notification-service/.env
```

### 3. Build and Start Services
```bash
# Build all services
docker-compose build

# Start infrastructure services first
docker-compose up -d user_db product_db transaction_db merchant_db notification_db warehouse_db redis rabbitmq

# Wait until databases are ready, then start application services
docker-compose up -d
```

### 4. Verify Installation
```bash
# Check all containers running
docker-compose ps

# Check API Gateway health
curl http://localhost:8080/health
```

## ⚙️ Configuration

### API Gateway Configuration (.env)
```env
APP_ENV=development
APP_PORT=8080

# Service URLs
USER_SERVICE_URL=http://user-service:8081
PRODUCT_SERVICE_URL=http://product-service:8082
MERCHANT_SERVICE_URL=http://merchant-service:8084
WAREHOUSE_SERVICE_URL=http://warehouse-service:8083
TRANSACTION_SERVICE_URL=http://transaction-service:8085
NOTIFICATION_SERVICE_URL=http://notification-service:8086

# JWT Configuration
JWT_SECRET_KEY=your-secret-key
JWT_ISSUER=jwt-warehouse
JWT_DURATION=24h

# Redis Configuration
REDIS_HOST=warehouse_redis
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_POOL_SIZE=10
```

### Service-Specific Configuration
Each service has similar configuration with additional:
- Database connection settings
- Port configuration
- Service-specific features

### Database Configuration
```env
# Database per service
POSTGRES_USER=postgres
POSTGRES_PASSWORD=lokal
POSTGRES_DB=warehouse_[service]_db
POSTGRES_HOST=postgres-[service]
POSTGRES_PORT=5432
POSTGRES_SSLMODE=disable
```

## 🎯 Usage

### 1. API Gateway Health Check
```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "OK",
  "message": "API Gateway is running",
  "services": {
    "user-service": "http://user-service:8081",
    "product-service": "http://product-service:8082",
    "merchant-service": "http://merchant-service:8084",
    "warehouse-service": "http://warehouse-service:8083",
    "transaction-service": "http://transaction-service:8085",
    "notification-service": "http://notification-service:8086"
  }
}
```

### 2. User Registration
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123",
    "phone": "+6281234567890"
  }'
```

### 3. User Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### 4. Access Protected Endpoints
```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 📚 API Documentation

### Authentication Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | User login |
| POST | `/api/v1/auth/logout` | User logout |
| POST | `/api/v1/auth/refresh` | Refresh JWT token |

### User Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/users/profile` | Get user profile | ✅ |
| PUT | `/api/v1/users/profile` | Update profile | ✅ |
| POST | `/api/v1/users/upload-photo` | Upload profile photo | ✅ |

### Product Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/products` | List products | ✅ |
| POST | `/api/v1/products` | Create product | ✅ |
| GET | `/api/v1/products/:id` | Get product by ID | ✅ |
| PUT | `/api/v1/products/:id` | Update product | ✅ |
| DELETE | `/api/v1/products/:id` | Delete product | ✅ |

### Category Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/categories` | List categories | ✅ |
| POST | `/api/v1/categories` | Create category | ✅ |
| PUT | `/api/v1/categories/:id` | Update category | ✅ |
| DELETE | `/api/v1/categories/:id` | Delete category | ✅ |

### Merchant Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/merchants` | List merchants | ✅ |
| POST | `/api/v1/merchants` | Create merchant | ✅ |
| GET | `/api/v1/merchants/:id` | Get merchant | ✅ |
| PUT | `/api/v1/merchants/:id` | Update merchant | ✅ |

### Transaction Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/transactions` | List transactions | ✅ |
| POST | `/api/v1/transactions` | Create transaction | ✅ |
| GET | `/api/v1/transactions/:id` | Get transaction | ✅ |
| PUT | `/api/v1/transactions/:id/status` | Update status | ✅ |

### Warehouse Management
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/warehouses` | List warehouses | ✅ |
| POST | `/api/v1/warehouses` | Create warehouse | ✅ |
| GET | `/api/v1/warehouses/:id/inventory` | Get inventory | ✅ |
| POST | `/api/v1/warehouses/:id/stock` | Update stock | ✅ |

## 🗄️ Database Schema

### User Service Database
```sql
-- Users table
users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  photo VARCHAR(255),
  phone VARCHAR(20),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Roles table
roles (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50) UNIQUE NOT NULL,
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User-Role junction table
user_role (
  user_id INTEGER REFERENCES users(id),
  role_id INTEGER REFERENCES roles(id),
  PRIMARY KEY (user_id, role_id)
);
```

### Product Service Database
```sql
-- Categories table
categories (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Products table
products (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  price DECIMAL(10,2) NOT NULL,
  category_id INTEGER REFERENCES categories(id),
  sku VARCHAR(100) UNIQUE,
  image_url VARCHAR(500),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Merchant Service Database
```sql
-- Merchants table
merchants (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE,
  phone VARCHAR(20),
  address TEXT,
  user_id INTEGER,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Merchant Products table
merchant_products (
  id SERIAL PRIMARY KEY,
  merchant_id INTEGER REFERENCES merchants(id),
  product_id INTEGER,
  price DECIMAL(10,2),
  stock INTEGER DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🔧 Development

### Setup Development Environment
```bash
# Clone repository
git clone https://github.com/your-username/micro-warehouse-api.git
cd micro-warehouse-api

# Install Go dependencies for all services
for service in api-gateway user-service product-service merchant-service transaction-service warehouse-service notification-service; do
  cd $service
  go mod tidy
  cd ..
done
```

### Running Individual Services for Development
```bash
# Start infrastructure services
docker-compose up -d user_db product_db redis rabbitmq

# Run specific service
cd user-service
go run main.go start
```

### Hot Reload Development (Optional)
Install air for hot reload:
```bash
go install github.com/cosmtrek/air@latest

# Run with hot reload
cd user-service
air
```

### Testing
```bash
# Run tests for all services
for service in user-service product-service merchant-service transaction-service warehouse-service notification-service; do
  echo "Testing $service..."
  cd $service
  go test ./...
  cd ..
done
```

### Code Structure
```
service-name/
├── main.go                 # Entry point
├── cmd/                   # CLI commands (Cobra)
│   ├── root.go
│   └── start.go
├── app/                   # Application setup
│   ├── app.go
│   ├── container.go       # Dependency injection
│   └── routes.go          # Route definitions
├── configs/               # Configuration
│   └── config.go
├── controller/            # HTTP handlers
│   ├── *_controller.go
│   ├── request/           # Request DTOs
│   └── response/          # Response DTOs
├── database/              # Database connection
│   └── postgres_database.go
├── middleware/            # Custom middleware
├── model/                 # Database models
├── repository/            # Data access layer
├── usecase/              # Business logic
└── pkg/                  # Shared utilities
    ├── jwt/
    ├── pagination/
    ├── validator/
    └── ...
```

## 🚀 Deployment

### Production Deployment with Docker

1. **Setup Production Environment**
```bash
# Copy and edit environment variables
cp api-gateway/env.example api-gateway/.env.prod

# Edit for production values
# - Change database passwords
# - Set secure JWT secrets
# - Configure external services
```

2. **Build Production Images**
```bash
# Build all services
docker-compose -f docker-compose.prod.yml build

# Push to registry (optional)
docker-compose -f docker-compose.prod.yml push
```

3. **Deploy**
```bash
# Deploy to production
docker-compose -f docker-compose.prod.yml up -d

# Check status
docker-compose -f docker-compose.prod.yml ps
```

### Kubernetes Deployment (Optional)
```yaml
# Example kubernetes deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-gateway
spec:
  replicas: 3
  selector:
    matchLabels:
      app: api-gateway
  template:
    metadata:
      labels:
        app: api-gateway
    spec:
      containers:
      - name: api-gateway
        image: warehouse/api-gateway:latest
        ports:
        - containerPort: 8080
        env:
        - name: REDIS_HOST
          value: "redis-service"
```

### Environment Variables for Production
```env
# Security
JWT_SECRET_KEY=super-secure-secret-key-change-in-production
BCRYPT_COST=12

# Database
POSTGRES_PASSWORD=super-secure-db-password

# Redis
REDIS_PASSWORD=secure-redis-password

# External Services
SUPABASE_URL=your-supabase-url
SUPABASE_KEY=your-supabase-key

# Email
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

## 📊 Monitoring

### Health Checks
```bash
# API Gateway health
curl http://localhost:8080/health

# Individual service health
curl http://localhost:8081/health  # User Service
curl http://localhost:8082/health  # Product Service
# ... etc
```

### Docker Health Monitoring
```bash
# Check container health status
docker-compose ps

# View logs
docker-compose logs -f api-gateway
docker-compose logs -f user-service

# Monitor resource usage
docker stats
```

### RabbitMQ Management
Access RabbitMQ management interface:
- URL: http://localhost:15672
- Username: guest
- Password: guest

### Database Monitoring
```bash
# Connect to specific database
docker exec -it postgres-user psql -U postgres -d warehouse_user_db

# Check database size
SELECT pg_size_pretty(pg_database_size('warehouse_user_db'));

# Monitor active connections
SELECT count(*) FROM pg_stat_activity;
```

### Redis Monitoring
```bash
# Connect to Redis
docker exec -it warehouse_redis redis-cli

# Monitor Redis
redis-cli INFO memory
redis-cli MONITOR
```

## 🤝 Contributing

We warmly welcome contributions from the community! This repository participates in **Hacktoberfest 2025** 🎃

[![Hacktoberfest](https://img.shields.io/badge/hacktoberfest-2025-orange.svg)](https://hacktoberfest.com)

### 🌟 For Hacktoberfest Contributors

This repository accepts contributions for Hacktoberfest 2025! Look for issues with labels:
- `hacktoberfest` - for Hacktoberfest participants
- `good first issue` - suitable for beginners
- `help wanted` - needs community help

### 📖 Complete Contribution Guide

For detailed guidelines on how to contribute, please read **[CONTRIBUTING.md](CONTRIBUTING.md)** which covers:

- Setup development environment
- Coding standards & guidelines
- Testing requirements
- Pull request process
- Area kontribusi yang dibutuhkan
- Good first issues for beginners

### 🚀 Quick Start for Contributors

1. **Fork & Clone repository**
2. **Setup development environment** (see [CONTRIBUTING.md](CONTRIBUTING.md))
3. **Select issue** what needs to be done
4. **Create a feature branch**
5. **Develop & test**  your changes
6. **Submit Pull Request**

### 💡 Contribution Area

- 🐛 **Bug Fixes** - Fix existing bugs
- 📚 **Documentation** - Improve docs & tutorials
- ✨ **New Features** - Add new features
- 🧪 **Testing** - Add unit & integration tests
- 🔒 **Security** - Security improvements
- 🚀 **Performance** - Performance optimization



## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Fiber](https://github.com/gofiber/fiber) - Fast HTTP framework
- [GORM](https://github.com/go-gorm/gorm) - ORM library
- [Viper](https://github.com/spf13/viper) - Configuration management
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [PostgreSQL](https://www.postgresql.org) - Database
- [Redis](https://redis.io) - Caching and session storage
- [RabbitMQ](https://www.rabbitmq.com) - Message queue

## 📞 Support

If you have questions or issues:

1. Check existing [Issues](https://github.com/your-username/micro-warehouse-api/issues)
2. Create new issue with complete details

## 🔮 Roadmap

### Version 2.0
- [ ] GraphQL API support
- [ ] WebSocket real-time notifications
- [ ] Advanced analytics dashboard
- [ ] Multi-tenant support
- [ ] Advanced inventory forecasting

### Version 2.1
- [ ] Mobile app integration
- [ ] Payment gateway integration
- [ ] Advanced reporting system
- [ ] Audit logging
- [ ] Performance monitoring

---
