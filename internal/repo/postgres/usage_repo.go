package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UsageRepo struct{ db *pgxpool.Pool }

func NewUsageRepo(db *pgxpool.Pool) *UsageRepo { return &UsageRepo{db: db} }

type DailyUsageRow struct {
	Day   time.Time
	Total int64
}

func (r *UsageRepo) GetDailyLast7(ctx context.Context, clientID string) ([]DailyUsageRow, error) {
	rows, err := r.db.Query(ctx, `
		SELECT day, total
		FROM daily_usage
		WHERE client_id=$1
		ORDER BY day DESC
		LIMIT 7
	`, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []DailyUsageRow
	for rows.Next() {
		var d time.Time
		var t int64
		if err := rows.Scan(&d, &t); err != nil {
			return nil, err
		}
		out = append(out, DailyUsageRow{Day: d, Total: t})
	}
	return out, nil
}

type TopRow struct {
	ClientID string
	Total    int64
}

type UsageHit struct {
	ClientID  string
	Timestamp time.Time
}

type UsageWriter interface {
	UpsertAggregates(ctx context.Context, hits []UsageHit) error
}

func (r *UsageRepo) GetTopLast24Hours(ctx context.Context, limit int) ([]TopRow, error) {
	// Sum hourly buckets in last 24 hours
	rows, err := r.db.Query(ctx, `
		SELECT client_id, SUM(total) AS total
		FROM hourly_usage
		WHERE hour >= NOW() - INTERVAL '24 hours'
		GROUP BY client_id
		ORDER BY total DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []TopRow
	for rows.Next() {
		var id string
		var t int64
		if err := rows.Scan(&id, &t); err != nil {
			return nil, err
		}
		out = append(out, TopRow{ClientID: id, Total: t})
	}
	return out, nil
}

func (w *UsageRepo) UpsertAggregates(ctx context.Context, hits []UsageHit) error {
	// naive loop upsert (OK for test). For prod, you’d group & batch.
	for _, h := range hits {
		day := h.Timestamp.Format("2006-01-02")
		hour := h.Timestamp.Truncate(time.Hour).UTC().Format(time.RFC3339)

		_, _ = w.db.Exec(ctx, `
			INSERT INTO daily_usage (client_id, day, total)
			VALUES ($1, $2::date, 1)
			ON CONFLICT (client_id, day)
			DO UPDATE SET total = daily_usage.total + 1
		`, h.ClientID, day)

		_, _ = w.db.Exec(ctx, `
			INSERT INTO hourly_usage (client_id, hour, total)
			VALUES ($1, $2::timestamptz, 1)
			ON CONFLICT (client_id, hour)
			DO UPDATE SET total = hourly_usage.total + 1
		`, h.ClientID, hour)
	}
	return nil
}
