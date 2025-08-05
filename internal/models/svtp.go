package models

import "time"

type GeoData struct {
	Status  string
	Message string
	Data    []Data
}

type Data struct {
	Login           string
	SessionId       int
	Subnet          string
	Mountpoint      string
	Station         string
	NtripAgent      string
	ConnectTime     int
	TimeSpan        int
	RecievedData    float64
	SentData        float64
	StatusCode      int
	Latency         int
	SvNum           int
	Lat             float64
	Lon             float64
	Height          float64
	StationDistance float64
	CreatedAt       time.Time
}

type PointData struct {
	Login            string    `json:"login"`
	Session_id       int       `json:"session_id"`
	Lat              float64   `json:"lat"`
	Lon              float64   `json:"lon"`
	Station_distance float64   `json:"station_distance"`
	CreatedAt        time.Time `json:"сreated_at"`
}

type LineBuilder struct {
	Login       string       `json:"login"`
	Session_id  int          `json:"session_id"`
	Coordinates [][2]float64 `json:"coordinates"`
	Start_time  time.Time    `json:"start_time"`
	End_time    time.Time    `json:"end_time"`
}
