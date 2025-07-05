<div align="center">
  <h1>🔐 Lares 游戏账号鉴权服务</h1>
  <p><em>专业游戏账号认证服务，具备双令牌架构和多平台登录支持</em></p>
</div>

<p align="center">
<a href="https://github.com/go-pantheon/lares/actions"><img src="https://github.com/go-pantheon/lares/workflows/Go/badge.svg" alt="构建状态"></a>
<a href="https://github.com/go-pantheon/lares/releases"><img src="https://img.shields.io/github/v/release/go-pantheon/lares" alt="最新版本"></a>
<a href="https://github.com/go-kratos/kratos"><img src="https://img.shields.io/badge/Kratos-v2.8.4-blue" alt="Kratos"></a>
<a href="https://pkg.go.dev/github.com/go-pantheon/lares"><img src="https://pkg.go.dev/badge/github.com/go-pantheon/lares" alt="GoDoc"></a>
<a href="https://goreportcard.com/report/github.com/go-pantheon/lares"><img src="https://goreportcard.com/badge/github.com/go-pantheon/lares" alt="Go Report Card"></a>
<a href="https://github.com/go-pantheon/lares/blob/main/LICENSE"><img src="https://img.shields.io/github/license/go-pantheon/lares" alt="许可证"></a>
<a href="https://deepwiki.com/go-pantheon/lares"><img src="https://deepwiki.com/badge.svg" alt="Ask DeepWiki"></a>
</p>

<p align="center">
  <a href="README.md">English</a> | <a href="README-zh.md">中文</a>
</p>

## 关于 Lares

**Lares** 是一个专业的游戏账号认证服务，为现代游戏应用提供安全、高性能的身份验证解决方案。Lares 采用双令牌架构（Session + AuthToken），原生支持 Apple/Google/Facebook 登录，并具备游戏优化功能，如基于染色选择和服务器路由等特性，将复杂的认证挑战转化为简单的解决方案，使游戏开发者能够快速实现身份验证，同时保持企业级的安全性和大规模性能。

## 关于 go-pantheon 生态系统

**go-pantheon** 是一个开箱即用的游戏服务器框架，基于微服务架构提供高性能、高可用的游戏服务器集群解决方案。Lares 作为身份验证中心，与其他核心服务协作构成完整的游戏服务生态系统：

- **Roma**: 游戏核心业务服务，负责游戏逻辑处理和数据管理
- **Janus**: 网关服务，负责客户端连接处理和请求转发
- **Lares**: 账户鉴权服务，负责用户认证和账户管理
- **Senate**: 后台管理服务，提供运营管理接口

## 核心优势

### 🔒 安全可靠
- **企业级加密**: 多重加密保护，防止账号被盗
- **防重放攻击**: 确保每次请求的唯一性和安全性
- **第三方平台验证**: 支持主流社交平台安全接入

### 🌐 多平台登录
- **传统账号**: 用户名/密码注册登录
- **Apple 登录**: 支持 Web 和 App 双端
- **Google 登录**: 标准 OAuth2 流程
- **Facebook 登录**: 社交账户快速接入

### ⚡ 高性能架构
- **微服务设计**: 支持水平扩展和分布式部署
- **双协议支持**: 同时提供 gRPC 和 HTTP 接口
- **快速Token生成**: 高效的令牌生成和验证机制

### 🎮 游戏生态无缝集成
- **独立验证**: Janus网关可自行验证Token，响应更快
- **运营支撑**: 内置公告系统，支持游戏运营
- **灵活扩展**: 模块化设计，易于添加新功能

## 系统架构

### 系统架构总览

**组件关系图：**

