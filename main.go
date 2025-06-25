package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ResponseData struct {
	Login            string    `json:"login"`
	Session_id       int       `json:"session_id"`
	Lat              float64   `json:"lat"`
	Lon              float64   `json:"lon"`
	Station_distance float64   `json:"station_distance"`
	CreatedAt        time.Time `json:"сreated_at"`
}

type GeoData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    []Data `json:"data"`
}

type Data struct {
	Login           string    `json:"login" gorm:"index:idx_login_сreated_at"`
	SessionId       int       `json:"session_id"`
	Subnet          string    `json:"subnet"`
	Mountpoint      string    `json:"mountpoint"`
	Station         string    `json:"station"`
	NtripAgent      string    `json:"ntrip_agent"`
	ConnectTime     int       `json:"connect_time"`
	TimeSpan        int       `json:"time_span"`
	RecievedData    float64   `json:"recieved_data"`
	SentData        float64   `json:"sent_data"`
	StatusCode      int       `json:"status_code"`
	Latency         int       `json:"latency"`
	SvNum           int       `json:"sv_num"`
	Lat             float64   `json:"lat"`
	Lon             float64   `json:"lon"`
	Height          float64   `json:"height"`
	StationDistance float64   `json:"station_distance"`
	CreatedAt       time.Time `json:"сreated_at" gorm:"index:idx_login_сreated_at"`
}

type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Type       string                 `json:"type"`
	Geometry   Geometry               `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

// Переопределяем имя таблицы
func (Data) TableName() string {
	return "loginonline" // желаемое имя таблицы
}

// Проверяем наличие файла окружения
func init() {
	// загружаем значения из .env в систему
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func ReadDataFromAPI() GeoData { // читаем и записываем данные с API
	url := os.Getenv("URL")
	var usersOnline GeoData    // записываем в переменную UsersOnline  данные из тела ответа
	resp, err := http.Get(url) // запрос с APi
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("Response status:", resp.Status)
	d, _ := io.ReadAll(resp.Body) // читаем данные и возвращаем тело ответа в байтах
	if err := json.Unmarshal(d, &usersOnline); err != nil {
		log.Fatal(err.Error())
	}
	return usersOnline
}

func Init() *gorm.DB {
	// функция подключения к БД
	dsn := "host=localhost user=postgres password=postgres dbname=OnlineUsersIist port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to database.")
	}
	err = db.AutoMigrate(&Data{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	return db
}

func RunTaskEverySecond(ctx context.Context, stop <-chan struct{}) {
	// Запись полученных с API данных в БД
	var db *gorm.DB = Init()
	ticker1 := time.NewTicker(time.Second)
	defer ticker1.Stop()
	for {
		select {
		case <-ticker1.C:
			fmt.Println("Running task every second")
			data := ReadDataFromAPI().Data // помещаем в переменную вычетанные данные DATA
			db.Create(&data)               // запись в БД
			log.Println("'Datetime' column added.")
			fmt.Println("Database entry complete")
		case <-stop:
			fmt.Println("no data received")
			return // выход из цикла
		case <-ctx.Done():
			fmt.Println("the user interrupted the program")
			return
		}
	}
}

func main() {
	var db *gorm.DB = Init() // запрос делаетя один раз в main, и далее везде используется для запросов
	// реализация в основном пототке graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	stop := make(chan struct{})
	time.Sleep(time.Second)
	go RunTaskEverySecond(ctx, stop) // если вынести функцию отделно, а потом
	//вызвать горутиной, то горутины синхронизируются (Channel Synchronization)
	// даем поработать алгоритму
	time.Sleep(3 * time.Second) //без этого гоурутина не успевает срабоать
	close(stop)                 // закрывает горутину main

	// Блок с сервером

	router := gin.Default()

	router.GET("/UsersOnline2", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Writer.Write([]byte("["))
		rows, err := db.Raw("SELECT * FROM loginonline").Rows()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error:": "DB error"})
		}
		defer rows.Close()
		inFirst := true
		for rows.Next() {
			var fD Data
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
	//curl http://localhost:8080/UsersOnline2

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
		rows, err := db.Raw("SELECT login, session_id, lat, lon, station_distance, created_at FROM loginonline WHERE login = ? AND session_id = ?", login, Session_id).Rows()
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
	//curl http://localhost:8080/UsersOnline2/nje232/4968

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
		rows, err := db.Raw("SELECT login, session_id, lat, lon, station_distance, created_at FROM loginonline WHERE login = ? AND created_at = ?", login, CreatedAt).Rows()
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
	//curl http://localhost:8080/UsersOnline2/nje232/date/2025-06-22T21:02:30.896313+03:00

	router.GET("/UsersOnline2/:login", func(c *gin.Context) {
		//функция с HTTP- стримингом
		login := c.Param("login")
		if login == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Login cannot be empty"})
			return
		}
		rows, err := db.Raw("SELECT login, session_id, lat, lon, station_distance,  created_at  FROM loginonline WHERE login = ?", login).Rows()
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
	//curl http://localhost:8080/UsersOnline2/aza235
	router.Run("localhost:8080")
}
func responseToGeojson(response ResponseData) (Feature, error) {
	return Feature{Type: "Feature",
		Geometry: Geometry{Type: "Point",
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
