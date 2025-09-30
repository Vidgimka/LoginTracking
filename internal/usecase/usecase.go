package usecase

import (
	"context"
	"log"
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

func (s *service) Run(ctx context.Context) error {
	return nil
}
