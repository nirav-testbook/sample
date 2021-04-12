package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"sample/tb/lesson"
	"sample/tb/question"
	repo "sample/tb/repo/mongo"

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
		Name:    "TB",
		Port:    8003,
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

	questionRepo, err := repo.NewQuestionRepo(db)
	if err != nil {
		log.Fatal(err)
	}

	lessonRepo, err := repo.NewLessonRepo(db)
	if err != nil {
		log.Fatal(err)
	}

	lessonService := lesson.NewService(lessonRepo, questionRepo)
	lessonService = lesson.NewLogService(lessonService, kitlog.With(logger, "service", "Lesson"))
	lessonHandler := lesson.NewHandler(lessonService)

	questionService := question.NewService(questionRepo)
	questionService = question.NewLogService(questionService, kitlog.With(logger, "service", "Lesson"))
	questionHandler := question.NewHandler(questionService)

	r := http.NewServeMux()

	r.Handle("/lesson", lessonHandler)
	r.Handle("/lesson/", lessonHandler)
	r.Handle("/question", questionHandler)
	r.Handle("/question/", questionHandler)

	log.Println("listening on", ":8003")
	err = http.ListenAndServe(":8003", r)
	if err != nil {
		log.Fatal(err)
	}
}
