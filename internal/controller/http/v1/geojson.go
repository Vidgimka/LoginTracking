package v1

import (
	"time"

	"github.com/Vidgimka/LoginTracking/internal/domain"
)

func ResponseToPointGeojson(response domain.PointData) (Feature, error) {
	return Feature{Type: "Feature",
		Geometry: Geometry{Type: "Point",
			Coordinates: response.Coordinates},
		Properties: map[string]interface{}{
			"Login":      response.Login,
			"Session_id": response.SessionId,
			"CreatedA":   response.CreatedAt,
		}}, nil
}

func ResponseToCoord(response domain.PointData) (domain.Coord, error) {
	return response.Coordinates, nil
}

type Feature2 struct {
	Type       string     `json:"type"`
	Geometry   Geometry2  `json:"geometry"`
	Properties Properties `json:"properties"`
}

type Geometry2 struct {
	Type        string        `json:"type"`
	Coordinates domain.Coords `json:"coordinates"`
}

type Properties struct {
	Login     string
	SessionId int
	CreatedAt time.Time
}

func LineToSessionIdGeojson(response domain.LineData) (Feature2, error) {
	return Feature2{Type: "Feature",
		Geometry: Geometry2{Type: "LineString",
			Coordinates: response.Coordinates},
		Properties: Properties{Login: response.Login,
			SessionId: response.SessionId,
			CreatedAt: response.CreatedAt,
		}}, nil
}
