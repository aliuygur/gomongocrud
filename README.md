# Sample Golang CRUD App

This is sample Golang application that using mongodb for a challenge.

## Build & Run (Source)

```
# build rest server
go build -o restserver ./cmd/restserver
# build graphql server
go build -o graphqlserver ./cmd/graphqlserver

# run restserver
MONGO_DSN=mongodb://localhost:27017 go run restserver
# run graphqlserver
MONGO_DSN=mongodb://localhost:27017 go run graphqlserver
```

## Build & Run (Docker)

```
docker build -t gomongocrud .

# run restserver
docker run --env MONGODB_DSN=mongodb://localhost:27017 gomongocrud restserver
# run graphqlserver
docker run --env MONGODB_DSN=mongodb://localhost:27017 gomongocrud graphqlserver
```

## Directory Tree

```
├── README.md
├── cmd
│   ├── gqlserver              <- Graphql adapter
│   └── restserver             <- Restful adapter
│       ├── api
│       │   ├── api.go
│       │   ├── oas3.yaml
│       │   ├── server.gen.go
│       │   └── types.gen.go
│       └── main.go
├── error.go
├── go.mod
├── go.sum
├── models.go                  <- domain layer
├── mongo_repo.go
├── tasks_service.go           <- use case layer
└── uuid4_generator.go
```
