package main

import (
	"context"
	"flag"
	"fmt"
	"log"
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

	repo "sample/user/repo/mongo"
	"sample/user/user"
	"sample/user/user/pb"
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

	consulClient, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		panic(err)
	}
	portNum, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
	svc := &consul.AgentServiceRegistration{
		ID:      "User" + " - " + addr,
		Name:    "User",
		Address: host,
		Port:    portNum,
	}
	kitConsulClient := kitconsul.NewClient(consulClient)
	err = kitConsulClient.Register(svc)
	if err != nil {
		panic(err)
	}

	logger := kitlog.NewJSONLogger(os.Stdout)

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
		log.Fatal(err)
	}

	userService := user.NewService(userRepo)
	userService = user.NewLogService(userService, kitlog.With(logger, "service", "User"))
	//userService = user.NewAuthService(userService, authService)
	userHandler := user.NewGRPCHandler(userService)

	log.Println("listening GRPC on", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
	pb.RegisterUserServer(baseServer, userHandler)

	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		errc <- baseServer.Serve(listener)
	}()

	logger.Log("exit", <-errc)
	err = kitConsulClient.Deregister(svc)
	if err != nil {
		logger.Log(err)
	}
	listener.Close()
}
