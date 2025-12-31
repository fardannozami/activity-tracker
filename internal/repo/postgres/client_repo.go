package postgres

import (
	"context"

	domainErr "github.com/fardannozami/activity-tracker/internal/domain/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientRepo struct {
	db *pgxpool.Pool
}

func NewClientRepo(db *pgxpool.Pool) *ClientRepo { return &ClientRepo{db: db} }

type ClientRow struct {
	ID         string
	Name       string
	Email      string
	APIKeyHash string
}

func (c *ClientRepo) Create(ctx context.Context, name, email, APIKeyHash string) (ClientRow, error) {
	id := uuid.New().String()
	query := `INSERT INTO clients (id,name,email,api_key_hash) VALUES ($1,$2,$3,$4)`
	_, err := c.db.Exec(ctx, query, id, name, email, APIKeyHash)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "23505":
				return ClientRow{}, domainErr.ErrClientEmailAlreadyExists
			}
		}
		return ClientRow{}, err
	}

	return ClientRow{ID: id, Name: name, Email: email, APIKeyHash: APIKeyHash}, nil
}
