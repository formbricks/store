package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	_ "github.com/lib/pq"

	"entgo.io/ent/dialect/sql"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/formbricks/formbricks-rewrite/apps/hub/internal/api"
	"github.com/formbricks/formbricks-rewrite/apps/hub/internal/config"
	"github.com/formbricks/formbricks-rewrite/apps/hub/internal/embedding"
	"github.com/formbricks/formbricks-rewrite/apps/hub/internal/enrichment"
	"github.com/formbricks/formbricks-rewrite/apps/hub/internal/ent"
	"github.com/formbricks/formbricks-rewrite/apps/hub/internal/queue"
	"github.com/formbricks/formbricks-rewrite/apps/hub/internal/webhook"
	"github.com/formbricks/formbricks-rewrite/apps/hub/internal/worker"
)

func main() {
	// Create a CLI app with Huma's service configuration
	cli := humacli.New(func(hooks humacli.Hooks, cfg *config.Config) {
		// Setup logger
		logLevel := slog.LevelInfo
		switch cfg.LogLevel {
		case "debug":
			logLevel = slog.LevelDebug
		case "warn":
			logLevel = slog.LevelWarn
		case "error":
			logLevel = slog.LevelError
		}

		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		}))

		// Connect to database
		drv, err := sql.Open("postgres", cfg.DatabaseURL)
		if err != nil {
			logger.Error("failed to connect to database", "error", err)
			os.Exit(1)
		}

		// Configure connection pool
		db := drv.DB()
		db.SetMaxOpenConns(cfg.DBMaxOpenConns)
		db.SetMaxIdleConns(cfg.DBMaxIdleConns)
		db.SetConnMaxLifetime(time.Duration(cfg.DBConnMaxLifetime) * time.Minute)
		db.SetConnMaxIdleTime(time.Duration(cfg.DBConnMaxIdleTime) * time.Minute)

		logger.Info("database connected",
			"url", cfg.DatabaseURL,
			"max_open_conns", cfg.DBMaxOpenConns,
			"max_idle_conns", cfg.DBMaxIdleConns,
			"conn_max_lifetime_min", cfg.DBConnMaxLifetime,
			"conn_max_idle_time_min", cfg.DBConnMaxIdleTime)

		// Create Ent client with the configured driver
		client := ent.NewClient(ent.Driver(drv))

		// Run migrations
		if err := client.Schema.Create(context.Background()); err != nil {
			logger.Error("failed to run migrations", "error", err)
			os.Exit(1)
		}

		// Create webhook dispatcher
		webhookURLs := cfg.GetWebhookURLs()
		dispatcher := webhook.NewDispatcher(webhookURLs, logger)
		if len(webhookURLs) > 0 {
			logger.Info("webhook dispatcher initialized", "urls", webhookURLs)
		} else {
			logger.Info("webhook dispatcher initialized with no URLs (webhooks disabled)")
		}

		// Initialize AI services and workers if configured
		var enricher *worker.Enricher
		var enrichmentQueue queue.Queue

		// Check if either enrichment or embedding is enabled
		if cfg.IsEnrichmentEnabled() || cfg.IsEmbeddingEnabled() {
			// Create queue (shared by both enrichment and embedding jobs)
			enrichmentQueue = queue.NewPostgresQueue(client)

			// Create enrichment service if configured
			var enrichmentService *enrichment.Service
			if cfg.IsEnrichmentEnabled() {
				enrichmentService = enrichment.NewService(
					cfg.OpenAIKey,
					cfg.OpenAIEnrichmentModel,
					cfg.EnrichmentTimeout,
					logger,
				)
				logger.Info("enrichment service initialized", "model", cfg.OpenAIEnrichmentModel)
			}

			// Create embedding service if configured
			var embeddingService *embedding.Service
			if cfg.IsEmbeddingEnabled() {
				embeddingService = embedding.NewService(
					cfg.OpenAIKey,
					cfg.OpenAIEmbeddingModel,
					cfg.EnrichmentTimeout,
					logger,
				)
				logger.Info("embedding service initialized", "model", cfg.OpenAIEmbeddingModel)
			}

			// Create worker pool (processes both types of jobs)
			pollInterval := time.Duration(cfg.EnrichmentPollInterval) * time.Second
			enricher = worker.NewEnricher(
				enrichmentQueue,
				enrichmentService,
				embeddingService,
				client,
				dispatcher,
				cfg.EnrichmentWorkers,
				pollInterval,
				logger,
			)
		}

		// Create server (pass queue for enqueueing jobs)
		server := api.NewServer(cfg, client, dispatcher, enrichmentQueue, logger)

		// Tell the CLI how to start the server
		hooks.OnStart(func() {
			logger.Info("starting Hub service",
				"port", cfg.Port,
				"environment", cfg.Environment,
				"docs_url", fmt.Sprintf("http://localhost:%d/docs", cfg.Port),
				"openapi_url", fmt.Sprintf("http://localhost:%d/openapi.json", cfg.Port))

			ctx := context.Background()

			// Start enrichment workers if configured
			if enricher != nil {
				go enricher.Start(ctx)
			}

			// Start HTTP server
			if err := server.Start(ctx); err != nil {
				logger.Error("server error", "error", err)
				os.Exit(1)
			}
		})

		// Handle graceful shutdown
		hooks.OnStop(func() {
			logger.Info("shutting down gracefully...")

			// Stop enrichment workers if running
			if enricher != nil {
				enricher.Stop()
			}

			// Shutdown webhook dispatcher with 30 second timeout
			if dispatcher != nil {
				if err := dispatcher.Shutdown(30 * time.Second); err != nil {
					logger.Error("webhook dispatcher shutdown error", "error", err)
				}
			}

			if err := client.Close(); err != nil {
				logger.Error("failed to close database connection", "error", err)
			}
		})
	})

	// Run the CLI - when passed no commands, it starts the server
	cli.Run()
}
