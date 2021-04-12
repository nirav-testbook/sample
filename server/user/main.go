package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	repo "sample/user/repo/mongo"
	"sample/user/user"

	kitlog "github.com/go-kit/kit/log"
	kitconsul "github.com/go-kit/kit/sd/consul"
	consul "github.com/hashicorp/consul/api"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	consulClient, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		panic(err)
	}

	err = kitconsul.NewClient(consulClient).Register(&consul.AgentServiceRegistration{
		Name:    "User",
		Port:    8002,
		Address: "http://127.0.0.1",
	})
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer mongoClient.Disconnect(ctx)
	db := mongoClient.Database("test")

	logger := kitlog.NewJSONLogger(os.Stdout)

	userRepo, err := repo.NewUserRepo(db)
	if err != nil {
		log.Fatal(err)
	}

	userService := user.NewService(userRepo)
	userService = user.NewLogService(userService, kitlog.With(logger, "service", "User"))
	//userService = user.NewAuthService(userService, authService)
	userHandler := user.NewHandler(userService)

	r := http.NewServeMux()
	r.Handle("/user", userHandler)
	r.Handle("/user/", userHandler)

	log.Println("listening on", ":8002")
	err = http.ListenAndServe(":8002", r)
	if err != nil {
		log.Fatal(err)
	}
}
