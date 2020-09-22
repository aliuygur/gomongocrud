package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alioygur/gomongocrud"
	"github.com/alioygur/gomongocrud/cmd/restserver/api"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// prepare mongo client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mustGetEnv("MONGODB_DSN")))
	if err != nil {
		log.Fatalf("mongodb connection failed: %s", err)
	}

	tasksService := &gomongocrud.TasksService{
		IDs: &gomongocrud.UUID4{},
		Tasks: &gomongocrud.MongoRepo{
			Client: client,
		},
	}
	tasksHandler := &api.Handler{Tasks: tasksService}

	swagger, err := api.GetSwagger()
	if err != nil {
		log.Fatalf("unable to load swagger spec: %s", err)
	}
	swagger.Servers = nil

	e := echo.New()
	e.Use(echomiddleware.Logger())
	e.Use(middleware.OapiRequestValidator(swagger))

	api.RegisterHandlers(e, tasksHandler)

	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%s", getDefaultPort())))

}

func getDefaultPort() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}

func mustGetEnv(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	panic("env required: " + key)
}
