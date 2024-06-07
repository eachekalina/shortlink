package service

import (
	"context"
	"github.com/eachekalina/shortlink/internal/model"
)

type Repository interface {
	CreateShortLink(ctx context.Context, link model.ShortLink) error
	GetLink(ctx context.Context, path string) (model.ShortLink, error)
}

type Generator interface {
	GeneratePath() (string, error)
}

type Service struct {
	repo Repository
	gen  Generator
}

func New(repo Repository, gen Generator) *Service {
	return &Service{
		repo: repo,
		gen:  gen,
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
	return path, nil
}

func (s *Service) GetLink(ctx context.Context, path string) (string, error) {
	link, err := s.repo.GetLink(ctx, path)
	if err != nil {
		return "", err
	}
	return link.Link, nil
}
