# Developer Guide

Welcome to the Echo REST API project! This guide will help you get started with local development.

## Prerequisites

- **Go 1.21+**
- **Docker & Docker Compose** - for PostgreSQL
- **Firebase credentials** - You'll need a Firebase project with Authentication enabled

## Local Development Setup

### 1. Install Dependencies

```bash
make setup
```

This installs:
- `swag` - Swagger documentation generator
- `migrate` - Database migration tool

### 2. Configure Environment

```bash
# The app reads from $HOME/.env by default when using make run-api
cp .env.example $HOME/.env
```

Edit `$HOME/.env` with your configuration:

```env
DATABASE_URL=postgres://world:hello@localhost:5432/echo_rest_api
PORT=8088
AUTH_PROVIDER=firebase
FIREBASE_CREDENTIALS={"type": "service_account", ...}
LOG_LEVEL=debug
```

### 3. Start Database

```bash
make run-db
```

### 4. Run Migrations

```bash
make migration-up
```

### 5. Generate API Docs (Optional)

```bash
make docs
```

### 6. Start the API

```bash
make run-api
```

API runs at `http://localhost:8088`
Swagger docs at `http://localhost:8088/docs/index.html`

## Architecture Overview

### Technology Stack

| Component | Technology |
|-----------|------------|
| Framework | Echo v4 |
| Auth | Firebase JWT |
| RBAC | Open Policy Agent (OPA) |
| Database | PostgreSQL with sqlx |
| DI Container | Uber-go/dig |
| Migrations | golang-migrate |

### Project Structure

```
cmd/api/          # Application entry point
config/           # Application configuration
infrastructure/   # Database connections, external services
migrations/sql/   # SQL migration files
modules/          # Feature modules (core, projects, etc.)
  ├── domains/    # Domain models/entities
  ├── dto/        # Request/Response DTOs
  ├── repositories/  # Data access layer
  ├── usecases/   # Business logic
  └── handlers/   # HTTP handlers
pkg/
  ├── authz/      # OPA Rego policies
  ├── middlewares/    # Auth middleware
  ├── wrapper/    # Response wrapper
  └── utils/      # Helper functions
```

### Authorization Flow

```
Request → Auth Middleware (Firebase JWT) → CheckPolicies Middleware (OPA) → Handler
```

1. **Auth Middleware** validates Firebase JWT and loads user context
2. **CheckPolicies Middleware** evaluates OPA Rego policies based on user roles
3. Endpoint permissions are defined in `pkg/authz/routes.json`

## Module System

Each module implements the `ModuleInstance` interface:

```go
type ModuleInstance interface {
    RegisterRepositories(container *dig.Container) error
    RegisterUseCases(container *dig.Container) error
    RegisterHandlers(g *echo.Group, container *dig.Container) error
}
```

Modules are registered in `cmd/api/di/di.go`:

```go
mapModules := map[string]core.ModuleInstance{
    "core":     core.Module,
    "projects": projects.Module,
}
```

## Creating a New Module

### Using the Generator

```bash
go run ./tools/mod/ gen -n ModuleName
```

This creates:
- `modules/modulename/mod.go` - Module registration
- `modules/modulename/domains/` - Domain models
- `modules/modulename/dto/` - Request/Response DTOs
- `modules/modulename/repositories/` - Data access
- `modules/modulename/usecases/` - Business logic
- `modules/modulename/handlers/` - HTTP handlers

### Manual Creation

1. Create module structure under `modules/yourmodule/`
2. Implement domain models in `domains/`
3. Implement repositories in `repositories/`
4. Implement use cases in `usecases/`
5. Implement handlers in `handlers/`
6. Register in `cmd/api/di/di.go`

### After Creating a Module

Always run:

```bash
make routes
```

This updates `pkg/authz/routes.json` with new endpoint permissions for OPA.

## API Development

### Response Format

Use `wrapper.Wrap()` for standardized responses:

