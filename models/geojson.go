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

func ResponseToPointGeojson(response ResponseData) (Feature, error) {
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

type Feature2 struct {
	Type       string                 `json:"type"`
	Geometry   Geometry2              `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

type Geometry2 struct {
	Type        string       `json:"type"`
	Coordinates [][2]float64 `json:"coordinates"`
}

func LineToSessionIdGeojson(response ResponseForLine) (Feature2, error) {
	return Feature2{Type: "Feature",
		Geometry: Geometry2{Type: "LineString",
			Coordinates: response.Coordinates},
		Properties: map[string]interface{}{
			"Login":      response.Login,
			"Session_id": response.Session_id,
			"Start_time": response.Start_time,
			"End_time":   response.End_time,
		}}, nil
}

// func LineToSessionIdGыыeojson(response ResponseForLine) (map[string]interface{}, error) {
// 	return map[string]interface{}{
// 		"Type": "Feature",
// 		"Geometry": map[string]interface{}{
// 			"Geometry":    "LineString",
// 			"Coordinates": response.Coordinates,
// 		},
// 		"Properties": map[string]interface{}{
// 			"Login":      response.Login,
// 			"Session_id": response.Session_id,
// 			"Start_time": response.Start_time,
// 			"End_time":   response.End_time,
// 		}}, nil
// }
