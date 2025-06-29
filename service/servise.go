package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Vidgimka/LoginTracking.git/api"
	"gorm.io/gorm"
)

type ServiceInterface interface {
	RunTaskEverySecond(db *gorm.DB, ctx context.Context, stop <-chan struct{})
}

type service struct {
	client api.HttpClientInterface
}

func NewService(httpClient api.HttpClientInterface) ServiceInterface {
	return &service{
		client: httpClient,
	}
}

func (s *service) RunTaskEverySecond(db *gorm.DB, ctx context.Context, stop <-chan struct{}) {
	// client := api.NewHttpClient()
	ticker1 := time.NewTicker(time.Second)
	defer ticker1.Stop()
	for {
		select {
		case <-ticker1.C:
			fmt.Println("Running task every second")
			data, err := s.client.ReadDataFromAPI()
			// client.ReadDataFromAPI()
			if err != nil {
				log.Fatal("GET error:", err)
			}
			db.Create(&data) // запись в БД
			log.Println("'Datetime' column added.")
			log.Println("Database entry complete")
		case <-stop:
			log.Println("no data received")
			return // выход из цикла
		case <-ctx.Done():
			log.Println("the user interrupted the program")
			return
		}
	}
}
