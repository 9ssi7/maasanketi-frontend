package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	log "github.com/mstrYoda/maasanketi.co/pkg/logger"
	"github.com/mstrYoda/maasanketi.co/repository"
	"github.com/mstrYoda/maasanketi.co/worker"
)

func main() {
	repo, err := repository.New()
	if err != nil {
		log.Logger().Panic("failed to create repository: " + err.Error())
	}
	defer repo.Close()

	// Create a new scheduler
	scheduler := gocron.NewScheduler(time.UTC)

	// Create graph generation service
	graphService := worker.NewGraphGenerationService(repo)

	// Schedule graph generation job to run every hour
	// This checks for completed surveys and generates graphs
	scheduler.Every(10).Minute().Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		log.Logger().Info("starting survey graph generation job")
		if err := graphService.GenerateGraphsForCompletedSurveys(ctx); err != nil {
			log.Logger().Error("error generating graphs: " + err.Error())
		}
	})

	// Start the scheduler in a separate goroutine
	scheduler.StartAsync()
	log.Logger().Info("worker started successfully")

	// Wait for termination signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	log.Logger().Info("worker shutting down...")
	scheduler.Stop()
}
