package repository

import (
	"context"
	"errors"
	"github.com/eachekalina/shortlink/internal/errs"
	"github.com/eachekalina/shortlink/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"
)

type Conn interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Repository struct {
	conn Conn
}

func New(conn Conn) *Repository {
	return &Repository{conn: conn}
}

func (r *Repository) CreateShortLink(ctx context.Context, link model.ShortLink) error {
	slog.Debug("saving link to database", "path", link.Path, "link", link.Link)
	_, err := r.conn.Exec(ctx, "insert into short_link (path, link) values ($1, $2);", link.Path, link.Link)
	if err != nil {
		slog.Error("failed to save link to database", "path", link.Path, "link", link.Link, "err", err)
		return err
	}
	return nil
}

func (r *Repository) GetLink(ctx context.Context, path string) (model.ShortLink, error) {
	slog.Debug("trying to get link from database", "path", path)
	row := r.conn.QueryRow(ctx, "select link from short_link where path = $1;", path)
	link := model.ShortLink{Path: path}
	err := row.Scan(&link.Link)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Debug("link not found in database", "path", path)
			return model.ShortLink{}, errs.ErrNotFound
		}
		slog.Error("failed to access the database", "path", path, "err", err)
		return model.ShortLink{}, err
	}
	return link, nil
}
