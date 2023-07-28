# Book Management API

This repo is a RESTful API for managing a library of books. Developed using Go, it uses the Gorilla Mux router for HTTP request routing and Gorm ORM for PostgreSQL. The API features endpoints for listing, adding, updating, and deleting books and collections from the database.

## Prerequisites

- Docker
- Docker Compose
- Go 1.16 or newer

## Getting Started

The project uses Docker Compose to manage its environment. Before running the project, you need to clone it to your local machine.

```
git clone https://github.com/gracesusanto/book-rest-api.git
cd book-rest-api
```

Then, you can start the environment using Docker Compose:

```
docker-compose up -d
```

This command pulls the necessary Docker images, creates a new Docker network, and starts two services:

- `postgresql` - A PostgreSQL 14 server, exposed at `localhost:5432`.

The PostgreSQL server has the following environment variables set:

- `POSTGRES_DB` - The name of the default database (`app`).
- `POSTGRES_USER` - The username of the default user (`app`).
- `POSTGRES_PASSWORD` - The password of the default user (`password`).

The server also mounts two volumes:

- `postgres` - The data of the PostgreSQL server.
- `./database` - The SQL scripts to initialize the database.

After you've started the environment, you can run the tests with the following command:

```
go test -v ./...
```

## Database Tables

![](tables.png?raw=true)

## Run Server

```bash
go run main.go
```

## Command Line Interface App

```bash
cd cmd
go build -o bookcli
```

## Unit Test

```bash
go test -v
```

## Documentation
The documentation is available at: http://github.com/gracesusanto/book-rest-api/tree/main/doc