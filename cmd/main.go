package main

import (
	"context"
	"errors"
	"github.com/eachekalina/shortlink/internal/cache"
	"github.com/eachekalina/shortlink/internal/handler"
	"github.com/eachekalina/shortlink/internal/pathgen"
	"github.com/eachekalina/shortlink/internal/repository"
	"github.com/eachekalina/shortlink/internal/service"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("failed to connect to database", "err", err)
		return
	}
	defer pool.Close()
	repo := repository.New(pool)
	gen := pathgen.NewGenerator(12)
	rootURL, err := url.Parse(os.Getenv("ROOT_URL"))
	if err != nil {
		slog.Error("invalid root URL", "err", err)
		return
	}

	redisCli := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	c := cache.New(redisCli, 5*time.Minute)
	s := service.New(repo, gen, c, rootURL)
	h := handler.New(s)
	r := mux.NewRouter()
	r.HandleFunc("/new", h.HandleCreateLink).
		Methods(http.MethodPost)
	r.HandleFunc("/{path:[A-Za-z0-9\\-_]+}", h.HandleLink).
		Methods(http.MethodGet)

	serv := http.Server{
		Addr:    os.Getenv("LISTEN_ADDR"),
		Handler: r,
	}

	slog.Info("server is started", "addr", serv.Addr)

	go func() {
		err := serv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server closed with an error", "err", err)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down server")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	err = serv.Shutdown(shutdownCtx)
	if err != nil {
		slog.Error("server shutdown failed", "err", err)
	}
}
