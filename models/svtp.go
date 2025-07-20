package models

import "time"

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

type ResponseData struct {
	Login            string    `json:"login"`
	Session_id       int       `json:"session_id"`
	Lat              float64   `json:"lat"`
	Lon              float64   `json:"lon"`
	Station_distance float64   `json:"station_distance"`
	CreatedAt        time.Time `json:"сreated_at"`
}
