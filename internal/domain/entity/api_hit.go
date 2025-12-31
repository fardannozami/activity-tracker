package entity

import "time"

type ApiHit struct {
	ClientId  string
	IP        string
	Endpoint  string
	Timestamp time.Time
}
