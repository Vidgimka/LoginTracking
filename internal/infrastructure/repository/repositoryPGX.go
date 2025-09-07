package repository

import (
	"context"
	"time"

	"github.com/Vidgimka/LoginTracking/internal/domain"
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

func (r *postgresPgx) CreateData(ctx context.Context, data []domain.Data) error
func (r *postgresPgx) GetLines(ctx context.Context, login string, start, end time.Time) ([]domain.LineData, error)
func (r *postgresPgx) GetPointByDatetime(ctx context.Context, login string, CreatedAt time.Time) ([]domain.PointData, error)
func (r *postgresPgx) GetPointByLogin(ctx context.Context, login string) ([]domain.PointData, error)
