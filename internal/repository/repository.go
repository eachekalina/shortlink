package repository

import (
	"context"
	"github.com/eachekalina/shortlink/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	_, err := r.conn.Exec(ctx, "insert into short_link (path, link) values ($1, $2);", link.Path, link.Link)
	return err
}

func (r *Repository) GetLink(ctx context.Context, path string) (model.ShortLink, error) {
	row := r.conn.QueryRow(ctx, "select link from short_link where path = $1;", path)
	link := model.ShortLink{Path: path}
	err := row.Scan(&link.Link)
	if err != nil {
		return model.ShortLink{}, err
	}
	return link, nil
}
