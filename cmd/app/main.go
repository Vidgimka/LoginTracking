package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Vidgimka/LoginTracking/internal/config"
	infrastructure "github.com/Vidgimka/LoginTracking/internal/infrastructure/database"
	"github.com/Vidgimka/LoginTracking/internal/infrastructure/repository"
	"github.com/Vidgimka/LoginTracking/internal/myhttp"
	"github.com/Vidgimka/LoginTracking/internal/myhttp/handlers"
	"github.com/Vidgimka/LoginTracking/internal/service"
)

func main() {
	var wg sync.WaitGroup
	config.LoadEnv()

	client := api.NewHttpClient()
	db, err := infrastructure.Init()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	repo := repository.NewPostgresGormRepo(db)
	service := service.NewService(client, repo)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	stop := make(chan struct{})
	wg.Add(1)
	go service.RunTaskEverySecond(ctx, stop, &wg)
	time.Sleep(1 * time.Second)
	close(stop)
	wg.Wait()

	handlers := handlers.NewHandlers(repo)
	router := myhttp.NewRouter(handlers)

	engine := router.SetRouter()
	engine.Run("localhost:8080")
}
