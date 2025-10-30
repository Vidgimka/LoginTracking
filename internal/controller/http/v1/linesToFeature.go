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
			Login: login,
		},
	}, nil
}

func LinesToFeature(login string, lines []domain.LineData) ([]Feature, error) {

	result := make([]Feature, 0, len(lines))
	for _, line := range lines {

		lineIoFeature := Feature{
			Type: feature,
			Geometry: Geometry{
				Type:        lineString,
				Coordinates: line.Coordinates.SliceCoords(),
			},
			Properties: Properties{
				Login:     login,
				SessionId: line.SessionId,
			},
		}
		result = append(result, lineIoFeature)
	}

	return result, nil
}
