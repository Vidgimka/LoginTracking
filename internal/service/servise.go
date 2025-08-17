package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Vidgimka/LoginTracking/internal/domain"
	"github.com/Vidgimka/LoginTracking/internal/infrastructure/client"
)

type svtpClient interface {
	GetUsersOnline(ctx context.Context) ([]client.Data, error)
}

type userRepositpry interface {
	// Create(data interface{}) error
	GetLines(ctx context.Context, login string, start, end time.Time) ([]domain.LineData, error)
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

func (s *service) RunTaskEverySecond(ctx context.Context, stop <-chan struct{}, wg *sync.WaitGroup) {
	ticker1 := time.NewTicker(time.Second)
	defer ticker1.Stop()
	defer wg.Done()
	for {
		select {
		case <-ticker1.C:
			fmt.Println("Running task every second")
			data, err := s.client.GetUsersOnline(ctx)
			if err != nil {
				log.Fatal("GET error:", err)
			}
			s.repo.Create(&data)
			log.Println("'datetime' column added.")
			log.Println("database entry complete")
		case <-stop:
			log.Println("closed by stop channel")
			return
		case <-ctx.Done():
			log.Println("the user cancelled the request")
			return
		}
	}
}
