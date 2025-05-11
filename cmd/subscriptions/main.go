package main

import (
	"github.com/blvxme/subpub"
	"github.com/sirupsen/logrus"
	"subscriptions/internal/config"
	"subscriptions/internal/handler"
	"subscriptions/internal/server"
)

func main() {
	logger := logrus.StandardLogger()
	config.ConfigureLogger(logger)
	cfg := config.NewConfig()
	sp := subpub.NewSubPub()
	h := handler.NewHandler(sp, logger)
	s := server.NewServer(h, logger, sp)
	if err := s.Start(cfg.Port); err != nil {
		logger.Panicf("Failed to start server: %+v\n", err)
	}
}
