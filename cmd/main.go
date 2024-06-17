package main

import (
	"context"
	"github.com/eachekalina/shortlink/internal/cache"
	"github.com/eachekalina/shortlink/internal/handler"
	"github.com/eachekalina/shortlink/internal/pathgen"
	"github.com/eachekalina/shortlink/internal/repository"
	"github.com/eachekalina/shortlink/internal/service"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println(err)
		return
	}
	defer pool.Close()
	repo := repository.New(pool)
	gen := pathgen.NewGenerator(12)
	rootURL, err := url.Parse(os.Getenv("ROOT_URL"))
	if err != nil {
		log.Println(err)
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
	http.ListenAndServe(os.Getenv("LISTEN_ADDR"), r)
}
