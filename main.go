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
// var UsersOnlineD Data

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

func Init() *gorm.DB { // логика, как при вызове этой функци будет создаваться БД
	var DB *gorm.DB
	dsn := "host=localhost user=postgres password=postgres dbname=OnlineUsersIist port=5432 sslmode=disable"
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("не подключилось к БД")
	}
	DB.AutoMigrate(&Data{})
	return DB
}

/* func GetDB() *gorm.DB { // проверяем подкюченеие к базе данных
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
// }*/

func RunTaskEverySecond(ctx context.Context, stop <-chan struct{}) { // совмещаем логику 2х функций, Init() созданию БД
	// и записи ReadFileData().Data в переменную Data конкретного куска данных полученных Data с API
	var Dbase *gorm.DB = Init()
	ticker1 := time.NewTicker(time.Second)
	defer ticker1.Stop()
	for {
		select {
		case <-ticker1.C:
			fmt.Println("Running task every second")
			Data := ReadFileData().Data // помещаем в переменную вычетанные данные DATA
			Dbase.Create(&Data)         // запись в БД
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

// функция, которая будет ходить в бд и возвращать данные,а не писать их в глобальную переменнюу
func GetUsersOnline(dbase *gorm.DB) []Data {
	var usersOnline2 []Data // записываем в переменную срез элементов Data usersOnline2  данные из базы данных
	//var dbase *gorm.DB = Init() // обращаемся к бд
	fmt.Println("Обращение в базу данных")
	dbase.Find(&usersOnline2) // записываем данные из БД в переменную
	return usersOnline2
}

// func Getlogin(c *gin.Context) {
// 	//var Login string
// 	Login := c.Param("Login") // записываем в переменную вычитанный логин из URL
// 	// запускаем цикл, который переберет все значения из переменной базы данных
// 	// и при собпадении всех записей с заданным логином выведет их в теле ответа
// 	for _, a := range usersOnline2 {
// 		if a.Login == Login {
// 			c.JSON(http.StatusOK, a)
// 			return
// 		}
// 	}
// 	c.JSON(http.StatusNotFound, gin.H{"message": "Login not found"})
// }

// func Testlogin() {
// 	var login string
// 	fmt.Fscan(os.Stdin, &login)
// 	for _, a := range usersOnline2 {
// 		if a.Login == login {
// 			fmt.Println(a)
// 			return
// 		}
// 	}
// }

func main() {
	// реализация в основном пототке graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	stop := make(chan struct{})
	time.Sleep(time.Second)
	go RunTaskEverySecond(ctx, stop) // если вынести функцию отделно, а потом
	//вызвать горутиной, то горутины синхронизируются (Channel Synchronization)
	// даем поработать алгоритму
	time.Sleep(5 * time.Second) //без этого гоурутина не успевает срабоать
	close(stop)                 // закрывает горутину main

	var dbase *gorm.DB = Init() // запрос делаетя один раз в main, и далее везде используется для запросов
	// Блок с сервером

	router := gin.Default()

	//var dbase *gorm.DB = Init() // обращаемся к бд
	//loginOnline := GetUsersOnline(dbase)

	router.GET("/UsersOnline2", func(c *gin.Context) {

		// делаем функцию, которая будет ходить в бд и возвращать данные,а не писать их в глобальную переменнюу
		//var Online2 []Data // записываем в переменную срез элементов Data usersOnline2  данные из базы данных
		loginOnline := GetUsersOnline(dbase)
		//var Dbase *gorm.DB = Init() // обращаемся к бд
		fmt.Println("Данные из БД полученны")
		//dbase.Find(&usersOnline2) // записываем данные из БД в переменную
		//fmt.Println(usersOnline2)
		c.JSON(http.StatusOK, loginOnline) //c.IndentedJSON(http.StatusOK, usersOnline2)
	})
	//curl http://localhost:8080/UsersOnline2

	router.GET("/UsersOnline2/:login", func(c *gin.Context) {
		loginOnline := GetUsersOnline(dbase)
		fmt.Println("Данные из БД полученны")
		dbase.Find(&loginOnline)
		login := c.Param("login") // записываем в переменную вычитанный логин из URL
		// запускаем цикл, который переберет все значения из переменной базы данных
		// и при собпадении всех записей с заданным логином выведет их в теле ответа
		for _, a := range loginOnline {
			if a.Login == login {
				c.JSON(http.StatusOK, a)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"message": "Login not found"})
	})
	// //curl http://localhost:8080/UsersOnline2/login

	router.Run("localhost:8080")

}
