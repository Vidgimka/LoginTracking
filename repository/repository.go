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
}

type PostgresGormRepo struct {
	db *gorm.DB //  пул соединений с базой данных
}

func NewPostgresGormRepo(db *gorm.DB) PostgresGormRepoInterfase {
	return &PostgresGormRepo{
		db: db,
	}
}

func (r *PostgresGormRepo) GetByAllUser(ctx context.Context, write io.Writer) error {
	write.Write([]byte("["))
	//проверка контекста перед запросом
	if err := ctx.Err(); err != nil {
		return err // обработать
	}

	rows, err := r.db.WithContext(ctx).Raw("SELECT * FROM data").Rows()
	if err != nil {
		return err // обработать
	}

	// // эту часть в роутер
	// c.Header("Content-Type", "application/json")
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error:": "DB error"})
	// }
	// //

	defer rows.Close()
	inFirst := true
	for rows.Next() {
		//проверяем конткест на отмену при заходе на итерацию цикла
		if err := ctx.Err(); err != nil {
			return err // // обработать
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
			fmt.Printf("Serializationerror %v", err)
			continue
		}
		inFirst = false
		write.Write(inJson)
		// c.Writer.Flush() - вручную вызвать в обработчике через джин контекст и http.Flusher и сбрость  через Flush()
	}
	write.Write([]byte("]"))
	return nil
}
