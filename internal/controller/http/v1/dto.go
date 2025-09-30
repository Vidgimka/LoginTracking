package v1

import "github.com/Vidgimka/LoginTracking/internal/domain"

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
	Type        string       `json:"type"`
	Coordinates domain.Coord `json:"coordinates"`
}