```mermaid
graph TB
    Client["🎮 游戏客户端"]
    Lares["🔐 Lares<br/>(账号鉴权服务)"]
    Janus["🔰 Janus<br/>(网关服务)"]
    Senate["📊 Senate<br/>(管理服务)"]

    subgraph Roma["⚙️ Roma 游戏服务集群"]
        direction LR
        Player["👤 Player服务"]
        Room["🏠 Room服务"]
        Team["🤝 Team服务"]
        Other["... 其他服务"]
    end

    subgraph ThirdParty["🌐 第三方平台"]
        direction LR
        Apple["Apple"]
        Google["Google"]
        Facebook["Facebook"]
    end

    subgraph OwnAccount["🔐 自有账号系统"]
        direction LR
        Register["📝 注册登录"]
        Session["🎫 Session管理"]
        Renewal["🔄 Session续期"]
    end

    Client -->|"1.身份验证请求"| Lares
    ThirdParty -->|"令牌验证"| Lares
    OwnAccount -->|"账号密码/Session验证"| Lares
    Lares -->|"返回Account Token + Session"| Client
    Client -->|"2.游戏连接+Token"| Janus
    Janus -->|"3.gRPC 调用"| Player
    Janus -->|"gRPC 调用"| Room
    Janus -->|"gRPC 调用"| Team
    Player -->|"业务响应"| Janus
    Room -->|"业务响应"| Janus
    Team -->|"业务响应"| Janus
    Janus -->|"4.游戏数据"| Client
    Lares -.->|"用户数据"| Senate

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

### Lares 内部架构

**Lares 服务架构图：**

```mermaid
graph TB
    subgraph Lares["🔐 Lares 账号鉴权服务"]
        direction TB

        subgraph Services["微服务模块"]
            direction LR

            subgraph Account["👤 Account 服务"]
                AccountAPI["API 层"]
                AccountBiz["业务层"]
                AccountData["数据层"]
            end

            subgraph Notice["📢 Notice 服务"]
                NoticeAPI["API 层"]
                NoticeBiz["业务层"]
                NoticeData["数据层"]
            end

            subgraph Future["🔮 未来服务"]
                ServerAPI["区服管理 API"]
                ServerBiz["区服管理业务"]
                ServerData["区服管理数据"]
            end
        end

        subgraph Security["🛡️ 安全组件"]
            direction LR
            JWT["JWT 验证器"]
            AES["AES 加密器"]
            Password["密码哈希器"]
            Session["会话管理器"]
        end

        subgraph Platform["🌐 平台集成"]
            direction LR
            AppleDomain["Apple 验证"]
            GoogleDomain["Google 验证"]
            FacebookDomain["Facebook 验证"]
        end
    end

    subgraph Infrastructure["基础设施"]
        direction LR
        PostgreSQL["PostgreSQL<br/>用户数据存储"]
        Etcd["etcd<br/>服务发现"]
    end

    subgraph External["外部服务"]
        direction LR
        AppleAuth["Apple 认证服务"]
        GoogleAuth["Google 认证服务"]
        FacebookAuth["Facebook 认证服务"]
    end

    AccountData -.->|"用户数据"| PostgreSQL
    NoticeData -.->|"公告数据"| PostgreSQL
    AccountAPI -.->|"服务注册"| Etcd
    NoticeAPI -.->|"服务注册"| Etcd

    AppleDomain -.->|"令牌验证"| AppleAuth
    GoogleDomain -.->|"令牌验证"| GoogleAuth
    FacebookDomain -.->|"令牌验证"| FacebookAuth

    AccountBiz -.->|"使用"| JWT
    AccountBiz -.->|"使用"| AES
    AccountBiz -.->|"使用"| Password
    AccountBiz -.->|"使用"| Session

    classDef serviceStyle fill:#fff3e0,stroke:#f57c00,stroke-width:2px
    classDef securityStyle fill:#e8f5e8,stroke:#388e3c,stroke-width:2px
    classDef infraStyle fill:#fce4ec,stroke:#c2185b,stroke-width:2px
    classDef externalStyle fill:#f3e5f5,stroke:#4a148c,stroke-width:2px

    class AccountAPI,AccountBiz,AccountData,NoticeAPI,NoticeBiz,NoticeData,ServerAPI,ServerBiz,ServerData serviceStyle
    class JWT,AES,Password,Session,AppleDomain,GoogleDomain,FacebookDomain securityStyle
    class PostgreSQL,Redis,Etcd infraStyle
    class AppleAuth,GoogleAuth,FacebookAuth externalStyle
