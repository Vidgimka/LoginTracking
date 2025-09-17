package application

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Vidgimka/LoginTracking/internal/domain"
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

func (s *service) SaveCurrentUsersLocation(ctx context.Context) {
	data, err := s.client.GetUsersOnline(ctx)
	if err != nil {
		log.Fatal("GET error:", err)
	}
	if err := s.repo.CreateData(ctx, data); err != nil {
		log.Fatal("repo.CreateData:", err)
	}
	log.Println("database entry complete")
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
			s.repo.CreateData(ctx, data)
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
