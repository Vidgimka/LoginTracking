package domain

import "time"

type Coord struct {
	Lat float64
	Lon float64
}

func (c Coord) SliceCoord() [2]float64 {
	coordinaties := [2]float64{c.Lat, c.Lon}
	return coordinaties
}

type Coords []Coord

func (c Coords) SliceCoords() [][2]float64 {
	coordinates := make([][2]float64, 0, len(c))
	for _, coord := range c {
		coordinates = append(coordinates, coord.SliceCoord())
	}
	return coordinates
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
	Login       string
	SessionId   int
	Coordinates Coord
	CreatedAt   time.Time
}

type LineData struct {
	Login       string
	SessionId   int
	Coordinates Coords
	CreatedAt   time.Time
}
