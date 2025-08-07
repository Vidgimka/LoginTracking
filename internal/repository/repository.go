package repository

import (
	"fmt"
	"time"

	"github.com/Vidgimka/LoginTracking.git/internal/domain"
	"github.com/lib/pq"
	"golang.org/x/net/context"
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

func (r *postgresGormRepo) Create(data interface{}) error {
	return r.db.Create(data).Error
}

func (r *postgresGormRepo) GetLines(ctx context.Context, login string, start, end time.Time) ([]domain.LineData, error) {
	if err := ctx.Err(); err != nil {
		return []domain.LineData{}, fmt.Errorf("ctx.Err: %w", err)
	}

	rows, err := r.db.WithContext(ctx).Raw("SELECT login, session_id, array_agg(ARRAY(lat, lon)) AS coordinates, MIN(created_at) AS start_time, MAX(created_at) AS end_time FROM data WHERE login = ? AND created_at BETWEEN ? AND ? GROUP BY session_id, login ORDER BY session_id ", login, start, end).Rows()
	if err != nil {
		return []domain.LineData{}, fmt.Errorf("db.Raw.Rows: %w", err)
	}
	defer rows.Close()

	var responce []domain.LineData

	for rows.Next() {
		if err := ctx.Err(); err != nil {
			return []domain.LineData{}, fmt.Errorf("ctx.Err: %w", err)
		}
		var line domain.LineData
		if err := rows.Scan(&line.Login, &line.Session_id, pq.Array(&line.Coordinates), &line.Start_time, &line.End_time); err != nil {
			return []domain.LineData{}, fmt.Errorf("rows.Scan: %w", err)
		}
		responce = append(responce, line)
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("rows.Err(): %v", err)
	}
	return responce, nil
}

func (db *postgresGormRepo) GetPointByDatetime(ctx context.Context, login string, CreatedAt string) ([]domain.PointData, error) {
	if err := ctx.Err(); err != nil {
		return []domain.PointData{}, fmt.Errorf("ctx.Err: %w", err)
	}

	rows, err := db.db.WithContext(ctx).Raw("SELECT login, session_id, lat, lon, station_distance, created_at FROM data WHERE login = ? AND created_at = ?", login, CreatedAt).Rows()
	if err != nil {
		return []domain.PointData{}, fmt.Errorf("db.Raw.Rows: %w", err)
	}
	defer rows.Close()

	var responce []domain.PointData

	for rows.Next() {
		if err := ctx.Err(); err != nil {
			return []domain.PointData{}, fmt.Errorf("ctx.Err: %w", err)
		}
		var point domain.PointData
		if err := rows.Scan(&point.Login, &point.Session_id, &point.Coordinates, &point.Station_distance, &point.CreatedAt); err != nil {
			return []domain.PointData{}, fmt.Errorf("rows.Scan: %w", err)
		}
		responce = append(responce, point)
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("rows.Err(): %v", err)
	}
	return responce, nil
}

func (r *postgresGormRepo) GetUsers(ctx context.Context) ([]domain.Data, error) {
	if err := ctx.Err(); err != nil {
		return []domain.Data{}, fmt.Errorf("ctx.Err: %w", err)
	}
	rows, err := r.db.WithContext(ctx).Raw("SELECT * FROM data").Rows()
	if err != nil {
		return []domain.Data{}, fmt.Errorf("db.Raw.Rows: %w", err)
	}
	defer rows.Close()
	var responce []domain.Data
	for rows.Next() {
		if err := ctx.Err(); err != nil {
			return []domain.Data{}, fmt.Errorf("ctx.Err: %w", err)
		}
		var user domain.Data
		if err := rows.Scan(&user.Login, &user.SessionId, &user.Subnet, &user.Mountpoint, &user.Station, &user.NtripAgent, &user.ConnectTime,
			&user.TimeSpan, &user.RecievedData, &user.SentData, &user.StatusCode, &user.Latency, &user.SvNum, &user.Coordinaties, &user.Height,
			&user.StationDistance, &user.CreatedAt); err != nil {
			return []domain.Data{}, fmt.Errorf("rows.Scan: %w", err)
		}
		responce = append(responce, user)
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("rows.Err(): %v", err)
	}
	return responce, nil
}

func (r *postgresGormRepo) GetUserByLogin(ctx context.Context, login string) ([]domain.PointData, error) {
	if err := ctx.Err(); err != nil {
		return []domain.PointData{}, fmt.Errorf("ctx.Err: %w", err)
	}
	rows, err := r.db.WithContext(ctx).Raw("SELECT login, session_id, lat, lon, station_distance,  created_at  FROM data WHERE login = ?", login).Rows()
	if err != nil {
		return []domain.PointData{}, fmt.Errorf("db.Raw.Rows: %w", err)
	}
	defer rows.Close()
	var responce []domain.PointData
	for rows.Next() {
		if err := ctx.Err(); err != nil {
			return []domain.PointData{}, fmt.Errorf("ctx.Err: %w", err)
		}
		var user domain.PointData
		if err := rows.Scan(&user.Login, &user.Session_id, &user.Coordinates, &user.Station_distance, &user.CreatedAt); err != nil {
			return []domain.PointData{}, fmt.Errorf("rows.Scan: %w", err)
		}
		responce = append(responce, user)
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("rows.Err(): %v", err)
	}
	return responce, nil
}

func (r *postgresGormRepo) GetUsersBySessionId(ctx context.Context, login string, Session_id string) ([]domain.PointData, error) {
	if err := ctx.Err(); err != nil {
		return []domain.PointData{}, fmt.Errorf("ctx.Err: %w", err)
	}

	rows, err := r.db.WithContext(ctx).Raw("SELECT login, session_id, lat, lon, station_distance, created_at FROM data WHERE login = ? AND session_id = ?", login, Session_id).Rows()
	if err != nil {
		return []domain.PointData{}, fmt.Errorf("db.Raw.Rows: %w", err)
	}
	defer rows.Close()
	var responce []domain.PointData
	for rows.Next() {
		if err := ctx.Err(); err != nil {
			return []domain.PointData{}, fmt.Errorf("ctx.Err: %w", err)
		}
		var user domain.PointData
		if err := rows.Scan(&user.Coordinates, &user.Station_distance, &user.CreatedAt); err != nil {
			return []domain.PointData{}, fmt.Errorf("rows.Scan: %w", err)
		}
		responce = append(responce, user)
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("rows.Err(): %v", err)
	}
	return responce, nil
}
