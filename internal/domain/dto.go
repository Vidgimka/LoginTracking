package domain

import "time"

type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Type       string                 `json:"type"`
	Geometry   Geometry               `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

type Geometry struct {
	Type        string `json:"type"`
	Coordinates Coord  `json:"coordinates"`
}

type Properties struct {
	Login     string
	SessionId int
	StartTime time.Time
	EndTime   time.Time
}
