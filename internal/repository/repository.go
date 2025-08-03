package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/Vidgimka/LoginTracking.git/internal/models"
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

func (r *postgresGormRepo) GetLines(ctx context.Context, write io.Writer, login string, start, end time.Time) error {
	write.Write([]byte(`{"type": "FeatureCollection","features": [`))
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("ctx.Err: %w", err)
	}

	var stringCoord []byte

	rows, err := r.db.WithContext(ctx).Raw("SELECT login, session_id, json_agg(json_build_array(lat, lon)) AS coordinates, MIN(created_at) AS start_time, MAX(created_at) AS end_time FROM data WHERE login = ? AND created_at BETWEEN ? AND ? GROUP BY session_id, login ORDER BY session_id ", login, start, end).Rows()
	if err != nil {
		return fmt.Errorf("db.Raw.Rows: %w", err)
	}
	defer rows.Close()
	firstElem := true
	for rows.Next() {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("ctx.Err: %w", err)
		}
		var response models.ResponseForLine
		if err := rows.Scan(&response.Login, &response.Session_id, &stringCoord, &response.Start_time, &response.End_time); err != nil {
			return fmt.Errorf("rows.Scan: %w", err)
		}

		err := json.Unmarshal(stringCoord, &response.Coordinates)
		if err != nil {
			fmt.Printf("Parce coordinate error: %v", err)
			continue
		}

		responseToGeojson, err := models.LineToSessionIdGeojson(response)
		if err != nil {
			fmt.Printf("Conver to geojson error %v", err)
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

		responseToGeojson, err := models.ResponseToPointGeojson(response)
		if err != nil {
			fmt.Printf("Conver to geojson error %v", err)
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
