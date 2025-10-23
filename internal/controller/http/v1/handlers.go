package v1

import (
	"context"
	"fmt"
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

func (h *handler) GetLinesByDate(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "data points not received"})
		return
	}

	switch visual {
	case "point":

	case "line":

	}

	// if visual == "line" {
	// 	lines, err := h.domain.BuildLines(c, points)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "line building not completed"})
	// 		return
	// 	} else if visual == "point" {

	// 	}
	// }

	//http://localhost:8080/loginytracking/v1/logins/tsb645/date?start=....&end=....&visual=...
	//2025-06-29T17:33:54.253593+03:00
}

// func (s *handler) GetlineCollection(c *gin.Context) {
// 	login := c.Param("login")
// 	if login == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "login cannot be empty"})
// 		return
// 	}
// 	start := c.Query("start")
// 	end := c.Query("end")
// 	if start == "" || end == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "login cannot be empty"})
// 		return
// 	}
// 	startTimeFormat, err := time.Parse("2006-01-02", start)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "start format date is not 2006-01-02"})
// 		return
// 	}
// 	endTimeFormat, err := time.Parse("2006-01-02", end)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "end format date is not 2006-01-02"})
// 		return
// 	}
// 	////////
// 	c.Header("Content-Type", "application/json")
// 	c.Writer.Write([]byte(`{"type": "FeatureCollection","features": [`))
// 	if err := repo.repo.GetLines(c.Request.Context(), c.Writer, login, startTimeFormat, endTimeFormat); err != nil {
// 		log.Printf("get login:%v", err)
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"db error": "login data not recirved"})
// 		return
// 	}
// 	///////////////////////////////////////////////////////////////
// 	firstElem := true
// 	responseToGeojson, err := models.LineToSessionIdGeojson(response)
// 	if err != nil {
// 		fmt.Printf("Conver to geojson error %v", err)
// 		continue
// 	}
// 	inJson, err := json.Marshal(responseToGeojson)
// 	if err != nil {
// 		fmt.Printf("Serializationerror %v", err)
// 		continue
// 	}
// 	if !firstElem {
// 		write.Write([]byte(","))
// 	}
// 	firstElem = false
// 	write.Write(inJson)
// 	///////////////////////////////////////////////////////////////
// 	//////////
// 	c.Writer.Write([]byte("]}"))
// 	c.Writer.Flush()
// }
// func (s *handler) GetUserByLoginAnDatetime(c *gin.Context) {
// 	c.Header("Content-Type", "application/json")
// 	login := c.Param("login")
// 	if login == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Login cannot be empty"})
// 		return
// 	}
// 	CreatedAt := c.Param("datetime")
// 	if CreatedAt == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "datetime cannot be empty"})
// 		return
// 	}
// 	if err := repo.repo.GetByDatetime(c.Request.Context(), c.Writer, login, CreatedAt); err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"db error": "login and datetime data not recirved"})
// 	}
// 	c.Writer.Flush()
// }
// func (s *handler) GetAllUsers(c *gin.Context) {
// 	c.Header("Content-Type", "application/json")
// 	if err := repo.repo.GetByAllUser(c.Request.Context(), c.Writer); err != nil {
// 		log.Printf("get all users:%v", err)
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"db error": "users not received"})
// 		return
// 	}
// 	c.Writer.Flush()
// }
// func (s *handler) GetUserByLogin(c *gin.Context) {
// 	c.Header("Content-Type", "application/json")
// 	login := c.Param("login")
// 	if login == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "login cannot be empty"})
// 		return
// 	}
// 	if err := repo.repo.GetByLogin(c.Request.Context(), c.Writer, login); err != nil {
// 		log.Printf("get login:%v", err)
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"db error": "login data not recirved"})
// 		return
// 	}
// 	c.Writer.Flush()
// }
// func (s *handler) GetUserByLoginAndSessionId(c *gin.Context) {
// 	c.Header("Content-Type", "application/json")
// 	login := c.Param("login")
// 	if login == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "login cannot be empty"})
// 		return
// 	}
// 	Session_id := c.Param("session_id")
// 	if Session_id == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id cannot be empty"})
// 		return
// 	}
// 	if err := repo.repo.GetBySessionId(c.Request.Context(), c.Writer, login, Session_id); err != nil {
// 		log.Printf("get by session_id: %v", err)
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"db error": "login and session_id data not recirved"})
// 	}
// 	c.Writer.Flush()
// }
