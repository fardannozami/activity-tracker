package usecase

import (
	"context"
	"errors"

	"github.com/fardannozami/activity-tracker/internal/domain/service"
	"github.com/fardannozami/activity-tracker/internal/repo/postgres"
)

type IssueTokenUC struct {
	Clients *postgres.ClientRepo
	Tokens  service.TokenService
}

func (uc *IssueTokenUC) Execute(ctx context.Context, email, apiKey string) (string, string, error) {
	client, ok, err := uc.Clients.GetByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}
	if !ok {
		return "", "", errors.New("invalid credentials")
	}

	if !service.ComparAPIKey(client.APIKeyHash, apiKey) {
		return "", "", errors.New("invalid credentials")
	}

	jwtStr, err := uc.Tokens.Issue(client.ID)
	return jwtStr, client.ID, err
}
