# ![ZG Data Guard](logo.png)
ZG Data Guard is a centralized tool designed to streamline and secure the management of multiple databases across various environments. It simplifies administration by providing a unified platform to handle database ecosystems, technologies, instances, predefined roles, databases, users, and access control. All secured through JWT-protected API.
By using this tool organizations can effectively centralize database access management simplifying administration, improve operational efficiency, maintain compliance through detailed auditing and logging and promoting good security practices.

> A practical use case is to manage access to database instances (clusters) and their databases in development, staging, and production environments, ensuring that only authorized users (people or applications) have access to connect and perform specific operations on database objects based on their roles.

## Features

---
1. [**Ecosystem Management**](#ecosystem-management)
1. [**Database Technologies Management**](#database-technologies-management)
1. [**Database Instances (Clusters) Management**](#database-instances-clusters-management)
1. [**Predefined Roles**](#predefined-roles)
1. [**Databases Management**](#databases-management)
1. [**Database Users Management**](#database-users-management)
1. [**Access Control Management**](#access-control-management)
1. [**API Secured by JWT Tokens**](#api-secured-by-jwt-tokens)

#### Ecosystem Management

Manage ecosystems where database instances (clusters) are running, such as AWS, Cloud, or On-premises environments.

#### Database Technologies Management

Handle various database technologies, like Elasticsearch 6.2, PostgreSQL 13, PostgreSQL 16, etc.

#### Database Instances (Clusters) Management

Manage database instances (clusters) within specific ecosystems. Initially supports PostgreSQL instances, with future extensibility for other technologies.

- **Operations:** Create, Read, Update
- **Additional Functions:**
  - **Test Connection:** Verify connectivity to the database instance.
  - **Synchronize Databases:** Update the list of databases within the instance.
  - **Create Predefined Roles:** Set up predefined roles in the instance context.
  - **Enable/Disable Instance:** Remove all defined accesses from all users when disabling; also disables all databases within the cluster.

#### Predefined Roles

Utilize predefined roles assigned to users in specific databases to enforce the principle of least privilege. The roles are defined as follows:

- **Roles:**
  - **User Read Only:** Read-only permissions on all schemas of the database.
  - **Developer:** DML permissions (SELECT, INSERT, UPDATE, DELETE) and usage of sequences, functions, and types in all schemas.
  - **DevOps:** DML and DDL permissions (CREATE, ALTER, TRUNCATE, DROP) on tables, functions, sequences, triggers, types, etc., in all schemas.
  - **Application:** Same as DevOps, intended for application users.
- **Notes:**
  - No role can grant or revoke privileges to itself or other roles.
  - No role has SUPERUSER permission.
  - Roles are designed following the **Principle of Least Privilege**.

#### Databases Management

Manage existing databases within instances and apply predefined roles to establish permissions on database objects like schemas, tables, functions, views, sequences, and types.

#### Database Users Management

Manage users who can be assigned to database instances or databases with specific roles (e.g., `foo.bar`, `john.doe`). It can be a user for a person or an application.

- **Operations:** Create, Read, Update, Enable/Disable Users

#### Access Control Management

Control users' access to instances/databases by granting or revoking connect permission, with comprehensive logging for auditing purposes.

- **Operations:**
  - **Grant Access:** Provide users access to one or more instances.
  - **Revoke Access:** Remove users' access from instances.
  - **Logging:** Record and display the results of binding and unbinding operations.

#### API Secured by JWT Tokens

The API is protected using JWT (JSON Web Tokens) for secure authentication and authorization, ensuring safe communication between clients and the server.

## Technologies Used

---

- GoLang 1.22+
- PostgreSQL 16+
- [Keycloak 26+](https://www.keycloak.org/) for OAuth2 and JWT
- AES-256 encryption for sensitive data protection
- Swagger for API documentation
- Makefile for task automation

 The following dependencies are used in this project (generated using [Glice](https://github.com/ribice/glice)):

```bash
+--------------------------------------+-------------------------------------------+--------------+
|              DEPENDENCY              |                  REPOURL                  |   LICENSE    |
+--------------------------------------+-------------------------------------------+--------------+
| github.com/go-chi/chi/v5             | https://github.com/go-chi/chi             | MIT          |
| github.com/go-chi/jwtauth            | https://github.com/go-chi/jwtauth         | MIT          |
| github.com/golang-migrate/migrate/v4 | https://github.com/golang-migrate/migrate | Other        |
| github.com/google/uuid               | https://github.com/google/uuid            | bsd-3-clause |
| github.com/joho/godotenv             | https://github.com/joho/godotenv          | MIT          |
| github.com/lib/pq                    | https://github.com/lib/pq                 | MIT          |
| github.com/stretchr/testify          | https://github.com/stretchr/testify       | MIT          |
| github.com/swaggo/http-swagger       | https://github.com/swaggo/http-swagger    | MIT          |
| github.com/swaggo/swag               | https://github.com/swaggo/swag            | MIT          |
| golang.org/x/oauth2                  | https://go.googlesource.com/oauth2        |              |
+--------------------------------------+-------------------------------------------+--------------+
```

1. Go-chi - HTTP Middleware Router
1. JWT Auth - JWT Authentication
1. Golang Migrate - Database Migrations
1. Google UUID - UUID generator
1. Godotenv - Environment variables
1. lib/pq - PostgreSQL driver
1. Testify/Assert - Asserting test results
1. Swaggo - Swagger documentation
1. OAuth2 - OAuth2 library

## Usage

---

### 1. Installation
- [**Docker**](#docker)
- [**Docker Compose**](#docker-compose)

### 2. Setup Project

#### Environment Variables
1. Configure the environment variables by creating a `.env` file in the root directory. Use the `.env.example` file as a template.
```sh
   cp .env.example .env
```
2. Update the `.env` file envs according to your preferences.

#### (TODO) Keycloak to Secure the API

1. Visit http://localhost:8080.
2. Log in with the credentials defined for Keycloak in the `.env` file.
3. Create a new realm, e.g., `zg-data-guard`.
4. Create a new client, e.g., `zg-data-guard-api`.
5. Configure the client with the following settings:
    - **Access Type:** Confidential
    - **Valid Redirect URIs:** `http://localhost:8081/*`
    - **Web Origins:** `http://localhost:8081`
    - **Client Protocol:** `openid-connect`
    - **Service Accounts Enabled:** On
    - **Authorization Enabled:** On
    - **Direct Access Grants Enabled:** On
    - **Standard Flow Enabled:** On
6. Create a new user and assign the user to the client.
7. Update the `.env` file with the Keycloak settings.
8. Restart the API server.
9. Access the API at [http://localhost:8081](http://localhost:8081).
10. Authenticate using the Keycloak credentials.
11. Access the protected endpoints. Use the Swagger documentation to test the API endpoints.

### 3. Running the API
- `docker compose up`: Run the services defined in the `docker-compose.yml` file.
- After the API is running:
  - Home page - [http://localhost:8081](http://localhost:8081)
  - Click in Login to access the Keycloak login page and authenticate.
  - Get the JWT token and use it to access the API endpoints.
  - You can use the Swagger UI to interact with the endpoints. The API can be accessed at [http://localhost:8081/docs/index.html](http://localhost:8081/docs/index.html).
  - Click on the `Authorize` button and enter the JWT token in the `Value` field with the `Bearer` prefix.

## Development Guide

---

### Installation of Tools

To set up the development environment, you need to install the following tools:

- [**GoLang 1.22+**](#golang-122)
- [**Docker**](#docker)
- [**Docker Compose**](#docker-compose)
- [**GolangCI-Lint**](#golangci-lint)
- [**Migrate**](#migrate)
- [**Swaggo**](#swaggo)
- [**Make**](#make)

#### GoLang 1.22+

Download and install GoLang from the [official website](https://go.dev/dl/). Follow the instructions for your operating system. After installation, verify the installation by running:

```bash
go version
```

#### Docker

Download and install Docker from the [official Docker website](https://www.docker.com/get-started). Follow the installation guide for your operating system. Verify the installation with:

```bash
docker --version
```

#### Docker Compose

Docker Compose is included with Docker Desktop for Windows and macOS. For Linux, install it separately by following the [official instructions](https://docs.docker.com/compose/install/). Verify the installation:

```bash
docker-compose --version
```

#### GolangCI-Lint

Install [GolangCI-Lint](https://golangci-lint.run/) for linting Go code:

```bash
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest
```

Alternatively, you can use Homebrew on macOS:

```bash
brew install golangci-lint
```

Verify the installation:

```bash
golangci-lint --version
```

Usage:

```bash
make lint
```

#### Migrate

Install [Migrate](https://github.com/golang-migrate/migrate/) for database migrations:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Ensure that your `GOPATH/bin` is in your `PATH` environment variable. Verify the installation:

```bash
migrate --version
```

To create a new migration, run:

```bash
make create_migration
```

#### Swaggo

Install [Swaggo](https://github.com/swaggo/swag) to generate Swagger documentation:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Verify the installation:

```bash
swag --version
```

#### Make

Ensure that `make` is installed on your system to use the provided `Makefile` for task automation.

- **On macOS** (if not already installed):

  ```bash
  xcode-select --install
  ```

- **On Linux** (using apt):

  ```bash
  sudo apt-get install build-essential
  ```

- **On Windows**:

  Install [Make for Windows](http://gnuwin32.sourceforge.net/packages/make.htm) or use a Unix-like environment like Git Bash.

Verify the installation:

```bash
make --version
```

After installing all the tools, you should be ready to set up and run the project.

### Directory Structure

---
```
.
├── cmd
│   ├── zg-data-guard
│       └── main.go     //main function start the server
├── config              //configurations for the project
├── docs                //swagger API documentation
├── internal
│   ├── database        //connector, migrations, sql files and storages
│   ├── dto             //data transfer objects
│   ├── entity          //database entities, models
│   ├── usecase         //business logic
│   ├── webserver       //http server, routes, handlers, middlewares
├── pkg                 //shared packages, utilities, security functions like crypto and jwt
└── testdata            //test data for unit tests, mocks
...
```

### Set up the Project

---
1. Fork this repo and `git clone` it to your local machine.
1. Configure the environment variables by creating a `.env` file in the root directory. Use the `.env.example` file as a template.
```sh
   cp .env.example .env
```
1. Update the `.env` file with your database connection details, [Keycloak settings](#setup-keycloak-to-secure-the-api), and other configurations.
1. To install dependencies, run:
```bash
make install
```

### Build

---
1. To build the project and generate the executable, run:

```bash
make build
```

2. To clean the project, run:

```bash
make clean
```

### Running API

--- 
* To run the app, execute:

```bash
make run
```

* To generate the API Swagger documentation and execute, run:

```bash
make run-with-docs
```

* If you want to generate docs without running the server, run:

```bash
make docs
```

* By default, the server will be available at: [http://localhost:8081](http://localhost:8081)
* The Swagger documentation will be available at: [http://localhost:8081/docs/index.html](http://localhost:8081/docs/index.html)
* Health check endpoint: [http://localhost:8081/healthcheck/info](http://localhost:8081/healthcheck/info)

### About Tests

--- 

To run the project tests, run in the terminal:

```bash
make test
```

To run the tests with more details, run:

```bash
make test-verbose
```

To run the project tests without cache and return the total number of tests executed:

```bash
make test-count
```

To run the tests and generate an HTML file with a complete coverage report for each file, run:

```bash
make test-cover-report
```

To validate the total test coverage percentage of all files, run:

```bash
make coverage
```

To validate the total test coverage percentage of business logic files, run:

```bash
make core-coverage
```

The minimum test coverage percentage is configured in the `Makefile` file in the `MIN_COVERAGE` and `MIN_CORE_COVERAGE` variables.

### Release

--- 

The `make release=<version>` command was created to be used by CI/CD pipeline.

## Credits

--- 

This project was created by the [ZG Soluções](https://www.linkedin.com/company/zg-solucoes/) team.

Enjoy!

![ZG Soluções](logo-zg.png)

[https://zgsolucoes.com.br/](https://zgsolucoes.com.br/)
