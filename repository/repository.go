package repository

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/Vidgimka/LoginTracking.git/models"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type PostgresGormRepoInterfase interface {
	GetByAllUser(ctx context.Context, write io.Writer) error
	GetByLogin(ctx context.Context, write io.Writer, login string) error
	GetBySessionId(ctx context.Context, write io.Writer, login string, Session_id string) error
	GetByDatetime(ctx context.Context, write io.Writer, login string, CreatedAt string) error
}

type postgresGormRepo struct {
	db *gorm.DB //  пул соединений с базой данных
}

func NewPostgresGormRepo(db *gorm.DB) PostgresGormRepoInterfase {
	return &postgresGormRepo{
		db: db,
	}
}

func (db *postgresGormRepo) GetByDatetime(ctx context.Context, write io.Writer, login string, CreatedAt string) error {
	write.Write([]byte("["))

	if err := ctx.Err(); err != nil {
		return fmt.Errorf("gin contrxt error: %w", err)
	}

	rows, err := db.db.WithContext(ctx).Raw("SELECT login, session_id, lat, lon, station_distance, created_at FROM data WHERE login = ? AND created_at = ?", login, CreatedAt).Rows()
	if err != nil {
		return fmt.Errorf("db request error: %w", err)
	}
	defer rows.Close()

	first := true
	for rows.Next() {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("db query iteration error: %w", err)
		}
		var r models.ResponseData
		if err := rows.Scan(&r.Login, &r.Session_id, &r.Lat, &r.Lon, &r.Station_distance, &r.CreatedAt); err != nil {
			fmt.Printf("Scan error %v", err)
			continue
		}
		inJson, err := json.Marshal(r) // тут скорее всего кодируем но чуть чуть другом методом
		if err != nil {
			fmt.Printf("Serialization error %v", err)
			continue
		}
		if !first {
			write.Write([]byte(","))
		}
		write.Write(inJson)
		first = false
	}
	write.Write([]byte("]"))
	return nil
}

func (r *postgresGormRepo) GetByAllUser(ctx context.Context, write io.Writer) error {
	write.Write([]byte("["))
	//проверка контекста перед запросом
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("gin contrxt error: %w", err)
	}
	rows, err := r.db.WithContext(ctx).Raw("SELECT * FROM data").Rows()
	if err != nil {
		return fmt.Errorf("db requers error: %w", err)
	}
	defer rows.Close()

	inFirst := true
	for rows.Next() {
		//проверяем конткест на отмену при заходе на итерацию цикла
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("db query iteration error: %w", err)
		}
		var fD models.Data
		if err := rows.Scan(&fD.Login, &fD.SessionId, &fD.Subnet, &fD.Mountpoint, &fD.Station, &fD.NtripAgent, &fD.ConnectTime,
			&fD.TimeSpan, &fD.RecievedData, &fD.SentData, &fD.StatusCode, &fD.Latency, &fD.SvNum, &fD.Lat, &fD.Lon, &fD.Height,
			&fD.StationDistance, &fD.CreatedAt); err != nil {
			fmt.Printf("Scan error: %v", err)
			continue
		}
		if !inFirst {
			write.Write([]byte(","))
		}
		inJson, err := json.Marshal(fD)
		if err != nil {
			fmt.Printf("Serialization error %v", err)
			continue
		}
		inFirst = false
		write.Write(inJson)
	}
	write.Write([]byte("]"))
	return nil
}

func (r *postgresGormRepo) GetByLogin(ctx context.Context, write io.Writer, login string) error {
	write.Write([]byte(`{"type": "FeatureCollection","features": [`))

	if err := ctx.Err(); err != nil {
		return fmt.Errorf("gin contrxt error: %w", err)
	}
	rows, err := r.db.WithContext(ctx).Raw("SELECT login, session_id, lat, lon, station_distance,  created_at  FROM data WHERE login = ?", login).Rows()
	if err != nil {
		return fmt.Errorf("db requers error: %w", err)
	}
	defer rows.Close()

	firstElem := true
	for rows.Next() {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("db query iteration error: %w", err)
		}
		var response models.ResponseData
		if err := rows.Scan(&response.Login, &response.Session_id, &response.Lat, &response.Lon, &response.Station_distance, &response.CreatedAt); err != nil {
			fmt.Printf("Scan error %v", err)
			continue
		}

		responseToGeojson, err := models.ResponseToGeojson(response)
		if err != nil {
			fmt.Printf("Scan error %v", err)
			continue
		}

		inJson, err := json.Marshal(responseToGeojson)
		if err != nil {
			fmt.Printf("Serializationerror %v", err)
			continue
		}
		if !firstElem {
			write.Write([]byte(","))
		}
		firstElem = false
		write.Write(inJson)

	}
	write.Write([]byte("]}"))
	return nil
}

func (r *postgresGormRepo) GetBySessionId(ctx context.Context, write io.Writer, login string, Session_id string) error {
	write.Write([]byte(`{"type": "FeatureCollection","features": [ {
      "type": "Feature",
      "geometry": {
        "type": "LineString",
        "coordinates": [`))

	if err := ctx.Err(); err != nil {
		return fmt.Errorf("gin contrxt error: %w", err)
	}

	rows, err := r.db.WithContext(ctx).Raw("SELECT login, session_id, lat, lon, station_distance, created_at FROM data WHERE login = ? AND session_id = ?", login, Session_id).Rows()
	if err != nil {
		return fmt.Errorf("db request error: %w", err)
	}
	defer rows.Close()
	first := true
	for rows.Next() {

		if err := ctx.Err(); err != nil {
			return fmt.Errorf("db query iteration error: %w", err)
		}

		var r models.ResponseData
		if err := rows.Scan(&r.Login, &r.Session_id, &r.Lat, &r.Lon, &r.Station_distance, &r.CreatedAt); err != nil {
			fmt.Printf("Scan error %v", err)
			continue
		}

		responseToGeojson, err := models.ResponseToCoord(r)
		if err != nil {
			fmt.Printf("Serialization error %v", err)
			continue
		}
		inJson, err := json.Marshal(responseToGeojson)
		if err != nil {
			fmt.Printf("Serialization error %v", err)
			continue
		}
		if !first {
			write.Write([]byte(","))
		}
		write.Write(inJson)
		first = false

	}
	write.Write([]byte(`]},"properties": {}}]}`))
	return nil
}