```

### 认证流程详解

**完整的用户认证流程：**

```mermaid
sequenceDiagram
    participant C as 🎮 游戏客户端
    participant L as 🔐 Lares
    participant TP as 🌐 第三方平台
    participant OA as 🔐 自有账号系统
    participant J as 🔰 Janus
    participant R as ⚙️ Roma服务

    Note over C,R: 方式一：第三方平台登录
    C->>L: 1.1 发送第三方登录请求
    L->>TP: 1.2 验证第三方令牌
    TP-->>L: 1.3 返回用户信息
    L->>L: 1.4 生成/更新用户账户
    L->>L: 1.5 生成Account Token + Session
    L-->>C: 1.6 返回Account Token + Session

    Note over C,R: 方式二：自有账号注册/登录
    C->>L: 2.1 发送注册/登录请求(用户名+密码)
    L->>OA: 2.2 验证账号密码
    OA-->>L: 2.3 返回验证结果
    L->>L: 2.4 生成/更新用户账户
    L->>L: 2.5 生成Account Token + Session
    L-->>C: 2.6 返回Account Token + Session

    Note over C,R: 方式三：Session快速登录
    C->>L: 3.1 发送Session登录请求
    L->>OA: 3.2 验证Session有效性
    OA-->>L: 3.3 Session验证通过
    L->>L: 3.4 生成新的Account Token
    L-->>C: 3.5 返回Account Token

    Note over C,R: 方式四：Session续期(独立操作)
    C->>L: 4.1 发送Session续期请求
    L->>OA: 4.2 验证Session有效性
    OA-->>L: 4.3 Session验证通过
    L->>L: 4.4 延长Session有效期
    L-->>C: 4.5 返回新的续期Session

    Note over C,R: 游戏连接阶段
    C->>J: 5.1 建立TCP连接
    C->>J: 5.2 发送握手请求(含Token)
    J->>J: 5.3 验证Account Token
    J-->>C: 5.4 握手成功响应

    Note over C,R: 游戏交互阶段
    C->>J: 6.1 发送游戏协议
    J->>R: 6.2 转发为gRPC调用
    R->>R: 6.3 处理游戏逻辑
    R-->>J: 6.4 返回处理结果
    J-->>C: 6.5 转发游戏响应
```

## 核心概念

### 🔑 Session 会话机制

Session 是客户端缓存的玩家登录标识，用于维护用户的登录状态：

```proto
message Session {
  int64 account_id = 1;    // 账户ID
  int64 timeout = 2;       // 过期时间戳
  string key = 3;          // 随机字符串
}
```

**特性：**
- 🕒 **自动过期**: 确保安全，防止长期滥用
- 🔐 **安全存储**: 加密保护，防止会话劫持
- 🔄 **灵活续期**: 支持延长有效期
- 📱 **免重复登录**: 提升用户体验

### 🎫 AuthToken 认证令牌

AuthToken 是客户端与 Janus 网关握手时携带的校验信息，包含完整的用户身份和路由信息：

```proto
message AuthToken {
  string rand = 1;         // 防重放攻击的随机字符串
  string color = 2;        // 染色标识
  int64 account_id = 3;    // 账户ID
  int64 server_id = 4;     // 服务器ID
  int64 timeout = 5;       // 过期时间戳
  int32 location = 6;      // 位置信息
  OnlineStatus status = 7; // 访问标识
  bool unencrypted = 8;    // 本次连接是否关闭加密
}
```

**特性：**
- 🔒 **安全传输**: 加密保护，确保数据安全
- 🔄 **信息完整**: 包含用户连接信息，无需额外查询
- 🛡️ **防攻击**: 确保令牌唯一性和安全性
- 🎯 **智能路由**: 支持负载均衡和分发
- ⚡ **快速验证**: 网关独立验证，响应更快
- 📍 **状态追踪**: 支持分布式游戏逻辑

## 服务模块

Lares 当前支持以下服务模块：

| 模块        | 状态     | 描述         | 功能                           |
| ----------- | -------- | ------------ | ------------------------------ |
| **Account** | ✅ 已实现 | 账户管理服务 | 用户注册、登录、第三方平台集成 |
| **Notice**  | ✅ 已实现 | 公告系统服务 | 公告发布、管理、客户端推送     |
| **Server**  | 🔮 规划中 | 区服管理服务 | 游戏区服创建、管理、负载均衡   |

### 身份验证支持

| 认证方式           | 状态       | 描述             | 特性                   |
| ------------------ | ---------- | ---------------- | ---------------------- |
| **用户名/密码**    | ✅ 完整支持 | 传统账号注册登录 | 安全加密、密码强度验证 |
| **Apple Sign In**  | ✅ 完整支持 | Apple官方登录    | 支持Web/App双端        |
| **Google OAuth**   | ✅ 完整支持 | Google账号登录   | 标准OAuth2流程         |
| **Facebook Login** | ✅ 完整支持 | Facebook社交登录 | 快速社交接入           |

## 技术栈

Lares 使用以下核心技术：

| 技术/组件         | 用途         | 版本     |
| ----------------- | ------------ | -------- |
| **Go**            | 主要开发语言 | 1.24+    |
| **go-kratos**     | 微服务框架   | v2.8.4   |
| **gRPC**          | 服务间通信   | v1.73.0  |
| **Protobuf**      | 数据序列化   | v1.36.6  |
| **etcd**          | 服务发现注册 | v3.6.1   |
| **PostgreSQL**    | 主数据库     | v5.7.5   |
| **OpenTelemetry** | 分布式追踪   | v1.37.0  |
| **Prometheus**    | 监控系统     | v1.22.0  |
| **Google Wire**   | 依赖注入     | v0.6.0   |
| **JWT**           | 令牌验证     | v4.5.2   |
| **Argon2**        | 密码哈希     | 内置支持 |

## 快速开始

### 环境要求

- **Go 1.24+** - 主要开发语言
- **PostgreSQL 13+** - 主数据库
- **etcd 3.5+** - 服务发现与配置中心
- **protoc** - Protocol Buffers 编译器

### 安装

```bash
# 1. 克隆项目
git clone https://github.com/go-pantheon/lares.git
cd lares

