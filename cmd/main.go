package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/Dorrrke/golang-test-task/docs"
	"github.com/Dorrrke/golang-test-task/internal/config"
	"github.com/Dorrrke/golang-test-task/internal/storage"
	"github.com/Dorrrke/golang-test-task/pkg/api"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// @title Blueprint Swagger API
// @version 1.0
// @description Swagger API for Golang Project Blueprint.

// @host localhost:8080
// @BasePath /test_task/api

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting service")
	log.Debug("Server config:", slog.Any("Config", cfg))

	conn := initDB(cfg.Storage, log)
	storage, err := storage.New(context.Background(), conn, log)
	if err != nil {
		log.Error("Error init database" + err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	api := api.New(log, storage, cfg.ServerConfig.Timeout)

	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))
	r.Route("/test_task/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/", api.AddUserHandler)
			r.Get("/{id}", api.GetUserHandler)
		})
		r.Get("/users", api.GetAllUsersHandler)
	})

	server := &http.Server{
		Addr:           cfg.ServerConfig.GetServerAddr(),
		Handler:        r,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	go func() {
		if err := run(server, log); err != nil {
			log.Error("Server is stoped, error:" + err.Error())
			os.Exit(1)
		}
	}()

	log.Info("Application started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("Application Shutting Down")

	if err := server.Shutdown(context.Background()); err != nil {
		log.Error("error occured on server shutting down:" + err.Error())
	}
	conn.Close()

}

func run(server *http.Server, logger *slog.Logger) error {
	const op = "main.Run"
	log := logger.With(slog.String("op", op))

	log.Debug("Server addr", slog.String("Addr:", server.Addr))
	log.Info("[SERVER STARTED]", slog.String("Server addr", server.Addr))
	return server.ListenAndServe()

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func initDB(stroageCfg config.StorageConfig, logger *slog.Logger) *pgxpool.Pool {
	var dbURL string = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", stroageCfg.User, stroageCfg.Pass, stroageCfg.Host, stroageCfg.Port, stroageCfg.DbName)
	const op = "server.AddUserHandler"
	log := logger.With(slog.String("op", op))
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Error("Error wile init db driver: " + err.Error())
		panic(err)
	}
	return pool

}
