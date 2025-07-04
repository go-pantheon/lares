<div align="center">
  <h1>🔐 Lares Game Account Authentication Service</h1>
  <p><em>High-performance game core business service framework for the go-pantheon ecosystem</em></p>
</div>

<p align="center">
<a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go" alt="Go"></a>
<a href="https://github.com/go-kratos/kratos"><img src="https://img.shields.io/badge/Kratos-v2.8.4-blue" alt="Kratos"></a>
<a href="https://github.com/go-pantheon/lares/actions/workflows/test.yml"><img src="https://github.com/go-pantheon/lares/workflows/Test/badge.svg" alt="Test Status"></a>
<a href="https://github.com/go-pantheon/lares/releases"><img src="https://img.shields.io/github/v/release/go-pantheon/lares" alt="Latest Release"></a>
<a href="https://pkg.go.dev/github.com/go-pantheon/lares"><img src="https://pkg.go.dev/badge/github.com/go-pantheon/lares" alt="GoDoc"></a>
<a href="https://goreportcard.com/report/github.com/go-pantheon/lares"><img src="https://goreportcard.com/badge/github.com/go-pantheon/lares" alt="Go Report Card"></a>
<a href="https://github.com/go-pantheon/lares/blob/main/LICENSE"><img src="https://img.shields.io/github/license/go-pantheon/lares" alt="License"></a>
<a href="https://deepwiki.com/go-pantheon/lares"><img src="https://deepwiki.com/badge.svg" alt="Ask DeepWiki"></a>
</p>

<p align="center">
  <a href="README.md">English</a> | <a href="README-zh.md">中文</a>
</p>

## About Lares

Lares is the account authentication center of the go-pantheon game server ecosystem, focusing on providing secure and reliable identity verification solutions for game applications. Lares supports native account systems and multiple third-party platform account registration and login, connecting users, third-party platforms, and game services through enterprise-grade security solutions, providing users with Account Tokens for Janus gateway handshake connections to start gaming.

