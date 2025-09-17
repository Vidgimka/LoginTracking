package v1

import "time"

func ResponseToPointGeojson(response PointData) (Feature, error) {
	return Feature{Type: "Feature",
		Geometry: Geometry{Type: "Point",
			Coordinates: response.Coordinates},
		Properties: map[string]interface{}{
			"Login":      response.Login,
			"Session_id": response.SessionId,
			"CreatedA":   response.CreatedAt,
		}}, nil
}

func ResponseToCoord(response PointData) (Coord, error) {
	return response.Coordinates, nil
}

type Feature2 struct {
	Type       string     `json:"type"`
	Geometry   Geometry2  `json:"geometry"`
	Properties Properties `json:"properties"`
}

type Geometry2 struct {
	Type        string `json:"type"`
	Coordinates Coords `json:"coordinates"`
}

type Properties struct {
	Login     string
	SessionId int
	StartTime time.Time
	EndTime   time.Time
}

func LineToSessionIdGeojson(response LineData) (Feature2, error) {
	return Feature2{Type: "Feature",
		Geometry: Geometry2{Type: "LineString",
			Coordinates: response.Coordinates},
		Properties: Properties{Login: response.Login,
			SessionId: response.SessionId,
			StartTime: response.StartTime,
			EndTime:   response.EndTime,
		}}, nil
}
