package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/eachekalina/shortlink/internal/cache"
	"github.com/eachekalina/shortlink/internal/handler"
	"github.com/eachekalina/shortlink/internal/pathgen"
	"github.com/eachekalina/shortlink/internal/repository"
	"github.com/eachekalina/shortlink/internal/service"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

type HandlerTestSuite struct {
	suite.Suite
	pool *pgxpool.Pool
	h    http.Handler
}

func (s *HandlerTestSuite) SetupSuite() {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println(err)
		return
	}
	s.pool = pool
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
	serv := service.New(repo, gen, c, rootURL)
	h := handler.New(serv)
	r := mux.NewRouter()
	r.HandleFunc("/new", h.HandleCreateLink).
		Methods(http.MethodPost)
	r.HandleFunc("/{path:[A-Za-z0-9\\-_]+}", h.HandleLink).
		Methods(http.MethodGet)
	s.h = r
}

func (s *HandlerTestSuite) TearDownSuite() {
	s.pool.Close()
}

func (s *HandlerTestSuite) TestCreateLink() {
	tests := []struct {
		name string
		link string
	}{
		{
			name: "success",
			link: "http://example.com/",
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			b, err := json.Marshal(handler.CreateLinkRequest{Link: test.link})
			if err != nil {
				s.FailNow("failed to marshal request", err)
			}
			req, err := http.NewRequest(http.MethodPost, "/new", bytes.NewReader(b))
			if err != nil {
				s.FailNow("invalid request", err)
			}

			rr := httptest.NewRecorder()
			s.h.ServeHTTP(rr, req)

			s.Equal(http.StatusCreated, rr.Code)
			var resp handler.CreateLinkResponse
			err = json.NewDecoder(rr.Body).Decode(&resp)
			if err != nil {
				s.FailNow("invalid response", err)
			}

			req, err = http.NewRequest(http.MethodGet, resp.Path, nil)
			if err != nil {
				s.FailNow("invalid request", err)
			}
			rr = httptest.NewRecorder()
			s.h.ServeHTTP(rr, req)

			s.Equal(http.StatusMovedPermanently, rr.Code)
			s.Equal(test.link, rr.Header().Get("Location"))
		})
	}
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
