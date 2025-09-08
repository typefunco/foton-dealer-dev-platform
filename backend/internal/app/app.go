package app

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typefunco/dealer_dev_platform/internal/delivery"
	"github.com/typefunco/dealer_dev_platform/internal/repository"
	"github.com/typefunco/dealer_dev_platform/internal/service/auth"
	"github.com/typefunco/dealer_dev_platform/internal/utils/jwt"
	"log/slog"
	"os"
	"sync"
)

func RunApp() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	pool, err := pgxpool.New(context.Background(), "example")
	if err != nil {
		//panic(err)
	}

	repo := repository.NewAuthRepository(pool, logger)
	JWTService := jwt.NewService()
	service := auth.NewService(repo, JWTService, logger)
	server := delivery.NewServer(service, logger)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.RunServer()
	}()
	wg.Wait()
}
