package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Vidgimka/LoginTracking/internal/domain"
	"github.com/Vidgimka/LoginTracking/internal/infrastructure/client"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresPgx struct {
	db *pgxpool.Pool
}

func NewPostgresPgxRepo(db *pgxpool.Pool) *postgresPgx {
	return &postgresPgx{
		db: db,
	}
}

func (r *postgresPgx) CreateData(ctx context.Context, usersOnline []client.Data) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("ctx.Err: %w", err)
	}

	sql := "INSERT INTO user_location (login, session_id, mountpoint, station, ntrip_agent, connect_time, time_span, recieved_data, sent_data, status_code, latency, sv_num, lat, lon, height, station_distance, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)"
	batch := &pgx.Batch{}

	for _, data := range usersOnline {
		batch.Queue(sql,
			data.Login, data.SessionId, data.Mountpoint, data.Station, data.NtripAgent, data.ConnectTime, data.TimeSpan, data.RecievedData, data.SentData, data.StatusCode, data.Latency, data.SvNum, data.Lat, data.Lon, data.StationDistance, data.CreatedAt)
	}
	result := r.db.SendBatch(ctx, batch)
	defer result.Close()

	for range usersOnline {
		_, err := result.Exec()
		if err != nil {
			return fmt.Errorf("result.Exec: %w", err)
		}

	}
	return result.Close()
}

func (r *postgresPgx) GetLines(ctx context.Context, login string, start, end time.Time) ([]domain.LineData, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("ctx.Err: %w", err)
	}
	sql := "SELECT login, session_id, lat, lon, MIN(created_at) AS start_time, MAX(created_at) AS end_time FROM data WHERE login = ? AND created_at BETWEEN ? AND ? GROUP BY session_id, login ORDER BY session_id "

	responce := make([]domain.LineData, 0)

	rows, err := r.db.Query(ctx, sql, login, start, end)
	if err != nil {
		return nil, fmt.Errorf("result.Exec: %w", err)
	}
	defer rows.Close()

	for rows.Next() {

		var line domain.LineData
		if err := rows.Scan(&line.Login, &line.SessionId); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		responce = append(responce, line)

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %v", err)
	}
	return responce, nil
}
func (r *postgresPgx) GetPointByDatetime(ctx context.Context, login string, CreatedAt time.Time) ([]domain.PointData, error)
func (r *postgresPgx) GetPointByLogin(ctx context.Context, login string) ([]domain.PointData, error)
