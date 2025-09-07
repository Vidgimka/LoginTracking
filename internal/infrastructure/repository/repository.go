package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Vidgimka/LoginTracking/internal/domain"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type postgresGormRepo struct {
	db *gorm.DB
}

func NewPostgresGormRepo(db *gorm.DB) *postgresGormRepo {
	return &postgresGormRepo{
		db: db,
	}
}

func (r *postgresGormRepo) CreateData(ctx context.Context, data []domain.Data) error {
	r.db.Create(&data)
	return nil
}

func (r *postgresGormRepo) GetLines(ctx context.Context, login string, start, end time.Time) ([]domain.LineData, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("ctx.Err: %w", err)
	}

	rows, err := r.db.WithContext(ctx).Raw("SELECT login, session_id, array_agg(ARRAY(lat, lon)) AS coordinates, MIN(created_at) AS start_time, MAX(created_at) AS end_time FROM data WHERE login = ? AND created_at BETWEEN ? AND ? GROUP BY session_id, login ORDER BY session_id ", login, start, end).Rows()
	if err != nil {
		return nil, fmt.Errorf("db.Raw.Rows: %w", err)
	}
	defer rows.Close()

	var responce []domain.LineData

	for rows.Next() {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("ctx.Err: %w", err)
		}
		var line domain.LineData
		if err := rows.Scan(&line.Login, &line.SessionId, pq.Array(&line.Coordinates), &line.StartTime, &line.EndTime); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		responce = append(responce, line)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %v", err)
	}
	return responce, nil
}

func (db *postgresGormRepo) GetPointByDatetime(ctx context.Context, login string, CreatedAt time.Time) ([]domain.PointData, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("ctx.Err: %w", err)
	}

	rows, err := db.db.WithContext(ctx).Raw("SELECT login, session_id, lat, lon, station_distance, created_at FROM data WHERE login = ? AND created_at = ?", login, CreatedAt).Rows()
	if err != nil {
		return nil, fmt.Errorf("db.Raw.Rows: %w", err)
	}
	defer rows.Close()

	var responce []domain.PointData

	for rows.Next() {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("ctx.Err: %w", err)
		}
		var point domain.PointData
		if err := rows.Scan(&point.Login, &point.SessionId, &point.Coordinates, &point.StationDistance, &point.CreatedAt); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		responce = append(responce, point)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %v", err)
	}
	return responce, nil
}

func (r *postgresGormRepo) GetPointByLogin(ctx context.Context, login string) ([]domain.PointData, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("ctx.Err: %w", err)
	}
	rows, err := r.db.WithContext(ctx).Raw("SELECT login, session_id, lat, lon, station_distance,  created_at  FROM data WHERE login = ?", login).Rows()
	if err != nil {
		return nil, fmt.Errorf("db.Raw.Rows: %w", err)
	}
	defer rows.Close()
	var responce []domain.PointData
	for rows.Next() {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("ctx.Err: %w", err)
		}
		var user domain.PointData
		if err := rows.Scan(&user.Login, &user.SessionId, &user.Coordinates, &user.StationDistance, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		responce = append(responce, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %v", err)
	}
	return responce, nil
}

func (r *postgresGormRepo) GetUsersBySessionId(ctx context.Context, login string, Session_id string) ([]domain.PointData, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("ctx.Err: %w", err)
	}

	rows, err := r.db.WithContext(ctx).Raw("SELECT login, session_id, lat, lon, station_distance, created_at FROM data WHERE login = ? AND session_id = ?", login, Session_id).Rows()
	if err != nil {
		return nil, fmt.Errorf("db.Raw.Rows: %w", err)
	}
	defer rows.Close()
	var responce []domain.PointData
	for rows.Next() {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("ctx.Err: %w", err)
		}
		var user domain.PointData
		if err := rows.Scan(&user.Coordinates, &user.StationDistance, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		responce = append(responce, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %v", err)
	}
	return responce, nil
}
