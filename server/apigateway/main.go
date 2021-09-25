package main

import (
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	kitlog "github.com/go-kit/kit/log"
	kitconsul "github.com/go-kit/kit/sd/consul"
	consul "github.com/hashicorp/consul/api"

	"sample/auth/auth"
	authclient "sample/auth/auth/client"
	"sample/common/healthcheck"
	"sample/tb/lesson"
	lessonclient "sample/tb/lesson/client"
	"sample/tb/question"
	questionclient "sample/tb/question/client"
	"sample/user/user"
	userclient "sample/user/user/client"
)

const (
	svcName = "Apigateway"
)

var (
	host         string
	port         string
	retryMax     = 1
	retryTimeout = 500 * time.Millisecond
)

func main() {
	flag.StringVar(&host, "h", "127.0.0.1", "host")
	flag.StringVar(&port, "p", "8000", "port")
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

	userInstancer := kitconsul.NewInstancer(kitConsulClient, logger, "User", nil, true)
	//userSvc := userclient.NewWithLB(userInstancer, retryMax, retryTimeout, logger, http.DefaultClient)
	userSvc := userclient.NewGRPCWithLB(userInstancer, retryMax, retryTimeout, logger)
	userSvc = user.NewLogService(userSvc, logger)
	userHandler := user.NewHandler(userSvc)

	authInstancer := kitconsul.NewInstancer(kitConsulClient, logger, "Auth", nil, true)
	authSvc := authclient.NewWithLB(authInstancer, retryMax, retryTimeout, logger, http.DefaultClient)
	authHandler := auth.NewHandler(authSvc)

	tbInstancer := kitconsul.NewInstancer(kitConsulClient, logger, "TB", nil, true)
	lessonSvc := lessonclient.NewWithLB(tbInstancer, retryMax, retryTimeout, logger, http.DefaultClient)
	lessonSvc = lesson.NewLogService(lessonSvc, logger)
	lessonHandler := lesson.NewHandler(lessonSvc)
	questionSvc := questionclient.NewWithLB(tbInstancer, retryMax, retryTimeout, logger, http.DefaultClient)
	questionHandler := question.NewHandler(questionSvc)

	r := http.NewServeMux()
	r.Handle("/user", userHandler)
	r.Handle("/user/", userHandler)
	r.Handle("/auth", authHandler)
	r.Handle("/auth/", authHandler)
	r.Handle("/lesson", lessonHandler)
	r.Handle("/lesson/", lessonHandler)
	r.Handle("/question", questionHandler)
	r.Handle("/question/", questionHandler)

	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		sig := <-c
		errc <- errors.New("received signal " + sig.String())
	}()

	go func() {
		logger.Log("running HTTP server on", addr)
		errc <- http.ListenAndServe(":"+port, r)
	}()

	err = <-errc
	logger.Log("exit error", err)
	reg.Deregister()
}
