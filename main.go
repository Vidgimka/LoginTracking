package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
var DB *gorm.DB

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

	//var UsersOnline GeoData // записываем в переменную UsersOnline  данные из тела ответа
	if err := json.Unmarshal(d, &UsersOnline); err != nil {
		log.Fatal(err.Error())
	}
	//fmt.Println(UsersOnline)
	return UsersOnline
}

func Init() *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=OnlineUsersIist port=5432 sslmode=disable"
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	//gorm.Open("postgres", "user=postgres password=postgres dbname=SvtpPrin sslmode=disable TimeZone=Asia/Shanghai")

	if err != nil {
		fmt.Println("не подключилось к БД")
	}
	DB.AutoMigrate(&Data{})
	return DB
}

var Dbase *gorm.DB

func GetDB() *gorm.DB { // проверяем подкюченеие к базе данных
	if Dbase == nil {
		Dbase = Init() // проинициализировали бд через Init и присвоили в переменную dbase, так как инициилизация return db
		var sleep = time.Duration(1)
		for Dbase == nil { // ждем 3 секунды если не подключается
			sleep = sleep * 3
			fmt.Printf("База данных не доступна. Пожождите %d секунды.\n", sleep)
			time.Sleep(sleep * time.Second)
			Dbase = Init() // еще раз кладем базу данных в переменную dbase
		}
	}
	return Dbase
}

func CreateEntry() {
	DB.Create(&Data{})
}

func main() {
	ReadFileData()
	//fmt.Println(UsersOnline.Data)
	Init()
	//CreateEntry()
	GetDB()
	//Dbase.Create(&Data{})

}
