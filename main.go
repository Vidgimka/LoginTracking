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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/**/
type Response struct {
	Login            string  `json:"login"`
	Session_id       int     `json:"session_id"`
	Lat              float64 `json:"lat"`
	Lon              float64 `json:"lon"`
	Station_distance float64 `json:"station_distance"`
}

// чтобы горм правильно определил схему, а именно теблицу DATA
// так данные в БД используются именно из DATA, остальные поля json не исп.
type GeoData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    []Data `json:"data"`
}

type Data struct {
	Login           string    `json:"login"`
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
	CreatedAt       time.Time `json:"datetime"`
}

// Переопределяем имя таблицы
func (Data) TableName() string {
	return "loginonline" // желаемое имя таблицы
}

//var usersOnline GeoData // записываем в переменную UsersOnline  данные из тела ответа
// var UsersOnlineD Data

func ReadFileData() GeoData { // читаем и записываем данные с API
	URL := "https://"
	var usersOnline GeoData    // записываем в переменную UsersOnline  данные из тела ответа
	resp, err := http.Get(URL) // запрос с APi
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

func Init() *gorm.DB { // функция подключения к БД, возвращает объект подключения к БД
	// логика, как при вызове этой функци будет создаваться БД
	//var DB *gorm.DB
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

func RunTaskEverySecond(ctx context.Context, stop <-chan struct{}) { // Запись полученных с API данных в БД
	// совмещаем логику 2х функций, Init() созданию БД
	// и записи ReadFileData().Data в переменную Data конкретного куска данных полученных Data с API
	var db *gorm.DB = Init()
	ticker1 := time.NewTicker(time.Second)
	defer ticker1.Stop()
	for {
		select {
		case <-ticker1.C:
			fmt.Println("Running task every second")
			data := ReadFileData().Data // помещаем в переменную вычетанные данные DATA
			db.Create(&data)            // запись в БД
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
		// делаем функцию, которая будет ходить в бд и возвращать данные,а не писать их в глобальную переменнюу
		//loginOnline := GetUsersOnline(db)
		var loginonline []Data
		response := db.Find(&loginonline)
		if response.Error != nil {
			log.Printf("Database query failed %v", response.Error)
			c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
			return
		}
		log.Println("Data from the database has been received")
		c.JSON(http.StatusOK, loginonline)
	})
	//curl http://localhost:8080/UsersOnline2

	router.GET("/UsersOnline2/:login", func(c *gin.Context) {
		var data []Data
		login := c.Param("login") // записываем в переменную вычитанный логин из URL
		// запускаем цикл, который переберет все значения из переменной баз данных
		// и при собпадении всех записей с заданным логином выведет их в теле ответа
		response := db.Where("login = ?", login).Find(&data)
		if response.Error != nil {
			log.Printf("Database query failed %v", response.Error)
			c.JSON(http.StatusNotFound, gin.H{"message": "Login not found"})
		}
		log.Println("Data from the database has been received")
		c.JSON(http.StatusOK, data)

	})
	//curl http://localhost:8080/UsersOnline2/uralgeometer1

	router.GET("/UsersOnline2/:login/:session_id", func(c *gin.Context) {
		var data []Data
		login := c.Param("login")
		session_id := c.Param("session_id")
		response := db.Where("login = ? AND session_id =?", login, session_id).Find(&data)
		if response.Error != nil {
			log.Printf("Database query failed %v", response.Error)
			c.JSON(http.StatusNotFound, gin.H{"message": "Login not found"})
		}
		log.Println("Data from the database has been received")
		c.JSON(http.StatusOK, data)

	})
	//curl http://localhost:8080/UsersOnline2/yea349/6088
	router.Run("localhost:8080")
}
