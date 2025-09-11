package app

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typefunco/dealer_dev_platform/internal/delivery"
	"github.com/typefunco/dealer_dev_platform/internal/repository"
	"github.com/typefunco/dealer_dev_platform/internal/service/auth"
	"github.com/typefunco/dealer_dev_platform/internal/service/performance"
	"github.com/typefunco/dealer_dev_platform/internal/utils/jwt"
	"log/slog"
	"os"
	"sync"
)

func RunApp() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	pool, err := pgxpool.New(context.Background(), "dsn")
	if err != nil {
		panic(err)
	}

	authRepo := repository.NewAuthRepository(pool, logger)
	performanceRepo := repository.NewPerformanceRepository(pool)
	JWTService := jwt.NewService()
	authService := auth.NewService(authRepo, JWTService, logger)
	perfService := performance.NewService(performanceRepo, logger)
	server := delivery.NewServer(authService, perfService, logger)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.RunServer()
	}()
	wg.Wait()
}
