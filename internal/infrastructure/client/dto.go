package client

import (
	"time"

	"github.com/Vidgimka/LoginTracking/internal/domain"
)

type getUsersOnlineResponse struct {
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

func (d Data) toEntity() domain.Data {
	return domain.Data{
		Login:        d.Login,
		SessionId:    d.SessionId,
		Subnet:       d.Subnet,
		Mountpoint:   d.Mountpoint,
		Station:      d.Station,
		NtripAgent:   d.NtripAgent,
		ConnectTime:  d.ConnectTime,
		TimeSpan:     d.TimeSpan,
		RecievedData: d.RecievedData,
		SentData:     d.SentData,
		StatusCode:   d.StatusCode,
		Latency:      d.Latency,
		SvNum:        d.SvNum,
		Coordinaties: domain.Coord{
			Lat: d.Lat,
			Lon: d.Lon,
		},
		Height:          d.Height,
		StationDistance: d.StationDistance,
		CreatedAt:       d.CreatedAt,
	}
}

func (r getUsersOnlineResponse) ToEntities() []domain.Data {
	entities := make([]domain.Data, 0, len(r.Data))
	for _, data := range r.Data {
		entities = append(entities, data.toEntity())
	}
	return entities
}
