package main

import (
	_ "insider-assignment/docs"
	"insider-assignment/internal/cache"
	"insider-assignment/internal/config"
	"insider-assignment/internal/handler"
	"insider-assignment/internal/repository"
	"insider-assignment/internal/service"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Insider Messaging API
// @version 1.0
// @description This API manages message sending and control for the Insider assignment project.
// @host localhost:8080
// @BasePath /api/v1
func main() {

	cfg := config.Load()
	repo, err := repository.New(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	redisCache, err := cache.NewCache(cfg.RedisHost, cfg.RedisPort)
	if err != nil {
		log.Println("redis not available", err)
	} else {
		defer redisCache.Close()

	}

	svc := service.NewSenderService(repo, redisCache, cfg.WebhookURL, cfg.AuthKey, cfg.Interval, cfg.BatchSize)
	svc.Start()

	h := handler.NewHandler(svc)
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/health", h.Health).Methods("GET")
	r.HandleFunc("/api/v1/control/{action}", h.Control).Methods("POST")
	r.HandleFunc("/api/v1/sent-messages", h.GetSent).Methods("GET")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Printf("Starting server on port %s", cfg.Port)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}

}
