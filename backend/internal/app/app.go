package app

import (
	"github.com/typefunco/dealer_dev_platform/internal/delivery"
	"log/slog"
	"os"
	"sync"
)

func RunApp() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	server := delivery.NewServer(logger)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.RunServer()
	}()
	wg.Wait()
}
