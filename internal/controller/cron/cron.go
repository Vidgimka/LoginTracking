package cron

import (
	"context"
	"fmt"
	"log"
	"time"
)

type service interface {
	SaveCurrentUsersLocation(context.Context) error
}

func RunTaskEverySecond(ctx context.Context, service service, stop <-chan struct{}) {
	ticker1 := time.NewTicker(time.Second)
	defer ticker1.Stop()
	for {
		select {
		case <-ticker1.C:
			fmt.Println("Running task every second")
			if err := service.SaveCurrentUsersLocation(ctx); err != nil {
				log.Printf("service.SaveCurrentUsersLocation: %v", err)
			}
		case <-stop:
			log.Println("closed by stop channel")
			return
		case <-ctx.Done():
			log.Println("the user cancelled the request")
			return
		}
	}
}
