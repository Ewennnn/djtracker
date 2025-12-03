package main

import (
	"djtracker/internal/api"
	"djtracker/internal/config"
	"djtracker/internal/service"
	"log"
	"log/slog"
	"os"
)

func main() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)

	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	s := service.New(logger, conf)
	if err := s.StartTracking(); err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(conf, logger, s)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
