package v1

import "github.com/Vidgimka/LoginTracking/internal/domain"

const (
	featureCollection string = "FeatureCollection"
	feature           string = "Feature"
	MultiPoint        string = "MultiPoint"
	LineString        string = "LineString"
)

type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Type       string     `json:"type"`
	Geometry   Geometry   `json:"geometry"`
	Properties Properties `json:"properties"`
}

type Geometry struct {
	Type        string        `json:"type"`
	Coordinates domain.Coords `json:"coordinates"`
}

type Properties struct {
	Login     string `json:"login"`
	SessionId int    `json:"session_id"`
}
