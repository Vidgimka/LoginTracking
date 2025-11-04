package usecase

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Vidgimka/LoginTracking/internal/config"
	"github.com/Vidgimka/LoginTracking/internal/controller/cron"
	v1 "github.com/Vidgimka/LoginTracking/internal/controller/http/v1"
	"github.com/Vidgimka/LoginTracking/internal/domain"
	"github.com/Vidgimka/LoginTracking/internal/infrastructure/client"
	"github.com/Vidgimka/LoginTracking/internal/infrastructure/database"
	"github.com/Vidgimka/LoginTracking/internal/infrastructure/repository"
)

type svtpClient interface {
	GetUsersOnline(ctx context.Context) ([]domain.Data, error)
}

type userRepositpry interface {
	CreateData(ctx context.Context, usersOnline []domain.Data) error
	GetPoints(ctx context.Context, login string, start, end time.Time) ([]domain.PointData, error)
}

type service struct {
	client svtpClient
	repo   userRepositpry
}

func NewService(httpClient svtpClient, db userRepositpry) *service {
	return &service{
		client: httpClient,
		repo:   db,
	}
}

func (s *service) SaveCurrentUsersLocation(ctx context.Context) error {
	data, err := s.client.GetUsersOnline(ctx)
	if err != nil {
		log.Fatal("GET error:", err)
	}
	if err := s.repo.CreateData(ctx, data); err != nil {
		log.Fatal("repo.CreateData:", err)
	}
	log.Println("database entry complete")
	return nil
}

func (s *service) BuildPointsByDate(ctx context.Context, login string, start, end time.Time) ([]domain.PointData, error) {
	pointData, err := s.repo.GetPoints(ctx, login, start, end)
	if err != nil {
		return nil, fmt.Errorf("repo.GetPoints: %w", err)
	}
	return pointData, nil
}

func Run(ctx context.Context, cfg *config.Config) error {
	ctxWithSignal, cancelWithSignal := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancelWithSignal()

	httpClient := &http.Client{
		Timeout: time.Duration(cfg.Client.ClientTimeOut),
	}
	svtpHttpClient := client.New(httpClient, cfg.Client.BaseUrl)

	pool, err := database.SetPool(ctxWithSignal, &cfg.DataBase)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	repo := repository.NewPostgresPgxRepo(pool)

	service := NewService(svtpHttpClient, repo)

	stop := make(chan struct{})
	go cron.RunTaskEverySecond(ctxWithSignal, service, stop)
	time.Sleep(1 * time.Second)
	close(stop)

	handlers := v1.NewHandlers(service)

	router := v1.NewRouter(handlers)

	engien := router.SetRouter()
	engien.Run("localhost:8080")

	return nil
}
