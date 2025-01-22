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

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/**/
type Data struct {
	Login            string  `json:"login"`
	Session_id       int     `json:"session_id"`
	Subnet           string  `json:"subnet"`
	Mountpoint       string  `json:"mountpoint"`
	Station          string  `json:"station"`
	Ntrip_agent      string  `json:"ntrip_agent"`
	Connect_time     int     `json:"connect_time"`
	Time_span        int     `json:"time_span"`
	Recieved_data    float64 `json:"recieved_data"`
	Sent_data        float64 `json:"sent_data"`
	Status_code      int     `json:"status_code"`
	Latency          int     `json:"latency"`
	Sv_num           int     `json:"sv_num"`
	Lat              float64 `json:"lat"`
	Lon              float64 `json:"lon"`
	Height           int     `json:"height"`
	Station_distance float64 `json:"station_distance"`
}

// чтобы горм правильно определил схему, а именно теблицу DATA
// так данные в БД используются именно из DATA, остальные поля json не исп.
type GeoData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    []Data `json:"data"`
}

var UsersOnline GeoData // записываем в переменную UsersOnline  данные из тела ответа
//var UsersOnlineD Data

func ReadFileData() GeoData { // читаем и записываем данные с API
	URL := "https://"

	resp, err := http.Get(URL) // запрос с APi
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("Response status:", resp.Status)
	/*
		io.Copy(os.Stdout, resp.Body) //тест корректного чтения  данных с API
		data := make([]byte, 1014)     // создаем байтовую переменную, чтобы записать респонс от API
			n, err := resp.Body.Read(data) // записываем тело ответа в переменную
			fmt.Println(string(data[:n]))  //вывод в консоль*/
	d, _ := io.ReadAll(resp.Body) // читаем данные и возвращаем тело ответа в байтах
	if err := json.Unmarshal(d, &UsersOnline); err != nil {
		log.Fatal(err.Error())
	}
	return UsersOnline
}

func Init() *gorm.DB {
	var DB *gorm.DB
	dsn := "host=localhost user=postgres password=postgres dbname=OnlineUsersIist port=5432 sslmode=disable"
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("не подключилось к БД")
	}
	DB.AutoMigrate(&Data{})
	return DB
}

// func GetDB() *gorm.DB { // проверяем подкюченеие к базе данных
// 	if Dbase == nil {
// 		Dbase = Init() // проинициализировали бд через Init и присвоили в переменную dbase, так как инициилизация return db
// 		var sleep = time.Duration(1)
// 		for Dbase == nil { // ждем 3 секунды если не подключается
// 			sleep = sleep * 3
// 			fmt.Printf("База данных не доступна. Пожождите %d секунды.\n", sleep)
// 			time.Sleep(sleep * time.Second)
// 			Dbase = Init() // еще раз кладем базу данных в переменную dbase
// 		}
// 	}
// 	return Dbase
// }

func RunTaskEverySecond(ctx context.Context, stop <-chan struct{}) {
	var Dbase *gorm.DB = Init()
	ticker1 := time.NewTicker(time.Second)
	defer ticker1.Stop()
	for {
		select {
		case <-ticker1.C:
			fmt.Println("Running task every second")
			Data := ReadFileData().Data
			Dbase.Create(&Data)
			fmt.Println("Запись в БД завершенна")
		case <-stop:
			fmt.Println("Данные не поступают")
			return // выход из цикла
		case <-ctx.Done():
			fmt.Println("Пользователь прервал программу")
			return
		}
	}
}

func main() {
	// реализация в основном пототке graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	// ctx, cancel := context.WithCancel(context.Background())
	// go func() {
	// 	exit := make(chan os.Signal, 1)
	// 	signal.Notify(exit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) // подписываем канал на чтение сист. сигнала
	// 	<-exit
	// 	cancel() // вызов фугкции отмены
	// }()

	stop := make(chan struct{})

	time.Sleep(time.Second)
	go RunTaskEverySecond(ctx, stop) // если вынести функцию отделно, а потом
	//вызвать горутиной, то горутины синхронизируются (Channel Synchronization)
	// даем поработать алгоритму
	time.Sleep(5 * time.Second) //без этого гоурутина не успевает срабоать
	close(stop)

}
