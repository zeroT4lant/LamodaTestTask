package main

import (
	"LamodaTestTask/app/internal/app"
	"LamodaTestTask/app/pkg/logging"
	"context"
	"log"
)

// Точка старта приложения в методе Run()
func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := logging.GetLogger(ctx)

	log.Print("logger initializing")
	ctx = logging.ContextWithLogger(ctx, logger)

	a, err := app.NewApp(ctx)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Run Application")
	err = a.Run(ctx)
	if err != nil {
		logger.Fatal(err)
	}
}
