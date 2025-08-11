package domain

import "time"

type Coord struct {
	Lat float64
	Lon float64
}

type Coords []Coord

func (c *Coord) SliceCoord() [2]float64 {
	coordinaties := [2]float64{c.Lat, c.Lon}
	return coordinaties
}

func (c *Coords) SliceCoords() [][2]float64 {
	coordinates := make([][2]float64, 0, len(*c))
	for _, coord := range *c {
		coordinates = append(coordinates, coord.SliceCoord())
	}
	return coordinates
}

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
	Coordinaties    Coord
	Height          float64
	StationDistance float64
	CreatedAt       time.Time
}

type PointData struct {
	Login            string
	Session_id       int
	Coordinates      Coord
	Station_distance float64
	CreatedAt        time.Time
}

type LineData struct {
	Login       string
	Session_id  int
	Coordinates Coords
	Start_time  time.Time
	End_time    time.Time
}
