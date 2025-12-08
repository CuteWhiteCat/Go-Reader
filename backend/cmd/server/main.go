package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/whitecat/go-reader/internal/api"
	"github.com/whitecat/go-reader/internal/api/handlers"
	"github.com/whitecat/go-reader/internal/config"
	"github.com/whitecat/go-reader/internal/repository"
	"github.com/whitecat/go-reader/internal/service"
)

func main() {
	// Setup logging
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	logrus.Info("Starting Go-Reader server...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := config.InitDatabase(cfg.Database.Path)
	if err != nil {
		logrus.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	bookRepo := repository.NewBookRepository(db)
	chapterRepo := repository.NewChapterRepository(db)
	tagRepo := repository.NewTagRepository(db)
	progressRepo := repository.NewProgressRepository(db)
	bookmarkRepo := repository.NewBookmarkRepository(db)

	// Initialize services
	bookService := service.NewBookService(bookRepo, chapterRepo, tagRepo)
	tagService := service.NewTagService(tagRepo)
	progressService := service.NewProgressService(progressRepo, bookmarkRepo)
	crawlerService := service.NewCrawlerServiceWithCoverDir(bookRepo, chapterRepo, cfg.Storage.CoversDir)

	// Initialize handlers
	bookHandler := handlers.NewBookHandler(bookService)
	tagHandler := handlers.NewTagHandler(tagService)
	progressHandler := handlers.NewProgressHandler(progressService)
	crawlerHandler := handlers.NewCrawlerHandler(crawlerService)

	// Setup router
	router := api.NewRouter(bookHandler, tagHandler, progressHandler, crawlerHandler)
	r := router.SetupRoutes()
	r.Handle("/covers/*", http.StripPrefix("/covers/", http.FileServer(http.Dir(cfg.Storage.CoversDir))))

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	logrus.Infof("Server listening on %s", addr)

	// Setup graceful shutdown
	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- http.ListenAndServe(addr, r)
	}()

	// Wait for shutdown signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		logrus.Fatalf("Server error: %v", err)
	case sig := <-shutdown:
		logrus.Infof("Received signal %v, shutting down...", sig)
	}

	logrus.Info("Server stopped")
}
