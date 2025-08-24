package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Vidgimka/LoginTracking/internal/config"
	"github.com/Vidgimka/LoginTracking/internal/infrastructure/client"
	"github.com/Vidgimka/LoginTracking/internal/infrastructure/database"
	"github.com/Vidgimka/LoginTracking/internal/infrastructure/repository"
	"github.com/Vidgimka/LoginTracking/internal/myhttp"
	"github.com/Vidgimka/LoginTracking/internal/myhttp/handlers"
	"github.com/Vidgimka/LoginTracking/internal/service"
)

func main() {
	var wg sync.WaitGroup

	cfg, err := config.NewConfig("", "")
	if err != nil {
		log.Fatalf("config initialization failed: %s", err)
	}

	httpClient := &http.Client{
		Timeout: time.Duration(cfg.Client.ClientTimeOut),
	}

	svtpHttpClient := client.New(httpClient, cfg.Client.BaseUrl)

	db, err := database.NewDatabase(&cfg.DataBase)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	repo := repository.NewPostgresGormRepo(db)
	service := service.NewService(svtpHttpClient, repo)

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
