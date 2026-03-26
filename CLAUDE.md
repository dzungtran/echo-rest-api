# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run Commands

```bash
make run-api          # Run the API (uses $HOME/.env by default)
make build-api        # Build binary to ./bin/
make test             # Run tests with coverage and benchmarks
make docs             # Generate Swagger docs (swag i)
make run-db           # Start PostgreSQL via docker compose
make migration-up     # Run database migrations
make migration-down   # Rollback last migration
make migration-create name=<name>  # Create new migration
make setup            # Install migrate and swag CLI tools
make routes           # Generate routes.json for OPA authorization
```

## Architecture

**Framework:** Echo v4 with Firebase Auth for JWT authentication and Open Policy Agent (OPA) for RBAC.

**Module System:** Each module (core, projects, etc.) implements `ModuleInstance` interface with `RegisterRepositories`, `RegisterUseCases`, and `RegisterHandlers`. Modules are registered in `cmd/api/di/di.go` via uber-go/dig container.

**Authorization Flow:** `Auth()` middleware validates Firebase JWT → `CheckPolicies()` middleware evaluates OPA Rego policies with user roles. Endpoint permissions are defined in `pkg/authz/routes.json`.

**Database:** PostgreSQL with sqlx, using master/slave instance pattern. Migrations in `migrations/sql/`.

**Response Format:** Use `wrapper.Wrap(handlerFunc)` for standardized JSON responses with `Data`, `Error`, `Status`, `Total`, `IncludeTotal` fields.

## Generate New Module

```bash
go run ./tools/mod/ gen -n <ModuleName>
```

This auto-generates domain, repository, usecase, and handler files. Always run `make routes` after to update authorization.

## Environment Variables

Key vars: `DATABASE_URL`, `PORT`, `AUTH_PROVIDER=firebase`, `FIREBASE_CREDENTIALS`, `AUTO_MIGRATE=true`, `LOG_LEVEL`. App loads `.env` from `$HOME/.env` when running with `make run-api`.

## Key Dependencies

- `github.com/labstack/echo/v4` - Web framework
- `github.com/open-policy-agent/opa` - Policy enforcement
- `github.com/uber-go/dig` - Dependency injection
- `github.com/firebase/firebase-admin-go/v4` - Authentication
- `github.com/jmoiron/sqlx` - Database
- `github.com/golang-migrate/migrate/v4` - Migrations
