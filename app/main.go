package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Vidgimka/LoginTracking.git/api"
	"github.com/Vidgimka/LoginTracking.git/config"
	"github.com/Vidgimka/LoginTracking.git/db"
	"github.com/Vidgimka/LoginTracking.git/models"
	"github.com/Vidgimka/LoginTracking.git/service"
	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Login            string    `json:"login"`
	Session_id       int       `json:"session_id"`
	Lat              float64   `json:"lat"`
	Lon              float64   `json:"lon"`
	Station_distance float64   `json:"station_distance"`
	CreatedAt        time.Time `json:"сreated_at"`
}

func main() {

	config.LoadEnv()

	client := api.NewHttpClient()
	service := service.NewService(client)

	db, err := db.Init()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	// реализация в основном пототке graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	stop := make(chan struct{})
	time.Sleep(time.Second)
	go service.RunTaskEverySecond(db, ctx, stop)
	time.Sleep(3 * time.Second)
	close(stop)

	// Блок с сервером

	router := gin.Default()

	router.GET("/UsersOnline2/:login/:session_id", func(c *gin.Context) {
		login := c.Param("login")
		if login == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Login cannot be empty"})
			return
		}
		Session_id := c.Param("session_id")
		if Session_id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "SessionId cannot be empty"})
			return
		}
		c.Header("Content-Type", "application/json")
		c.Writer.Write([]byte(`{"type": "FeatureCollection","features": [ {
      "type": "Feature",
      "geometry": {
        "type": "LineString",
        "coordinates": [`))
		rows, err := db.Raw("SELECT login, session_id, lat, lon, station_distance, created_at FROM data WHERE login = ? AND session_id = ?", login, Session_id).Rows()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error:": "DB error"})
		}
		defer rows.Close()
		first := true
		for rows.Next() {
			var r ResponseData
			if err := rows.Scan(&r.Login, &r.Session_id, &r.Lat, &r.Lon, &r.Station_distance, &r.CreatedAt); err != nil {
				fmt.Printf("Scan error %v", err)
				continue
			}

			// responseToGeojson, err := responseToGeojson(r)
			responseToGeojson, err := responseToCoord(r)
			if err != nil {
				fmt.Printf("Serializationerror %v", err)
				continue
			}
			inJson, err := json.Marshal(responseToGeojson)
			if err != nil {
				fmt.Printf("Serializationerror %v", err)
				continue
			}
			if !first {
				c.Writer.Write([]byte(","))
			}
			c.Writer.Write(inJson)
			first = false
			c.Writer.Flush()
		}
		c.Writer.Write([]byte(`]},"properties": {}}]}`))
	})
	// http://localhost:8080/UsersOnline2/nje232/4968

	router.GET("/UsersOnline2", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Writer.Write([]byte("["))
		rows, err := db.Raw("SELECT * FROM data").Rows()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error:": "DB error"})
		}
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
				c.Writer.Write([]byte(","))
			}
			inJson, err := json.Marshal(fD)
			if err != nil {
				fmt.Printf("Serializationerror %v", err)
				continue
			}
			inFirst = false
			c.Writer.Write(inJson)
			c.Writer.Flush()
		}
		c.Writer.Write([]byte("]"))
	})
	// http://localhost:8080/UsersOnline2

	router.GET("/UsersOnline2/:login/date/:datetime", func(c *gin.Context) {
		login := c.Param("login")
		if login == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Login cannot be empty"})
			return
		}
		CreatedAt := c.Param("datetime")
		if CreatedAt == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "SessionId cannot be empty"})
			return
		}
		c.Header("Content-Type", "application/json")
		c.Writer.Write([]byte("["))
		rows, err := db.Raw("SELECT login, session_id, lat, lon, station_distance, created_at FROM data WHERE login = ? AND created_at = ?", login, CreatedAt).Rows()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error:": "DB error"})
		}
		defer rows.Close()
		first := true
		for rows.Next() {
			var r ResponseData
			if err := rows.Scan(&r.Login, &r.Session_id, &r.Lat, &r.Lon, &r.Station_distance, &r.CreatedAt); err != nil {
				fmt.Printf("Scan error %v", err)
				continue
			}

			inJson, err := json.Marshal(r) // тут скорее всего кодируем но чуть чуть другом методом
			if err != nil {
				fmt.Printf("Serializationerror %v", err)
				continue
			}
			if !first {
				c.Writer.Write([]byte(","))
			}

			c.Writer.Write(inJson)
			first = false
			c.Writer.Flush()
		}
		c.Writer.Write([]byte("]"))
	})
	// http://localhost:8080/UsersOnline2/nje232/date/2025-06-22T21:02:30.896313+03:00

	router.GET("/UsersOnline2/:login", func(c *gin.Context) {
		//функция с HTTP- стримингом
		login := c.Param("login")
		if login == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Login cannot be empty"})
			return
		}
		rows, err := db.Raw("SELECT login, session_id, lat, lon, station_distance,  created_at  FROM data WHERE login = ?", login).Rows()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		}
		defer rows.Close()
		// начинаме вручнкю заполнять json ответ
		c.Header("Content-Type", "application/json")
		c.Writer.Write([]byte(`{"type": "FeatureCollection","features": [`))
		firstElem := true
		for rows.Next() {
			var response ResponseData
			if err := rows.Scan(&response.Login, &response.Session_id, &response.Lat, &response.Lon, &response.Station_distance, &response.CreatedAt); err != nil {
				fmt.Printf("Scan error %v", err)
				continue
			}
			responseToGeojson, err := responseToGeojson(response)
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
				c.Writer.Write([]byte(","))
			}
			firstElem = false
			c.Writer.Write(inJson)
			c.Writer.Flush() // сбрасываем буфер, чтобы сервак не накапливал инфц в буфере
		}
		c.Writer.Write([]byte("]}"))
	})
	// http://localhost:8080/UsersOnline2/aza235
	router.Run("localhost:8080")
}
func responseToGeojson(response ResponseData) (models.Feature, error) {
	return models.Feature{Type: "Feature",
		Geometry: models.Geometry{Type: "Point",
			Coordinates: []float64{response.Lon, response.Lat}},
		Properties: map[string]interface{}{
			"Login":      response.Login,
			"Session_id": response.Session_id,
			"CreatedA":   response.CreatedAt,
		}}, nil
}

func responseToCoord(response ResponseData) ([]float64, error) {
	return []float64{response.Lon, response.Lat}, nil
}