```go
package handlers

import "github.com/dzungtran/echo-rest-api/pkg/wrapper"

type Handler struct{}

func (h *Handler) GetUser(c echo.Context) wrapper.Response {
    user := getUserFromDB()
    return wrapper.Response{
        Data:   user,
        Status: 200,
    }
}

// In handler registration:
handlers := []*Handler{h}
for _, h := range handlers {
    e.GET("/users/:id", wrapper.Wrap(h.GetUser), middManager.Auth(), middManager.CheckPolicies())
}
```

### Response Structure

```json
{
  "success": true,
  "data": { ... },
  "metadata": {
    "total": 100
  }
}
```

### Error Response

```go
return wrapper.Response{
    Error:  errors.New("user not found"),
    Status: 404,
}
```

### Request DTOs

Place DTOs in `modules/yourmodule/dto/requests.go`:

```go
package dto

type CreateUserRequest struct {
    Email string `json:"email" validate:"required,email"`
    Name  string `json:"name" validate:"required"`
}
```

## Database Migrations

### Create Migration

```bash
make migration-create name=add_users_table
```

Creates files:
- `migrations/sql/20220315000000_add_users_table.up.sql`
- `migrations/sql/20220315000000_add_users_table.down.sql`

### Run Migrations

```bash
make migration-up
```

### Rollback

```bash
make migration-down
```

### Migration Content Example

```sql
-- up
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- down
DROP TABLE users;
```

## Authentication

### Firebase JWT

1. Client sends `Authorization: Bearer <token>` header
2. `Auth()` middleware validates token with Firebase
3. User is loaded and set in context

### Getting Current User in Handlers

```go
func (h *Handler) GetUser(c echo.Context) wrapper.Response {
    user := c.Get(constants.ContextKeyUser).(*domains.User)
    // ...
}
```

### Context Keys

Defined in `pkg/constants/context.go`:
- `ContextKeyUser` - Current authenticated user
- `ContextKeyOrg` - Current organization (for org-scoped endpoints)
- `ContextKeyProject` - Current project

## Authorization (OPA Policies)

### Policy Files

- `pkg/authz/data.json` - Static policy data (roles, permissions)
- `pkg/authz/routes.json` - Endpoint definitions
- Rego policies evaluated at runtime

### How It Works

1. `routes.json` defines which roles can access which endpoints
2. OPA evaluates policies based on:
   - User's roles (from JWT)
   - Request context (org, project ownership)
   - Endpoint being accessed

### Adding New Endpoints

After creating new handlers:

#### 1. Register the Route

Routes are registered with an explicit **route name** in `<action_name>:<resource_name>` format:

```go
// In your handler registration (modules/yourmodule/handlers/)
apiV1 := g.Group("admin/resources", middManager.Auth(), middManager.CheckPolicies())
apiV1.GET("", wrapper.Wrap(handler.Fetch)).Name = "list:resource"
apiV1.POST("", wrapper.Wrap(handler.Create)).Name = "create:resource"

apiV1Resource := g.Group("admin/resources/:resourceId",
    middManager.Auth(),
    middlewares.RequireResourceIdInParam("resourceId"),
    middManager.CheckPoliciesWithOrg(), // or CheckPoliciesWithProject()
)
apiV1Resource.GET("", wrapper.Wrap(handler.GetByID)).Name = "read:resource"
apiV1Resource.PUT("", wrapper.Wrap(handler.Update)).Name = "update:resource"
apiV1Resource.DELETE("", wrapper.Wrap(handler.Delete)).Name = "delete:resource"
```

**Route name format:** `<action_name>:<resource_name>` (e.g., `list:project`, `create:org`)
- Actions: `list`, `create`, `read`, `update`, `delete`
- Resource: singular noun (e.g., `org`, `project`, `resource`)

#### 2. Generate Routes

```bash
make routes
```

This introspects all Echo routes and writes them to `pkg/authz/routes.json`:

```json
{
  "/admin/resources": {
    "GET": "list:resource",
    "POST": "create:resource"
  },
  "/admin/resources/:resourceId": {
    "DELETE": "delete:resource",
    "GET": "read:resource",
    "PUT": "update:resource"
  }
}
```

The route **name** (`list:resource`, `create:resource`, etc.) is the permission identifier used in OPA policies.

#### 3. Add Permissions to `data.json`

Edit `pkg/authz/data.json` to grant roles access to the new permissions:

