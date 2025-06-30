package repository

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/Vidgimka/LoginTracking.git/models"
	"gorm.io/gorm"
)

type PostgresGormRepoInterfase interface {
	GetByAllUser(write io.Writer) error
}

type PostgresGormRepo struct {
	db *gorm.DB //  пул соединений с базой данных
}

func NewPostgresGormRepo() PostgresGormRepoInterfase {
	return &PostgresGormRepo{
		db: &gorm.DB{},
	}
}

func (r *PostgresGormRepo) GetByAllUser(write io.Writer) error {

	// отсюда
	// c.Writer.Write([]byte("["))
	write.Write([]byte("["))

	rows, err := r.db.Raw("SELECT * FROM data").Rows()
	// db.Raw("SELECT * FROM data").Rows()
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
		var fD models.Data
		if err := rows.Scan(&fD.Login, &fD.SessionId, &fD.Subnet, &fD.Mountpoint, &fD.Station, &fD.NtripAgent, &fD.ConnectTime,
			&fD.TimeSpan, &fD.RecievedData, &fD.SentData, &fD.StatusCode, &fD.Latency, &fD.SvNum, &fD.Lat, &fD.Lon, &fD.Height,
			&fD.StationDistance, &fD.CreatedAt); err != nil {
			fmt.Printf("Scan error: %v", err)
			continue
		}
		if !inFirst {
			// c.Writer.Write([]byte(","))
			write.Write([]byte(","))
		}
		inJson, err := json.Marshal(fD)
		if err != nil {
			fmt.Printf("Serializationerror %v", err)
			continue
		}
		inFirst = false
		// c.Writer.Write(inJson)
		write.Write(inJson)
		// c.Writer.Flush() - вручную вызвать в обработчике через джин контекст и http.Flusher и сбрость  через Flush()
	}
	// c.Writer.Write([]byte("]"))
	write.Write([]byte("]"))
	// тут конец
	return nil
}