For more information, please visit: [deepwiki/go-pantheon/lares](https://deepwiki.com/go-pantheon/lares)

## About the go-pantheon Ecosystem

**go-pantheon** is an out-of-the-box game server framework that provides high-performance, highly available game server cluster solutions based on microservice architecture. Lares serves as the authentication hub, collaborating with other core services to form a complete game service ecosystem:

- **Roma**: Game core business service, responsible for game logic processing and data management
- **Janus**: Gateway service, responsible for client connection handling and request forwarding
- **Lares**: Account authentication service, responsible for user authentication and account management
- **Senate**: Backend management service, providing operational management interfaces

## Core Advantages

### 🔒 Secure & Reliable
- **Enterprise-grade Encryption**: Multi-layer encryption protection to prevent account theft
- **Replay Attack Prevention**: Ensures uniqueness and security of each request
- **Third-party Platform Verification**: Supports secure integration with mainstream social platforms

### 🌐 Multi-platform Login
- **Traditional Account**: Username/password registration and login
- **Apple Login**: Supports both Web and App platforms
- **Google Login**: Standard OAuth2 flow
- **Facebook Login**: Quick social account integration

### ⚡ High-performance Architecture
- **Microservice Design**: Supports horizontal scaling and distributed deployment
- **Dual Protocol Support**: Provides both gRPC and HTTP interfaces
- **Fast Token Generation**: Efficient token generation and verification mechanisms

### 🎮 Seamless Game Ecosystem Integration
- **Independent Verification**: Janus gateway can verify tokens independently for faster response
- **Operational Support**: Built-in announcement system supporting game operations
- **Flexible Extension**: Modular design, easy to add new features

## System Architecture

### System Architecture Overview

**Component Relationship Diagram:**

```mermaid
graph TB
    Client["🎮 Game Client"]
    Lares["🔐 Lares<br/>(Authentication Service)"]
    Janus["🔰 Janus<br/>(Gateway Service)"]
    Senate["📊 Senate<br/>(Management Service)"]

    subgraph Roma["⚙️ Roma Game Service Cluster"]
        direction LR
        Player["👤 Player Service"]
        Room["🏠 Room Service"]
        Team["🤝 Team Service"]
        Other["... Other Services"]
    end

    subgraph ThirdParty["🌐 Third-party Platforms"]
        direction LR
        Apple["Apple"]
        Google["Google"]
        Facebook["Facebook"]
    end

    subgraph OwnAccount["🔐 Native Account System"]
        direction LR
        Register["📝 Register/Login"]
        Session["🎫 Session Management"]
        Renewal["🔄 Session Renewal"]
    end

    Client -->|"1.Authentication Request"| Lares
    ThirdParty -->|"Token Verification"| Lares
    OwnAccount -->|"Username/Password & Session Verification"| Lares
    Lares -->|"Return Account Token + Session"| Client
    Client -->|"2.Game Connection + Token"| Janus
    Janus -->|"3.gRPC Call"| Player
    Janus -->|"gRPC Call"| Room
    Janus -->|"gRPC Call"| Team
    Player -->|"Business Response"| Janus
    Room -->|"Business Response"| Janus
    Team -->|"Business Response"| Janus
    Janus -->|"4.Game Data"| Client
    Lares -.->|"User Data"| Senate

    classDef clientStyle fill:#e1f5fe,stroke:#01579b,stroke-width:2px
    classDef laresStyle fill:#fff3e0,stroke:#f57c00,stroke-width:3px
    classDef serviceStyle fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    classDef gatewayStyle fill:#e8f5e8,stroke:#1b5e20,stroke-width:2px
    classDef thirdPartyStyle fill:#fce4ec,stroke:#c2185b,stroke-width:2px
    classDef ownAccountStyle fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px

    class Client clientStyle
    class Lares laresStyle
    class Janus gatewayStyle
    class Senate serviceStyle
    class Player,Room,Team,Other serviceStyle
    class Apple,Google,Facebook thirdPartyStyle
    class Register,Session,Renewal ownAccountStyle
```

### Lares Internal Architecture

**Lares Service Architecture Diagram:**

```mermaid
graph TB
    subgraph Lares["🔐 Lares Account Authentication Service"]
        direction TB

        subgraph Services["Microservice Modules"]
            direction LR

            subgraph Account["👤 Account Service"]
                AccountAPI["API Layer"]
                AccountBiz["Business Layer"]
                AccountData["Data Layer"]
            end

            subgraph Notice["📢 Notice Service"]
                NoticeAPI["API Layer"]
                NoticeBiz["Business Layer"]
                NoticeData["Data Layer"]
            end

            subgraph Future["🔮 Future Services"]
                ServerAPI["Server Management API"]
                ServerBiz["Server Management Business"]
                ServerData["Server Management Data"]
            end
        end

        subgraph Security["🛡️ Security Components"]
            direction LR
            JWT["JWT Validator"]
            AES["AES Encryptor"]
            Password["Password Hasher"]
            Session["Session Manager"]
        end

        subgraph Platform["🌐 Platform Integration"]
            direction LR
            AppleDomain["Apple Verification"]
            GoogleDomain["Google Verification"]
            FacebookDomain["Facebook Verification"]
        end
    end

    subgraph Infrastructure["Infrastructure"]
        direction LR
        PostgreSQL["PostgreSQL<br/>User Data Storage"]
        Etcd["etcd<br/>Service Discovery"]
    end

    subgraph External["External Services"]
        direction LR
        AppleAuth["Apple Authentication Service"]
        GoogleAuth["Google Authentication Service"]
        FacebookAuth["Facebook Authentication Service"]
    end

    AccountData -.->|"User Data"| PostgreSQL
    NoticeData -.->|"Notice Data"| PostgreSQL
    AccountAPI -.->|"Service Registration"| Etcd
    NoticeAPI -.->|"Service Registration"| Etcd

    AppleDomain -.->|"Token Verification"| AppleAuth
    GoogleDomain -.->|"Token Verification"| GoogleAuth
    FacebookDomain -.->|"Token Verification"| FacebookAuth

    AccountBiz -.->|"Use"| JWT
    AccountBiz -.->|"Use"| AES
    AccountBiz -.->|"Use"| Password
    AccountBiz -.->|"Use"| Session

    classDef serviceStyle fill:#fff3e0,stroke:#f57c00,stroke-width:2px
    classDef securityStyle fill:#e8f5e8,stroke:#388e3c,stroke-width:2px
    classDef infraStyle fill:#fce4ec,stroke:#c2185b,stroke-width:2px
    classDef externalStyle fill:#f3e5f5,stroke:#4a148c,stroke-width:2px

    class AccountAPI,AccountBiz,AccountData,NoticeAPI,NoticeBiz,NoticeData,ServerAPI,ServerBiz,ServerData serviceStyle
    class JWT,AES,Password,Session,AppleDomain,GoogleDomain,FacebookDomain securityStyle
    class PostgreSQL,Redis,Etcd infraStyle
    class AppleAuth,GoogleAuth,FacebookAuth externalStyle
```

### Authentication Flow Details

**Complete User Authentication Flow:**

```mermaid
sequenceDiagram
    participant C as 🎮 Game Client
    participant L as 🔐 Lares
    participant TP as 🌐 Third-party Platform
    participant OA as 🔐 Native Account System
    participant J as 🔰 Janus
    participant R as ⚙️ Roma Service

    Note over C,R: Method 1: Third-party Platform Login
    C->>L: 1.1 Send third-party login request
    L->>TP: 1.2 Verify third-party token
    TP-->>L: 1.3 Return user information
    L->>L: 1.4 Generate/update user account
    L->>L: 1.5 Generate Account Token + Session
    L-->>C: 1.6 Return Account Token + Session

    Note over C,R: Method 2: Native Account Registration/Login
    C->>L: 2.1 Send register/login request (username + password)
    L->>OA: 2.2 Verify account credentials
    OA-->>L: 2.3 Return verification result
    L->>L: 2.4 Generate/update user account
    L->>L: 2.5 Generate Account Token + Session
    L-->>C: 2.6 Return Account Token + Session

    Note over C,R: Method 3: Session Quick Login
    C->>L: 3.1 Send Session login request
    L->>OA: 3.2 Verify Session validity
    OA-->>L: 3.3 Session verification passed
    L->>L: 3.4 Generate new Account Token
    L-->>C: 3.5 Return Account Token

    Note over C,R: Method 4: Session Renewal (Independent Operation)
    C->>L: 4.1 Send Session renewal request
    L->>OA: 4.2 Verify Session validity
    OA-->>L: 4.3 Session verification passed
    L->>L: 4.4 Extend Session validity
    L-->>C: 4.5 Return new renewed Session

    Note over C,R: Game Connection Phase
    C->>J: 5.1 Establish TCP connection
    C->>J: 5.2 Send handshake request (with Token)
    J->>J: 5.3 Verify Account Token
    J-->>C: 5.4 Handshake success response

    Note over C,R: Game Interaction Phase
    C->>J: 6.1 Send game protocol
    J->>R: 6.2 Forward as gRPC call
    R->>R: 6.3 Process game logic
    R-->>J: 6.4 Return processing result
    J-->>C: 6.5 Forward game response
```

## Core Concepts

### 🔑 Session Mechanism

Session is the player login identifier cached on the client side, used to maintain user login status:

```proto
message Session {
  int64 account_id = 1;    // Account ID
  int64 timeout = 2;       // Expiration timestamp
  string key = 3;          // Random string
}
```

**Features:**
- 🕒 **Auto Expiration**: Ensures security, prevents long-term abuse
- 🔐 **Secure Storage**: Encrypted protection, prevents session hijacking
- 🔄 **Flexible Renewal**: Supports extending validity period
- 📱 **No Repeated Login**: Improves user experience

### 🎫 AuthToken Authentication Token

AuthToken is the verification information carried by the client when handshaking with the Janus gateway, containing complete user identity and routing information:

```proto
message AuthToken {
  string rand = 1;         // Random string for replay attack prevention
  string color = 2;        // Color identifier
  int64 account_id = 3;    // Account ID
  int64 server_id = 4;     // Server ID
  int64 timeout = 5;       // Expiration timestamp
  int32 location = 6;      // Location information
  OnlineStatus status = 7; // Access identifier
  bool unencrypted = 8;    // Whether to disable encryption for this connection
}
```

**Features:**
- 🔒 **Secure Transmission**: Encrypted protection, ensures data security
- 🔄 **Complete Information**: Contains user connection information, no additional queries needed
- 🛡️ **Attack Prevention**: Ensures token uniqueness and security
- 🎯 **Smart Routing**: Supports load balancing and distribution
- ⚡ **Fast Verification**: Gateway independent verification, faster response
- 📍 **Status Tracking**: Supports distributed game logic

## Service Modules

Lares currently supports the following service modules:

| Module      | Status        | Description            | Functions                                                  |
| ----------- | ------------- | ---------------------- | ---------------------------------------------------------- |
| **Account** | ✅ Implemented | Account management     | User registration, login, third-party platform integration |
| **Notice**  | ✅ Implemented | Announcement system    | Notice publishing, management, client push                 |
| **Server**  | 🔮 Planned     | Game server management | Game server creation, management, load balancing           |

### Authentication Support

| Authentication Method | Status         | Description                   | Features                                        |
| --------------------- | -------------- | ----------------------------- | ----------------------------------------------- |
| **Username/Password** | ✅ Full Support | Traditional account reg/login | Secure encryption, password strength validation |
| **Apple Sign In**     | ✅ Full Support | Apple official login          | Supports Web/App dual platform                  |
| **Google OAuth**      | ✅ Full Support | Google account login          | Standard OAuth2 flow                            |
| **Facebook Login**    | ✅ Full Support | Facebook social login         | Quick social access                             |

## Technology Stack

Lares uses the following core technologies:

| Technology/Component | Purpose                      | Version  |
| -------------------- | ---------------------------- | -------- |
| **Go**               | Primary development language | 1.24+    |
| **go-kratos**        | Microservice framework       | v2.8.4   |
| **gRPC**             | Inter-service communication  | v1.73.0  |
| **Protobuf**         | Data serialization           | v1.36.6  |
| **etcd**             | Service discovery & registry | v3.6.1   |
| **PostgreSQL**       | Main database                | v5.7.5   |
| **OpenTelemetry**    | Distributed tracing          | v1.37.0  |
| **Prometheus**       | Monitoring system            | v1.22.0  |
| **Google Wire**      | Dependency injection         | v0.6.0   |
| **JWT**              | Token verification           | v4.5.2   |
| **Argon2**           | Password hashing             | Built-in |

## Quick Start

### Requirements

- **Go 1.24+** - Primary development language
- **PostgreSQL 13+** - Main database
- **etcd 3.5+** - Service discovery and configuration center
- **protoc** - Protocol Buffers compiler

### Installation

```bash
# 1. Clone the project
git clone https://github.com/go-pantheon/lares.git
cd lares

# 2. Initialize development environment
make init

# 3. Install dependencies
go mod download
```

### Configuration

```bash
# 1. Copy configuration template
cp app/account/configs.tmpl/config.yaml app/account/configs/config.yaml

# 2. Edit configuration file
vim app/account/configs/config.yaml
```

**Key Configuration Items:**

```yaml
# Service configuration
server:
  http:
    addr: 0.0.0.0:8001
  grpc:
    addr: 0.0.0.0:9001

# Database configuration
data:
  postgresql:
    source: "postgres://user:password@localhost:5432/lares?sslmode=disable"

# Third-party platform configuration
platform:
  apple:
    client_id: "your_apple_client_id"
    team_id: "your_apple_team_id"
    key_id: "your_apple_key_id"
  google:
    aud: "your_google_client_id"
  facebook:
    app_id: "your_facebook_app_id"
    app_secret: "your_facebook_app_secret"

# Security configuration
secret:
  token_key: "your_32_byte_token_encryption_key"
  session_key: "your_32_byte_session_encryption_key"
  platform_key: "your_32_byte_platform_encryption_key"
```

### Start Services

```bash
# 1. Generate code
make generate

# 2. Build service
make build

# 3. Start account service
make run app=account

# 4. Start notice service
make run app=notice
```

### Service Verification

```bash
# Check service status
curl http://localhost:8001/accounts/v1/dev/ping

# Test user registration
curl -X POST http://localhost:8001/accounts/v1/username/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "SecurePassword123!",
    "color": "blue"
  }'
```

## API Documentation

### Account Management Interface

#### User Registration
```http
POST /accounts/v1/username/register
Content-Type: application/json

{
  "username": "player123",
  "password": "SecurePassword123!",
  "color": "blue"
}

# Response Example
{
  "info": {
    "account_id": "encoded_account_id",
    "register": true,
    "token": "encrypted_auth_token",      # AuthToken, for gateway handshake
    "token_timeout": 1704067200,
    "session": "encrypted_session",       # Session, client cache
    "session_timeout": 1704153600,
    "state": "random_state_string"
  }
}
```

#### User Login
```http
POST /accounts/v1/username/login
Content-Type: application/json

{
  "username": "player123",
  "password": "SecurePassword123!",
  "color": "blue"
}

# Response Example
{
  "info": {
    "account_id": "encoded_account_id",
    "register": false,
    "token": "encrypted_auth_token",      # AuthToken, for gateway handshake
    "token_timeout": 1704067200,
    "session": "encrypted_session",       # Session, client cache
    "session_timeout": 1704153600,
    "state": "random_state_string"
  }
}
```

#### Third-party Login
```http
# Apple Login
POST /accounts/v1/apple/login
Content-Type: application/json

{
  "token": "apple_id_token",
  "color": "blue"
}

# Google Login
POST /accounts/v1/google/login
Content-Type: application/json

{
  "token": "google_id_token",
  "color": "blue"
}

# Facebook Login
POST /accounts/v1/fb/login
Content-Type: application/json

{
  "token": "facebook_access_token",
  "color": "blue"
}
```

#### Get Game Token
```http
POST /accounts/v1/token
Content-Type: application/json

{
  "account_id": "encoded_account_id",
  "session": "encrypted_session_token",
  "color": "blue"
}
```

#### Refresh Session
```http
POST /accounts/v1/refresh
Content-Type: application/json

{
  "account_id": "encoded_account_id",
  "session": "current_session_token"
}
```

### Notice System Interface

#### Get Notice List
```http
GET /notices/v1/list
Content-Type: application/json

{
  "page": 1,
  "size": 10
}
```

## Contributing

We welcome contributions! Please follow the process below:

1. Fork the project
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Create a Pull Request

### Development Standards

- Follow Go official code standards
- Use `golangci-lint` for code checking
- Write unit tests, ensure code coverage > 80%
- Update relevant documentation and API descriptions
- Ensure CI/CD pipeline passes

## License

This project is open sourced under the [MIT License](https://github.com/go-pantheon/lares/blob/main/LICENSE).

---

<div align="center">
✨ **Lares Game Account Authentication Service** - Your passport to enter the gaming world

🏛️ _Part of the go-pantheon ecosystem_ 🏛️
</div>
