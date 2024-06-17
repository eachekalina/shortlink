package service

import (
	"context"
	"github.com/eachekalina/shortlink/internal/model"
	"log"
	"net/url"
)

type Repository interface {
	CreateShortLink(ctx context.Context, link model.ShortLink) error
	GetLink(ctx context.Context, path string) (model.ShortLink, error)
}

type Generator interface {
	GeneratePath() (string, error)
}

type Cache interface {
	PutLink(ctx context.Context, link model.ShortLink) error
	GetLink(ctx context.Context, path string) (string, error)
}

type Service struct {
	repo    Repository
	gen     Generator
	cache   Cache
	rootURL *url.URL
}

func New(repo Repository, gen Generator, cache Cache, rootURL *url.URL) *Service {
	return &Service{
		repo:    repo,
		gen:     gen,
		cache:   cache,
		rootURL: rootURL,
	}
}

func (s *Service) CreateShortLink(ctx context.Context, link string) (string, error) {
	path, err := s.gen.GeneratePath()
	if err != nil {
		return "", err
	}
	shortLink := model.ShortLink{
		Path: path,
		Link: link,
	}
	err = s.repo.CreateShortLink(ctx, shortLink)
	if err != nil {
		return "", err
	}
	err = s.cache.PutLink(ctx, shortLink)
	if err != nil {
		log.Printf("Failed to write the log: %v\n", err)
	}
	return s.rootURL.JoinPath(path).String(), nil
}

func (s *Service) GetLink(ctx context.Context, path string) (string, error) {
	linkStr, err := s.cache.GetLink(ctx, path)
	if err == nil {
		return linkStr, nil
	}
	link, err := s.repo.GetLink(ctx, path)
	if err != nil {
		return "", err
	}
	return link.Link, nil
}
