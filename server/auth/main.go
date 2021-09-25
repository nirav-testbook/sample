package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	kitlog "github.com/go-kit/kit/log"
	kitconsul "github.com/go-kit/kit/sd/consul"
	consul "github.com/hashicorp/consul/api"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"sample/auth/auth"
	repo "sample/auth/repo/mongo"
	"sample/common/healthcheck"
	userclient "sample/user/user/client"
)

const (
	svcName = "Auth"
)

var (
	host         string
	port         string
	retryMax     = 1
	retryTimeout = 500 * time.Millisecond
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	flag.StringVar(&host, "h", "127.0.0.1", "host")
	flag.StringVar(&port, "p", "8001", "port")
	flag.Parse()
	addr := host + ":" + port
	svcId := svcName + ":" + addr
	portNum, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}

	logger := kitlog.NewJSONLogger(os.Stdout)
	logger = kitlog.With(logger, "svc", svcName, "svcId", svcId)

	consulClient, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		panic(err)
	}
	kitConsulClient := kitconsul.NewClient(consulClient)
	reg := kitconsul.NewRegistrar(kitConsulClient, &consul.AgentServiceRegistration{
		ID:      svcId,
		Name:    svcName,
		Address: host,
		Port:    portNum,
		Check: &consul.AgentServiceCheck{
			CheckID:                        svcId,
			TTL:                            "5s",
			DeregisterCriticalServiceAfter: "24h",
		},
	}, logger)
	reg.Register()
	go func() {
		healthcheck.InitConsulHealthCheck(consulClient.Agent(), logger, svcId, time.Second*3)
	}()

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

	userInstancer := kitconsul.NewInstancer(kitConsulClient, logger, "User", nil, true)
	userService := userclient.NewWithLB(userInstancer, retryMax, retryTimeout, logger, http.DefaultClient)

	authService := auth.NewService(sessionRepo, userService)
	authService = auth.NewLogService(authService, kitlog.With(logger, "service", "Auth"))
	authHandler := auth.NewHandler(authService)

	r := http.NewServeMux()
	r.Handle("/auth", authHandler)
	r.Handle("/auth/", authHandler)

	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("running HTTP server on", addr)
		errc <- http.ListenAndServe(":"+port, r)
	}()

	err = <-errc
	logger.Log("exit error", err)
	reg.Deregister()
}
