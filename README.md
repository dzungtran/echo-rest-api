[![Build Status](https://app.travis-ci.com/dzungtran/echo-rest-api.svg?branch=main)](https://app.travis-ci.com/dzungtran/echo-rest-api)
[![codecov](https://codecov.io/gh/dzungtran/echo-rest-api/branch/main/graph/badge.svg?token=hxaHIVyoBN)](https://codecov.io/gh/dzungtran/echo-rest-api)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/dzungtran/echo-rest-api/blob/master/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/dzungtran/echo-rest-api.svg)](https://pkg.go.dev/github.com/dzungtran/echo-rest-api)
[![GoReportCard Badge](https://goreportcard.com/badge/github.com/dzungtran/echo-rest-api)](https://goreportcard.com/report/github.com/dzungtran/echo-rest-api)

# Echo REST API boilerplate

A Golang restful API boilerplate based on Echo framework v4. Includes tools for module generation, db migration, authorization, authentication and more.
Any feedback and pull requests are welcome and highly appreciated. Feel free to open issues just for comments and discussions.

<!--toc-->
- [Echo REST API boilerplate](#echo-rest-api-boilerplate)
    * [Overview](#overview)
    * [Features](#features)
    * [Running the project](#running-the-project)
    * [Environment variables](#environment-variables)
    * [Commands](#commands)
    * [Folder structure](#folder-structure)
    * [Open source refs](#open-source-refs)
    * [Contributing](#contributing)
    * [TODOs](#todos)

<!-- tocstop -->

## Overview

![Request processing flow - Sequence Diagram](out/docs/diagrams/overview/request_flow.svg)

## Features

- [x] User Auth functionality (Signup, Login, Forgot Password, Reset Password, 2FA) using [Ory/Kratos](https://github.com/ory/kratos).
- [x] REST API using [labstack/echo](https://github.com/labstack/echo).
- [x] DB Migration using [golang-migrate/migrate](https://github.com/golang-migrate/migrate).
- [x] Configs via environmental variables.
- [x] Unit tests.
- [x] Dependency injection using [uber-go/dig](https://github.com/uber-go/dig).
- [x] Role based access control using [Open Policy Agent](https://github.com/open-policy-agent/opa).
- [x] Module generation, quickly create model, usecase, api handler.
- CLI support. try: `make tools` using [spf13/cobra](https://github.com/spf13/cobra).

## Running the project

- Make sure you have docker installed.
- Copy `.env.example` to `.env.docker`
- Add a new line `127.0.0.1	echo-rest.local` to `/etc/hosts` file.
- Run `docker compose up -d`.
- Go to `localhost:8088` to verify if the server works.

## Environment variables

By default, when you run application with `make run-api` command, it will look at `$HOME/.env` for exporting environment variabels.
Setting your config as Environment Variables is recommended as by 12-Factor App.

<details>
    <summary>Variables Defined in the project </summary>

| Name                   | Type    | Description                                                      | Example value                                 |
|------------------------|---------|------------------------------------------------------------------|-----------------------------------------------|
| DATABASE_URL           | string  | Data source URL for main DB                                      | postgres://world:hello@postgres/echo_rest_api |
| KRATOS_API_ENDPOINT    | string  | Public endpoint of Kratos                                        | http://kratos:4433/                           |
| KRATOS_WEBHOOK_API_KEY | string  | Api key for Kratos integration                                   | very-very-very-secure-api-key                 |
| PORT                   | integer | Http port (accepts also port number only for heroku compability) | 8088                                          |
| AUTO_MIGRATE           | boolean | Enable run migration every time the application starts           | true                                          |
| ENV                    | string  | Environment name                                                 | development                                   |
| REDIS_URL              | string  | Optional                                                         | redis://redis:6379                            |

</details>

## Commands

| Command                                  | Description                                                 |
|------------------------------------------|-------------------------------------------------------------|
| `make run-api`                           | Start REST API application                                  |
| `make build-api`                         | Build application binary                                    |
| `make setup`                             | Run commands to setup development env                       |
| `make run-db`                            | Run DB docker container on local                            |
| `go run ./tools/mod/ gen`                | Generate module component codes.                            |
| `make migration-create [migration_name]` | Create migration files. migration_name should be snake case |
| `make git-hooks`                         | Setup git hooks                                             |
| `make routes`                            | Generate routes file for authorization                      |

## Folder structure

```
.
├── 3rd-parties         # Thirdparty configs
├── cmd
│   └── api             # Main package of API service
├── config              # Application configs struct
├── delivery
│   ├── defines
│   ├── http
│   ├── requests
│   ...        
├── docs                # Content documentation and PlantUML for charts and diagrams
├── domains
├── infrastructure
├── migrations
│   └── sql             # Migration files
├── out                 # Output folder of PlantUML
├── pkg
│   ├── authz           # Contents Rego rule files for RBAC
│   ├── constants
│   ├── cue             # Contents cue files for data validation
│   ...
│   └── utils           # Contents helper functions
├── repositories
│   ├── postgres
│   └── redis
├── tests
├── tools
│   ├── modtool         # Module generation
│   ├── routes          # Generate routes file for Authorization
│   └── scripts         # Some helpful bash commands
└── usecases
```

## Open source refs
- https://www.ory.sh/docs/kratos/self-service
- https://cuelang.org/docs/about/
- https://www.openpolicyagent.org/docs/latest/
- https://echo.labstack.com/guide/


## Contributing

Please open issues if you want the template to add some features that is not in todos.

Create a PR with relevant information if you want to contribute in this template.

## TODOs

- [x] Update docker compose for ory/kratos.
- [x] Update README.md.
- [ ] Write more tests.
- [ ] Add support Heroku.
