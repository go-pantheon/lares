# Lares

Lares is a high-performance account system service framework based on microservice architecture, developed in Go. This framework aims to provide scalable user authentication and account management infrastructure, supporting various login methods and authentication processes. Lares is a core component of the go-pantheon ecosystem, responsible for implementing user account management and authentication logic.

## go-pantheon Ecosystem

**go-pantheon** is an out-of-the-box game server framework providing high-performance, highly available game server cluster solutions based on microservices architecture. Lares, as the account management implementation component, works alongside other core services to form a complete game service ecosystem:

- **Roma**: Game core logic services
- **Janus**: Gateway service for client connection handling and request forwarding
- **Lares**: Account service for user authentication and account management
- **Senate**: Backend management service providing operational interfaces

### Core Features

- 🚀 Microservice account system architecture built with [go-kratos](https://github.com/go-kratos/kratos)
- 🔒 Multi-platform login support (username/password, Apple, Google, Facebook, etc.)
- 🛡️ Enterprise-grade secure communication protocol and authentication mechanism
- 📊 Real-time monitoring and distributed tracing
- 🔄 Gray release and hybrid deployment support
- 🔍 Developer-friendly debugging environment
- 🔑 High-performance token generation and validation mechanism

### Service Layer Features

- gRPC for inter-service communication
- Dual HTTP/gRPC API support
- Session management and token generation
- Secure encryption and replay attack prevention
- High concurrency handling capability

## System Architecture

The relationship between Lares and other go-pantheon components is illustrated below:

```
   (1)HTTPS Token Request
┌────────┐                  ┌────────┐
│        │------------------>        │
│ Client │                  │ Lares  │
│        │<-----------------|        │
└────────┘   Return Token   └────────┘
    │
    │  (2)TCP Connection
    │  & Token Handshake
    │
    ▼                       ┌────────────────┐
┌────────┐                  │                │
│        │  (4)Game Protocol│     Roma       │
│ Janus  │<---------------->│    (Hidden)    │
│        │    gRPC Tunnel   │                │
└────────┘                  └────────────────┘
    ▲                              │
    │                              │
    │  (3)Game Protocol            │
    │     TCP                      ▼
    │                        ┌──────────────┐
└────────┘                   │ Senate Admin │
  Client                     └──────────────┘
```

### Authentication Flow Details

```
┌────────┐        ┌────────┐        ┌────────┐        ┌────────┐
│ Client │        │ Lares  │        │ Janus  │        │  Roma  │
└───┬────┘        └───┬────┘        └───┬────┘        └───┬────┘
    │                 │                 │                 │
    │ 1.HTTPS Token Request             │                 │
    │---------------->│                 │                 │
    │                 │                 │                 │
    │ Return Token    │                 │                 │
    │<----------------│                 │                 │
    │                 │                 │                 │
    │ 2.Establish TCP Connection        │                 │
    │---------------------------------->│                 │
    │                 │                 │                 │
    │ Send Token in Handshake           │                 │
    │---------------------------------->│                 │
    │                 │                 │                 │
    │                 │                 │ (Internal Token │
    │                 │                 │  Verification)  │
    │                 │                 │                 │
    │ TCP Handshake Success Response    │                 │
    │<----------------------------------│                 │
    │                 │                 │                 │
    │ 3.Send Login Protocol (TCP)       │                 │
    │---------------------------------->│                 │
    │                 │                 │                 │
    │                 │                 │ 4.Select Roma Service │
    │                 │                 │ Based on Token Info   │
    │                 │                 │                 │
    │                 │                 │ Establish gRPC Tunnel │
    │                 │                 │---------------->│
    │                 │                 │                 │
    │                 │                 │ Tunnel Established      
    │                 │                 │<----------------│
    │                 │                 │                 │
    │ Game Protocol Messages (TCP)      │ Forward as gRPC │
    │---------------------------------->│---------------->│
    │                 │                 │                 │
    │                 │                 │ Game Logic Processing │
    │                 │                 │<----------------│
    │ Response Messages (TCP)           │                 │
    │<----------------------------------│                 │
    │                 │                 │                 │
```

Lares internally adopts a microservice architecture, with services communicating via gRPC:

```
                    ┌─────────────────────────────┐
                    │           Lares             │
                    │                             │
┌─────────────┐     │  ┌───────────┐ ┌─────────┐ │
│  Config     │◀───▶│  │Account    │ │Auth     │ │
│  (etcd)     │     │  │Interface  │ │Service  │ │
└─────────────┘     │  └───────────┘ └─────────┘ │
                    │        ▲            ▲      │
┌─────────────┐     │        │            │      │
│  Monitoring │◀───▶│        └─────┬──────┘      │
│(Prometheus) │     │              │             │
└─────────────┘     │         ┌────▼────┐        │
                    │         │ Redis   │        │
┌─────────────┐     │         │ Session │        │
│  Tracing    │◀───▶│         │ Store   │        │
│     (OT)    │     │                             │
└─────────────┘     └─────────────────────────────┘
```

## Project Overview

The Lares framework is built on the Go-Kratos microservice framework, supporting both gRPC and HTTP protocols, and integrates core components such as etcd registry, Redis cache, and PostgreSQL database. The framework design follows Domain-Driven Design (DDD) principles, achieving high cohesion and low coupling of business logic through clear service boundaries and domain models.

## Technology Stack

Lares utilizes the following core technologies:

| Technology/Component | Purpose | Version |
|---------|------|------|
| Go | Primary development language | 1.23+ |
| go-kratos | Microservice framework | v2.8.4 |
| gRPC | Inter-service communication | v1.71.1 |
| Protobuf | Data serialization | v1.36.6 |
| etcd | Service discovery & registry | v3.5.21 |
| Redis | Caching and session storage | v9.7.3 |
| PostgreSQL | Data storage | v1.5.11 |
| OpenTelemetry | Distributed tracing | v1.35.0 |
| Prometheus | Monitoring system | v1.22.0 |
| Google Wire | Dependency injection | v0.6.0 |
| JWT | Token generation and validation | v4.5.2 |
| Buf | API management | Latest |

## Key Features

- **Microservice Architecture**: Built on Go-Kratos with service registry, discovery, and load balancing
- **Multi-Protocol Support**: Simultaneous support for gRPC and HTTP interfaces
- **Configuration Center**: etcd-based configuration center with dynamic updates
- **Multi-Platform Login**: Support for username/password, Apple, Google, Facebook, and other login methods
- **Distributed Tracing**: OpenTelemetry integration for distributed tracing
- **Service Monitoring**: Prometheus metrics collection
- **Dependency Injection**: Google Wire for dependency injection
- **Code Generation**: Simplified development through Protobuf and code generation tools
- **High-Performance Session Management**: Redis-based high-performance session storage and management
- **Secure Token Mechanism**: Support for JWT token generation and validation

## Core Components

### Application Services (app/)

- **account**: Business functionality related to account management
- **notice**: Notification service business functionality

### API Definitions (api/)

- **server**: Server-side internal API definitions
- **interface**: Client-facing interface definitions

### Common Libraries (pkg/)

- **security**: Security-related functionality
- **util**: Utility functions

## Requirements

- Go 1.23+
- Protobuf
- etcd
- Redis
- PostgreSQL

## Quick Start

### Initialize Environment

```bash
make init
```

### Generate API Code

```bash
make proto
make api
```

### Build Services

```bash
make build
```

### Start Services

```bash
# Start all services
make run

# Start a specific service
make run app=account
```

## Integration with go-pantheon Components

Integration of Lares with other go-pantheon components typically follows these steps:

### Integration with Janus Gateway

1. Configure Janus service registry information to ensure discovery by Lares
2. Set up Janus service routing rules in Lares
3. Configure AES encryption keys for Token in both Lares and Janus
4. Client first obtains authentication token from Lares, then establishes connection with Janus

```yaml
# Janus configuration example
services:
  - name: account
    discovery:
      type: etcd
      address: ["127.0.0.1:2379"]
    endpoints:
      - protocol: grpc
        port: 9000
      - protocol: http
        port: 8000
```

### Integration with Roma Game Services

1. Janus service verifies client handshake requests using tokens generated by Lares
2. Other services can query and update account information through internal interfaces provided by Lares

```proto
# api/server/account/interface/account/v1/account.proto
service AccountInterface {
  // Get TCP handshake token
  rpc Token (TokenRequest) returns (TokenResponse) {
    option (google.api.http) = {
      post: "/accounts/v1/token"
      body: "*"
    };
  }
}
```

### Integration with Senate Backend Management

1. Ensure Lares services expose necessary management APIs
2. Call `api/server/account/admin` interfaces in the Senate service

## Project Structure

```
.
├── api/                # API definitions
│   └── server/         # Server-side API
│       ├── account/    # Account-related API
│       └── notice/     # Notice-related API
├── app/                # Application services
│   ├── account/        # Account service
│   └── notice/         # Notice service
├── deps/               # Local dependencies
├── gen/                # Generated code
├── pkg/                # Common libraries
│   ├── security/       # Security-related
│   └── util/           # Utility functions
└── third_party/        # Third-party dependencies
```

## Port Conventions

### Account Service

- HTTP Ports:
  - Internal: 8101
  - External: 8001
- gRPC Ports:
  - Internal: 9101
  - External: 9001

## Development Guide

### Development Workflow

1. Define API interfaces using Protobuf
2. Generate interface code with `make proto` and `make api`
3. Implement service logic based on business requirements
4. Use Wire for dependency injection
5. Write unit tests
6. Build and deploy services

### Adding New Services

Steps to create a new service:

1. Create Proto definitions for the new service in the `api/server/` directory
2. Generate API code: `make proto`
3. Create a new service directory in `app/`
4. Copy and modify the framework code from existing services
5. Generate dependency injection code using Wire: `make wire`
6. Implement service logic

### Multi-Platform Login Support

Lares supports various login methods:

- Username/password login
- Apple login
- Google login
- Facebook login

To add a new login method:

1. Define relevant interfaces in `api/server/account/interface/account/v1/account.proto`
2. Implement login logic in `app/account/internal/service`
3. Update security configuration and token generation logic

## Troubleshooting

### 1. Service Registration Failure

**Issue**: Service cannot register with etcd

**Solution**:
- Check if etcd is running properly
- Verify that the etcd address in the configuration file is correct
- Check network connectivity

### 2. Code Generation Errors

**Issue**: Generated code has errors after running `make proto`

**Solution**:
- Ensure all necessary protoc plugins are installed
- Check proto file syntax
- Run `make init` to reinstall all dependency tools

### 3. Configuration Errors During Service Startup

**Issue**: Service fails to start with configuration errors

**Solution**:
- Check configuration files in the `configs/` directory
- Ensure all necessary environment variables are set
- Reference template files in the `configs.tmpl/` directory

## Contributing

1. Fork this repository
2. Create a feature branch
3. Submit changes
4. Ensure all tests pass
5. Submit a Pull Request

## License

This project is licensed under the terms specified in the LICENSE file.
