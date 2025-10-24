package v1

import "github.com/Vidgimka/LoginTracking/internal/domain"

func PointToFeature(login string, pointsData []domain.PointData) (Feature, error) {

	sliceCoordPoint := make([][2]float64, 0, len(pointsData))

	for _, point := range pointsData {
		sliceCoordPoint = append(sliceCoordPoint, point.Coordinates.SliceCoord())
	}

	return Feature{
		Type: feature,
		Geometry: Geometry{
			Type:        multiPoint,
			Coordinates: sliceCoordPoint,
		},
		Properties: Properties{
			Login: lineString,
		},
	}, nil
}

func LinesToFeature(login string, lines []domain.LineData) (Feature, error) {

	
	for _,line := range lines{

	}

	return Feature{
		Type: feature,
		Geometry: Geometry{
			Type:        multiPoint,
			Coordinates: ,
		},
		Properties: login,
	}, nil
}
