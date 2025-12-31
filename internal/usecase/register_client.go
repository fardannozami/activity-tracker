package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/fardannozami/activity-tracker/internal/domain/service"
	"github.com/fardannozami/activity-tracker/internal/repo/postgres"
)

type RegisterClientUC struct {
	Client *postgres.ClientRepo
}

type RegisterOut struct {
	ClientID string `json:"client_id"`
	APIKey   string `json:"api_key"`
}

func NewRegisterClientUC(client *postgres.ClientRepo) *RegisterClientUC {
	return &RegisterClientUC{Client: client}
}

func (r *RegisterClientUC) Execute(ctx context.Context, name, email string) (RegisterOut, error) {
	apiKey := randomHex(32)
	hash, err := service.HashAPIKey(apiKey)
	if err != nil {
		return RegisterOut{}, err
	}

	row, err := r.Client.Create(ctx, name, email, hash)
	if err != nil {
		return RegisterOut{}, err
	}

	return RegisterOut{ClientID: row.ID, APIKey: apiKey}, nil
}

func randomHex(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}
