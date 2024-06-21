package service

import (
	"context"
	"fmt"
	"github.com/eachekalina/shortlink/internal/model"
	"log/slog"
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
	PutLink(ctx context.Context, path string, link string) error
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
		return "", fmt.Errorf("failed to generate path: %w", err)
	}
	shortLink := model.ShortLink{
		Path: path,
		Link: link,
	}
	err = s.repo.CreateShortLink(ctx, shortLink)
	if err != nil {
		return "", fmt.Errorf("failed to create short link: %w", err)
	}
	err = s.cache.PutLink(ctx, path, link)
	if err != nil {
		slog.Warn("failed to write to cache", "path", path)
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
		return "", fmt.Errorf("failed to get the link: %w", err)
	}
	err = s.cache.PutLink(ctx, path, link.Link)
	if err != nil {
		slog.Warn("failed to write to cache", "path", path)
	}
	return link.Link, nil
}
