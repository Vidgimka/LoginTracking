package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Vidgimka/LoginTracking/internal/domain"
	"github.com/gin-gonic/gin"
)

type servicePointData interface {
	BuildPointsByDate(ctx context.Context, login string, start, end time.Time) ([]domain.PointData, error)
}

type domainPointData interface {
	BuildLines(ctx context.Context, points []domain.PointData) ([]domain.LineData, error)
}

type handler struct {
	service servicePointData
	domain  domainPointData
}

func NewHandlers(service servicePointData) *handler {
	return &handler{
		service: service,
	}
}

func parseInputData(start, end string) (time.Time, time.Time, error) {
	startInTime, err := time.Parse("2/1/2006", start)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("time.Parse.startInTime: %w", err)
	}
	endInTime, err := time.Parse("2/1/2006", end)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("time.Parse.endInTime: %w", err)
	}
	return startInTime, endInTime, nil
}

func (h *handler) GetPositionByDate(c *gin.Context) {
	login := c.Param("login")
	if login == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "login cannot be empty"})
		return
	}
	start := c.Query("start")
	if start == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start time cannot be empty"})
		return
	}
	end := c.Query("end")
	if end == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end time cannot be empty"})
		return
	}
	startInTime, endInTime, err := parseInputData(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse date"})
		return
	}
	visual := c.Query("visual")
	if visual == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data type cannot be empty"})
		return
	}

	points, err := h.service.BuildPointsByDate(c, login, startInTime, endInTime)
	if err != nil {
		log.Println("errors data points not received")
		return
	}
	pointToFeature, err := PointToFeature(login, points)
	if err != nil {
		log.Println("error preparing data for geojson")
		return
	}

	switch visual {
	case "point":
		geoJsonData, err := NewGeoJsonMessageV3(pointToFeature)
		if err != nil {
			log.Println("geojson generation error")
			return
		}
		geojson, err := json.MarshalIndent(geoJsonData, "", " ")
		if err != nil {
			log.Println("error preparing data for geojson")
			return
		}
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(geojson)
	case "line":
		lines, err := h.domain.BuildLines(c, points)
		if err != nil {
			log.Println("line building not completed")
			return
		}
		lineToFeature, err := LinesToFeature(login, lines)
		geoJsonData, err := NewGeoJsonMessageV3(pointToFeature, lineToFeature)
		if err != nil {
			log.Println("geojson generation error")
			return
		}
		geojson, err := json.MarshalIndent(geoJsonData, "", " ")
		if err != nil {
			log.Println("error preparing data for geojson")
			return
		}
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(geojson)
		//http://localhost:8080/loginytracking/v1/logins/tsb645/date?start=....&end=....&visual=...
		//2025-06-29T17:33:54.253593+03:00
	}
}
