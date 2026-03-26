# Echo REST API Boilerplate

A Golang RESTful API boilerplate based on Echo framework v4. Includes tools for module generation, DB migration, authorization, authentication, and more.

Any feedback and pull requests are welcome and highly appreciated. Feel free to open issues for comments and discussions.

<!--toc-->
- [Echo REST API Boilerplate](#echo-rest-api-boilerplate)
  - [How to Use This Template](#how-to-use-this-template)
  - [Overview](#overview)
  - [Features](#features)
  - [Development](#development)
    - [Prerequisites](#prerequisites)
    - [Setup](#setup)
    - [Development Workflow](#development-workflow)
    - [Code Generation](#code-generation)
  - [Running the Project](#running-the-project)
    - [Using Docker (Recommended)](#using-docker-recommended)
    - [Local Development](#local-development)
  - [Swagger Docs](#swagger-docs)
  - [Environment Variables](#environment-variables)
  - [Commands](#commands)
  - [Folder Structure](#folder-structure)
  - [Open Source Refs](#open-source-refs)
  - [Contributing](#contributing)
  - [TODOs](#todos)
<!-- tocstop -->

## How to Use This Template

> **DO NOT FORK** - This template is meant to be used via the **[Use this template](https://github.com/dzungtran/echo-rest-api/generate)** feature.

1. Click on **[Use this template](https://github.com/dzungtran/echo-rest-api/generate)**
2. Give your project a name (recommend all lowercase with underscores, e.g., `my_awesome_project`)
3. Wait until the first CI run finishes (GitHub Actions will process the template and commit to your new repo)
4. Clone your new project and start coding!

> **NOTE**: Wait until the first CI run on GitHub Actions completes before cloning your new project.

## Overview

This is a production-ready Go REST API boilerplate using Echo v4 with:
- Firebase Authentication for JWT-based auth
- Open Policy Agent (OPA) for RBAC
- PostgreSQL with sqlx
- Uber-go/dig for dependency injection
- Modular architecture for easy feature addition

## Features

- [x] User authentication (Signup, Login, Forgot Password, Reset Password, 2FA) using **Firebase Auth**
- [x] REST API using [labstack/echo](https://github.com/labstack/echo)
- [x] DB migration using [golang-migrate/migrate](https://github.com/golang-migrate/migrate)
- [x] Modular structure
- [x] Configuration via environment variables
- [x] Unit tests
- [x] Dependency injection using [uber-go/dig](https://github.com/uber-go/dig)
- [x] Role-based access control using [Open Policy Agent](https://github.com/open-policy-agent/opa)
- [x] Module generation - quickly create models, usecases, and API handlers
- [x] CLI support via [spf13/cobra](https://github.com/spf13/cobra)
- [x] API docs generation using [swaggo](https://github.com/swaggo/swag)

## Development

### Prerequisites

- Go 1.21+
- Docker and Docker Compose
- [swag](https://github.com/swaggo/swag) CLI (`go install github.com/swaggo/swag/cmd/swag@v1.16.3`)
- [migrate](https://github.com/golang-migrate/migrate) CLI

### Setup

```bash
# Install dependencies
make setup

# Copy environment file
cp .env.example .env

# Start database
make run-db

# Run migrations
make migration-up

# Generate API docs
make docs
```

### Development Workflow

```bash
# Run the API
make run-api

# Run tests
make test

# Build binary
make build-api

# Create a new migration
make migration-create name=add_new_table

# Generate a new module
go run ./tools/mod/ gen -n ModuleName

# After generating a module, update routes
make routes
```

### Code Generation

Generate a new module with:
```bash
go run ./tools/mod/ gen -n ModuleName
```

This creates the domain, repository, usecase, and handler files. Run `make routes` after to update authorization.

## Running the Project

### Using Docker (Recommended)

```bash
# Copy and edit environment file
cp .env.example .env.docker

# Start all services
docker compose up -d

# Verify API is running
curl http://localhost:8088
```

### Local Development

```bash
# Start database
make run-db

# Run API
make run-api
```

API is available at `http://localhost:8088`. Swagger docs at `http://localhost:8088/docs/index.html`.

## Swagger Docs

Generate API documentation:

```bash
make docs
```

After generation, run the app and open `http://localhost:8088/docs/index.html`.

## Environment Variables

By default, when running with `make run-api`, the application looks for `$HOME/.env`. Environment variables are recommended per the 12-Factor App methodology.

<details>
    <summary>Variables defined in the project</summary>

| Name                        | Type   | Description                                            | Example value                                |
|-----------------------------|--------|--------------------------------------------------------|---------------------------------------------|
| DATABASE_URL                | string | PostgreSQL connection string                           | postgres://world:hello@postgres/echo_rest_api |
| PORT                        | int    | HTTP port (also accepts port number for Heroku)       | 8088                                        |
| AUTO_MIGRATE               | bool   | Enable migration on application startup               | true                                        |
| ENV                         | string | Environment name                                      | development                                 |
| AUTH_PROVIDER              | string | Authentication provider                               | firebase                                    |
| FIREBASE_CREDENTIALS       | JSON   | Firebase admin key                                    | {firebase_admin_key}                        |
| FIREBASE_AUTH_CREDENTIALS  | JSON   | Firebase auth key                                     | {firebase_auth_key}                        |
</details>

## Commands

| Command                                  | Description                                              |
|------------------------------------------|----------------------------------------------------------|
| `make run-api`                           | Start the REST API                                      |
| `make build-api`                         | Build the application binary to `./bin/`                |
| `make test`                              | Run tests with coverage                                 |
| `make setup`                             | Install development dependencies (swag, migrate)        |
| `make run-db`                            | Start PostgreSQL via Docker Compose                     |
| `make migration-up`                      | Run database migrations                                 |
| `make migration-down`                    | Rollback the last migration                             |
| `make migration-create name=<name>`      | Create new migration files (use snake_case for name)   |
| `make docs`                              | Generate Swagger API documentation                       |
| `make routes`                            | Generate routes file for OPA authorization               |
| `make git-hooks`                         | Install pre-commit hooks                                |
| `go run ./tools/mod/ gen -n ModuleName` | Generate a new module (e.g., Booking)                  |

## Folder Structure

```
.
├── cmd
│   └── api                 # Main application entry point
├── config                  # Application configuration
├── docs                    # Swagger documentation
├── infrastructure          # Database and external services
├── migrations
│   └── sql                 # SQL migration files
├── modules
│   ├── core                # Core module (users, orgs, auth)
│   ├── projects            # Demo module
│   └── shared              # Shared code across modules
├── pkg
│   ├── authz               # OPA Rego rules for RBAC
│   ├── constants           # Error definitions, context keys
│   ├── contexts            # Request context utilities
│   ├── cue                 # CUE validation (deprecated)
│   ├── hook                # Event hooks system
│   ├── logger              # Zap logger wrapper
│   ├── middlewares         # Auth and RBAC middleware
│   ├── sql-tools           # SQL builder, transaction helpers
│   ├── utils               # Helper functions
│   └── wrapper             # HTTP response wrapper
├── tests                   # Test utilities
└── tools
    ├── mod                 # Module generation CLI
    └── routes              # Routes generation for authorization
```

## Open Source Refs

- [CUE Language](https://cuelang.org/docs/about/)
- [Open Policy Agent](https://www.openpolicyagent.org/docs/latest/)
- [Echo Framework](https://echo.labstack.com/guide/)
- [Firebase Auth Admin](https://firebase.google.com/docs/auth/admin/)
- [Firebase Go SDK](https://pkg.go.dev/firebase.google.com/go/auth)

## Contributing

Please open issues for features not in the TODOs. Create a PR with relevant information if you'd like to contribute to this template.

## TODOs

- [x] Update docker compose for ory/kratos
- [x] Update README.md
- [x] Update API docs
- [ ] Write more tests
- [ ] Replace CUE validation with go-playground/validator
