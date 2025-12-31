package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ApiHitRepo struct {
	db *pgxpool.Pool
}

func NewApiHitRepo(db *pgxpool.Pool) *ApiHitRepo { return &ApiHitRepo{db: db} }

func (r *ApiHitRepo) BulkInsert(ctx context.Context, rows []struct {
	ClientID string
	IP       string
	Endpoint string
	TS       string // RFC3339 string
}) error {
	if len(rows) == 0 {
		return nil
	}

	// Build: INSERT INTO api_hits (client_id, ip, endpoint, ts) VALUES ...
	sql := `INSERT INTO api_hits (client_id, ip, endpoint, ts) VALUES `
	args := make([]any, 0, len(rows)*4)

	// ($1,$2,$3,$4),($5,$6,$7,$8)...
	placeholder := 1
	for i := range rows {
		if i > 0 {
			sql += ","
		}
		sql += "("
		for j := 0; j < 4; j++ {
			if j > 0 {
				sql += ","
			}
			sql += "$" + itoa(placeholder)
			placeholder++
		}
		sql += ")"

		args = append(args, rows[i].ClientID, rows[i].IP, rows[i].Endpoint, rows[i].TS)
	}

	_, err := r.db.Exec(ctx, sql, args...)
	return err
}

func itoa(n int) string {
	// tiny helper
	if n == 0 {
		return "0"
	}
	b := make([]byte, 0, 16)
	for n > 0 {
		b = append(b, byte('0'+n%10))
		n /= 10
	}
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}
