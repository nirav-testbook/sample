package main

import (
	"context"
	"errors"
	"flag"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	kitlog "github.com/go-kit/kit/log"
	kitconsul "github.com/go-kit/kit/sd/consul"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	consul "github.com/hashicorp/consul/api"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	"sample/common/healthcheck"
	repo "sample/user/repo/mongo"
	"sample/user/user"
	"sample/user/user/pb"
)

const (
	svcName = "User"
)

var (
	host string
	port string
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	flag.StringVar(&host, "h", "127.0.0.1", "host")
	flag.StringVar(&port, "p", "8002", "port")
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
			TTL:                            "30s",
			DeregisterCriticalServiceAfter: "24h",
		},
	}, logger)
	reg.Register()
	go func() {
		healthcheck.InitConsulHealthCheck(consulClient.Agent(), logger, svcId, 10*time.Second)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer mongoClient.Disconnect(ctx)
	db := mongoClient.Database("test")

	userRepo, err := repo.NewUserRepo(db)
	if err != nil {
		panic(err)
	}

	userService := user.NewService(userRepo)
	userService = user.NewLogService(userService, logger)
	userHandler := user.NewGRPCHandler(userService)

	logger.Log("running GRPC server on", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
	pb.RegisterUserServer(baseServer, userHandler)

	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		sig := <-c
		errc <- errors.New("received signal " + sig.String())
	}()

	go func() {
		errc <- baseServer.Serve(listener)
	}()

	err = <-errc
	logger.Log("exit error", err)
	reg.Deregister()
	listener.Close()
}
