package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"sample/auth/auth"
	repo "sample/auth/repo/mongo"
	userclient "sample/user/user/client"

	kitlog "github.com/go-kit/kit/log"
	kitconsul "github.com/go-kit/kit/sd/consul"
	consul "github.com/hashicorp/consul/api"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	retryMax     = 3
	retryTimeout = 500 * time.Millisecond
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	logger := kitlog.NewJSONLogger(os.Stdout)

	consulClient, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		panic(err)
	}

	consul := kitconsul.NewClient(consulClient)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer mongoClient.Disconnect(ctx)
	db := mongoClient.Database("test")

	sessionRepo, err := repo.NewSessionRepo(db)
	if err != nil {
		log.Fatal(err)
	}

	userInstancer := kitconsul.NewInstancer(consul, logger, "User", nil, true)
	userService := userclient.NewWithLB(userInstancer, retryMax, retryTimeout, logger, http.DefaultClient)

	authService := auth.NewService(sessionRepo, userService)
	authService = auth.NewLogService(authService, kitlog.With(logger, "service", "Auth"))
	authHandler := auth.NewHandler(authService)

	r := http.NewServeMux()
	r.Handle("/auth", authHandler)
	r.Handle("/auth/", authHandler)

	log.Println("listening on", ":8001")
	err = http.ListenAndServe(":8001", r)
	if err != nil {
		log.Fatal(err)
	}
}