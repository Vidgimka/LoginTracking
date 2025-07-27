package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Vidgimka/LoginTracking.git/api"
	"gorm.io/gorm"
)

type ServiceInterface interface {
	RunTaskEverySecond(db *gorm.DB, ctx context.Context, stop <-chan struct{}, wg *sync.WaitGroup)
}

type service struct {
	client api.HttpClientInterface
}

func NewService(httpClient api.HttpClientInterface) ServiceInterface {
	return &service{
		client: httpClient,
	}
}

func (s *service) RunTaskEverySecond(db *gorm.DB, ctx context.Context, stop <-chan struct{}, wg *sync.WaitGroup) {
	ticker1 := time.NewTicker(time.Second)
	defer ticker1.Stop()
	defer wg.Done()
	for {
		select {
		case <-ticker1.C:
			fmt.Println("Running task every second")
			data, err := s.client.ReadDataFromAPI()
			if err != nil {
				log.Fatal("GET error:", err)
			}
			db.Create(&data)
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
