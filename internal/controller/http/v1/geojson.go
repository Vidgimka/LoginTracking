package v1

import (
	"github.com/Vidgimka/LoginTracking/internal/domain"
)


func NewGeoJsonMessege(pointsData []domain.PointData) (FeatureCollection, error) {

	sliceCoordPoint := make([][2]float64,0,len(pointsData))
	for  _, point := range pointsData{
		
	}


	return FeatureCollection{
		Type: feature,
		Features: []Feature{Feature{
			Type: feature, 
			Geometry: Geometry{
				Type: MultiPoint, 
				Coordinates:, }},
	}
}
}

// func NewGeoJsonMessege(response domain.PointData) (Feature2, error) {
// 	return Feature2{Type: "Feature",
// 		Geometry: Geometry{Type: "Point",
// 			Coordinates: response.Coordinates},
// 		Properties: map[string]interface{}{
// 			"Login":      response.Login,
// 			"Session_id": response.SessionId,
// 			"CreatedA":   response.CreatedAt,
// 		}}, nil
// }

func ResponseToCoord(response domain.PointData) (domain.Coord, error) {
	return response.Coordinates, nil
}
