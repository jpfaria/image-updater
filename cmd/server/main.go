package main

import (
	"context"
	"os"

	"github.com/jpfaria/image-updater/internal/config"
	"github.com/jpfaria/image-updater/internal/server"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/rs/zerolog/v1"
	"github.com/xgodev/boost/wrapper/log"
)

func init() {
	// Configure logging
	os.Setenv("BOOST_FACTORY_ZEROLOG_CONSOLE_FORMATTER", "JSON")
	os.Setenv("BOOST_FACTORY_ZEROLOG_CONSOLE_LEVEL", "DEBUG")
}

func main() {
	// Initialize Boost
	boost.Start()

	// Set up logger
	log.Set(zerolog.NewLogger())
	log.Info("Starting Image Updater service")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create and start server
	ctx := context.Background()
	srv, err := server.New(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Start the server
	if err := srv.Start(ctx); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