```json
{
  "endpoints_acl": { ... },
  "roles_chart": {
    "owner": {
      "access": ["delete:org", "update:org", "delete:resource", "update:resource"]
    },
    "manager": {
      "access": ["invite:org", "create:resource", "update:resource", "delete:resource"],
      "owner": "owner"
    },
    "viewer": {
      "access": ["read:org", "list:org", "list:resource", "read:resource"],
      "owner": "manager"
    }
  }
}
```

#### 4. Role Hierarchy

Roles inherit from `owner` chain. A `viewer` has `owner: "manager"` which means they inherit all `manager` permissions, and `manager` has `owner: "owner"`.

### Testing Policies

```bash
# Test if a user can access an endpoint
opa eval --bundle pkg/authz/ "data.authz.allow" -i test_input.json
```

## Testing

### Run All Tests

```bash
make test
```

### Run Tests for Specific Package

```bash
go test -v ./modules/core/...
```

### Write Tests

```go
package usecases

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestUserUsecase_CreateUser(t *testing.T) {
    uc := NewUserUsecase(mockRepo)
    user, err := uc.CreateUser(context.Background(), &dto.CreateUserRequest{
        Email: "test@example.com",
        Name:  "Test User",
    })
    assert.NoError(t, err)
    assert.NotNil(t, user)
}
```

## Common Tasks

### Add a New API Endpoint

1. Create/Update handler in `modules/yourmodule/handlers/`
2. Use `wrapper.Wrap()` for response
3. Register route with explicit route name in `<action_name>:<resource_name>` format:
   ```go
   apiV1 := g.Group("admin/resources", middManager.Auth(), middManager.CheckPolicies())
   apiV1.GET("", wrapper.Wrap(handler.Fetch)).Name = "list:resource"
   ```
4. Run `make docs` to update Swagger
5. Run `make routes` to update `pkg/authz/routes.json`
6. Add permission (e.g., `list:resource`) to the appropriate roles in `pkg/authz/data.json`

### Add a New Database Field

1. Create migration:
   ```bash
   make migration-create name=add_field_to_users
   ```
2. Edit migration files
3. Run `make migration-up`
4. Update domain model in `domains/`
5. Update repository queries if needed

### Add Environment Variable

1. Add to `config/app_config.go`
2. Add to `.env.example`
3. Document in this DEVELOP.md

## Troubleshooting

### Database Connection Issues

```bash
# Check if PostgreSQL is running
docker compose ps postgres

# Test connection
psql $DATABASE_URL -c "SELECT 1"
```

### Migration Failures

```bash
# Check migration status
migrate -path migrations/sql -database $DATABASE_URL version

# Force migrate to specific version
migrate -path migrations/sql -database $DATABASE_URL force 20220315000000
```

### Firebase Auth Issues

1. Verify `FIREBASE_CREDENTIALS` is valid JSON
2. Check Firebase console - ensure Authentication is enabled
3. Ensure the token is not expired

### Swagger Not Updating

```bash
make docs
# Restart the API server
```

### OPA Policy Issues

If access is denied unexpectedly:
1. Run `make routes` to regenerate `pkg/authz/routes.json`
2. Check that the route name (e.g., `list:resource`, `create:resource`) appears in `routes.json` for your endpoint
3. Verify the permission is granted in `pkg/authz/data.json` under the appropriate role
4. Ensure the user's role includes the required permission (check role hierarchy)

## Useful Commands

| Command | Description |
|---------|-------------|
| `make run-api` | Start API server |
| `make build-api` | Build binary to `./bin/` |
| `make test` | Run tests with coverage |
| `make docs` | Generate Swagger docs |
| `make routes` | Update OPA routes |
| `make migration-up` | Run migrations |
| `make migration-down` | Rollback last migration |
| `make migration-create name=x` | Create new migration |

## Additional Resources

- [Echo Framework](https://echo.labstack.com/guide/)
- [OPA Documentation](https://www.openpolicyagent.org/docs/latest/)
- [Firebase Auth](https://firebase.google.com/docs/auth/admin/)
- [golang-migrate](https://github.com/golang-migrate/migrate)
