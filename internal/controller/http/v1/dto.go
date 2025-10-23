package v1

import "github.com/Vidgimka/LoginTracking/internal/domain"

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

func PointDataBuilder(login string, pointsData []domain.PointData) (FeatureCollection, error) {

	sliceCoordPoint := make([][2]float64, 0, len(pointsData))

	for _, point := range pointsData {
		sliceCoordPoint = append(sliceCoordPoint, point.Coordinates.SliceCoord())
	}

	return FeatureCollection{
		Type: feature,
		Features: []Feature{Feature{
			Type: feature,
			Geometry: Geometry{
				Type:        multiPoint,
				Coordinates: sliceCoordPoint},
			Properties: Properties{
				Login: login,
			}},
		},
	}, nil
}
