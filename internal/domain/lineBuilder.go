package domain

import (
	"context"
	"fmt"
)

func BuildLines(ctx context.Context, points []PointData) ([]LineData, error) {

	// проверить конткст на ошибкуct
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("ctx.Err: %w", err)
	}

	mapWithGroupsBySessionId := make(map[int]*LineData, len(points))

	for _, point := range points {

		if groupWithSessionId, exists := mapWithGroupsBySessionId[point.SessionId]; !exists {
			mapWithGroupsBySessionId[point.SessionId] = &LineData{
				Login:       point.Login,
				SessionId:   point.SessionId,
				Coordinates: Coords{point.Coordinates},
				CreatedAt:   point.CreatedAt}
		} else {
			groupWithSessionId.Coordinates = append(groupWithSessionId.Coordinates, point.Coordinates)
		}
	}

	var result []LineData

	for _, group := range mapWithGroupsBySessionId {

		result = append(result, *group)
	}

	return result, nil
}
