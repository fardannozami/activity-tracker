package entity

import "time"

type Client struct {
	ID         string
	Name       string
	Email      string
	APIKeyHash string
	CreatedAt  time.Time
}
