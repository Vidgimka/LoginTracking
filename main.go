package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/**/

type GeoData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    []struct {
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
	} `json:"data"`
}

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

	var UsersOnline GeoData // записываем в переменную UsersOnline  данные из тела ответа
	if err := json.Unmarshal(d, &UsersOnline); err != nil {
		//panic(err)
		log.Fatal(err.Error())
	}
	//fmt.Println(UsersOnline)
	return UsersOnline

}

var DB *gorm.DB

func Init() *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=OnlineUsersIist port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	//gorm.Open("postgres", "user=postgres password=postgres dbname=SvtpPrin sslmode=disable TimeZone=Asia/Shanghai")

	if err != nil {
		fmt.Println("не подключилось к БД")
	}

	db.AutoMigrate(&GeoData{})
	return db
}

func main() {
	//var UsOn GeoData
	//UsOn = ReadFileData()
	Init()
	//fmt.Println(UsOn)
}
