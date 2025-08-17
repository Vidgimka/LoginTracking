package client

import (
	"time"

	"github.com/Vidgimka/LoginTracking/internal/domain"
)

type GetUsersOnlineResponse struct {
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

func (d Data) ToEntity() domain.Data {
	var entity domain.Data
	entity.Login = d.Login
	entity.SessionId = d.SessionId
	entity.Subnet = d.Subnet
	entity.Mountpoint = d.Mountpoint
	entity.Station = d.Station
	entity.NtripAgent = d.NtripAgent
	entity.ConnectTime = d.ConnectTime
	entity.TimeSpan = d.TimeSpan
	entity.RecievedData = d.RecievedData
	entity.SentData = d.SentData
	entity.StatusCode = d.StatusCode
	entity.Latency = d.Latency
	entity.SvNum = d.SvNum
	entity.Coordinaties.Lat = d.Lat
	entity.Coordinaties.Lon = d.Lon
	entity.Height = d.Height
	entity.StationDistance = d.StationDistance
	entity.CreatedAt = d.CreatedAt
	return entity
}

func (r GetUsersOnlineResponse) ToEntitys() []domain.Data {
	var entitys []domain.Data
	for _, data := range r.Data {
		entitys = append(entitys, data.ToEntity())
	}
	return entitys
}
