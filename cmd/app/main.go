package main

import (
	"context"
	"log"

	"github.com/Vidgimka/LoginTracking/internal/config"
	"github.com/Vidgimka/LoginTracking/internal/usecase"
)

func main() {

	cfg, err := config.NewConfig("", "")
	if err != nil {
		log.Fatalf("config initialization failed: %s", err)
	}
	ctx := context.Background()

	if err := usecase.Run(ctx, cfg); err != nil {
		log.Fatalf("application stert error: %v", err)
	}
}
