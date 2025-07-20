package models

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
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

func ResponseToGeojson(response ResponseData) (Feature, error) {
	return Feature{Type: "Feature",
		Geometry: Geometry{Type: "Point",
			Coordinates: []float64{response.Lon, response.Lat}},
		Properties: map[string]interface{}{
			"Login":      response.Login,
			"Session_id": response.Session_id,
			"CreatedA":   response.CreatedAt,
		}}, nil
}

func ResponseToCoord(response ResponseData) ([]float64, error) {
	return []float64{response.Lon, response.Lat}, nil
}
