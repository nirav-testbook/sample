package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	kitlog "github.com/go-kit/kit/log"
	kitconsul "github.com/go-kit/kit/sd/consul"
	consul "github.com/hashicorp/consul/api"

	"sample/auth/auth"
	authclient "sample/auth/auth/client"
	"sample/tb/lesson"
	lessonclient "sample/tb/lesson/client"
	"sample/tb/question"
	questionclient "sample/tb/question/client"
	"sample/user/user"
	userclient "sample/user/user/client"
)

const (
	port = 8000
)

var (
	retryMax     = 1
	retryTimeout = 500 * time.Millisecond
)

func main() {
	consulClient, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		panic(err)
	}
	kitConsulClient := kitconsul.NewClient(consulClient)

	err = kitConsulClient.Register(&consul.AgentServiceRegistration{
		Name:    "Apigateway",
		Port:    port,
		Address: "http://127.0.0.1",
	})
	if err != nil {
		panic(err)
	}

	logger := kitlog.NewJSONLogger(os.Stdout)

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

	portStr := strconv.Itoa(port)
	log.Println("listening on", ":"+portStr)
	err = http.ListenAndServe(":"+portStr, r)
	if err != nil {
		log.Fatal(err)
	}
}
