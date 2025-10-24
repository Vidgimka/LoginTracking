package v1

const (
	featureCollection string = "FeatureCollection"
	feature           string = "Feature"
	multiPoint        string = "MultiPoint"
	lineString        string = "LineString"
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
	Type        string       `json:"type"`
	Coordinates [][2]float64 `json:"coordinates"`
}

type Properties struct {
	Login string `json:"login"`
}