# 2. 初始化开发环境
make init

# 3. 安装依赖
go mod download
```

### 配置

```bash
# 1. 复制配置模板
cp app/account/configs.tmpl/config.yaml app/account/configs/config.yaml

# 2. 编辑配置文件
vim app/account/configs/config.yaml
```

**关键配置项：**

```yaml
# 服务配置
server:
  http:
    addr: 0.0.0.0:8001
  grpc:
    addr: 0.0.0.0:9001

# 数据库配置
data:
  postgresql:
    source: "postgres://user:password@localhost:5432/lares?sslmode=disable"

# 第三方平台配置
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

# 安全配置
secret:
  token_key: "your_32_byte_token_encryption_key"
  session_key: "your_32_byte_session_encryption_key"
  platform_key: "your_32_byte_platform_encryption_key"
```

### 启动服务

```bash
# 1. 生成代码
make generate

# 2. 构建服务
make build

# 3. 启动账户服务
make run app=account

# 4. 启动公告服务
make run app=notice
```

### 验证服务

```bash
# 检查服务状态
curl http://localhost:8001/accounts/v1/dev/ping

# 测试用户注册
curl -X POST http://localhost:8001/accounts/v1/username/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "SecurePassword123!",
    "color": "blue"
  }'
```

## API 文档

### 账户管理接口

#### 用户注册
```http
POST /accounts/v1/username/register
Content-Type: application/json

{
  "username": "player123",
  "password": "SecurePassword123!",
  "color": "blue"
}

# 响应示例
{
  "info": {
    "account_id": "encoded_account_id",
    "register": true,
    "token": "encrypted_auth_token",      # AuthToken，用于网关握手
    "token_timeout": 1704067200,
    "session": "encrypted_session",       # Session，客户端缓存
    "session_timeout": 1704153600,
    "state": "random_state_string"
  }
}
```

#### 用户登录
```http
POST /accounts/v1/username/login
Content-Type: application/json

{
  "username": "player123",
  "password": "SecurePassword123!",
  "color": "blue"
}

# 响应示例
{
  "info": {
    "account_id": "encoded_account_id",
    "register": false,
    "token": "encrypted_auth_token",      # AuthToken，用于网关握手
    "token_timeout": 1704067200,
    "session": "encrypted_session",       # Session，客户端缓存
    "session_timeout": 1704153600,
    "state": "random_state_string"
  }
}
```

#### 第三方登录
```http
# Apple 登录
POST /accounts/v1/apple/login
Content-Type: application/json

{
  "token": "apple_id_token",
  "color": "blue"
}

# Google 登录
POST /accounts/v1/google/login
Content-Type: application/json

{
  "token": "google_id_token",
  "color": "blue"
}

# Facebook 登录
POST /accounts/v1/fb/login
Content-Type: application/json

{
  "token": "facebook_access_token",
  "color": "blue"
}
```

#### 获取游戏令牌
```http
POST /accounts/v1/token
Content-Type: application/json

{
  "account_id": "encoded_account_id",
  "session": "encrypted_session_token",
  "color": "blue"
}
```

#### 刷新会话
```http
POST /accounts/v1/refresh
Content-Type: application/json

{
  "account_id": "encoded_account_id",
  "session": "current_session_token"
}
```

### 公告系统接口

#### 获取公告列表
```http
GET /notices/v1/list
Content-Type: application/json

{
  "page": 1,
  "size": 10
}
```

## 贡献指南

欢迎贡献代码！请遵循以下流程：

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

### 开发规范

- 遵循 Go 语言官方代码规范
- 使用 `golangci-lint` 进行代码检查
- 编写单元测试，保证代码覆盖率 > 80%
- 更新相关文档和 API 说明
- 确保 CI/CD 流程通过

## 许可证

本项目基于 [MIT License](https://github.com/go-pantheon/lares/blob/main/LICENSE) 开源协议。

---

<div align="center">
✨ **Lares 游戏账号鉴权服务** - 用户进入游戏世界的通行证

🏛️ _go-pantheon 生态系统的一部分_ 🏛️
</div>
