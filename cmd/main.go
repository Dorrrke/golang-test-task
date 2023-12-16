package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/Dorrrke/golang-test-task/internal/api"
	"github.com/Dorrrke/golang-test-task/internal/config"
	"github.com/Dorrrke/golang-test-task/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// @title Blueprint Swagger API
// @version 1.0
// @description Swagger API for Golang Project Blueprint.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email martin7.heinz@gmail.com

// @license.name MIT
// @license.url https://github.com/MartinHeinz/go-project-blueprint/blob/master/LICENSE

// @BasePath /api/v1

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting service")
	log.Debug("Server config:", slog.Any("Config", cfg))

	conn := initDB(cfg.StoragePath, log)
	storage := storage.New(conn, log)
	defer conn.Close()

	server := api.New(log, storage, cfg.ServerConfig.Timeout)

	err := run(*server, *cfg, log)
	if err != nil {
		panic(err)
	}

}

func run(s api.Server, cfh config.Config, logger *slog.Logger) error {
	const op = "main.Run"
	log := logger.With(slog.String("op", op))
	r := chi.NewRouter()

	r.Route("/test_task/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/", s.AddUserHandler)
			r.Get("/{id}", s.GetUserHandler)
		})
		r.Get("/users", s.GetAllUsersHandler)
	})

	log.Debug("Server addr", slog.String("Addr:", cfh.ServerConfig.GetServerAddr()))
	return http.ListenAndServe(cfh.ServerConfig.GetServerAddr(), r)

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

func initDB(DBAddr string, logger *slog.Logger) *pgxpool.Pool {
	const op = "server.AddUserHandler"
	log := logger.With(slog.String("op", op))
	pool, err := pgxpool.New(context.Background(), DBAddr)
	if err != nil {
		log.Error("Error wile init db driver: " + err.Error())
		panic(err)
	}
	return pool

}
