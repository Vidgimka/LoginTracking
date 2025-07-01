package handlers

import (
	"time"

	"github.com/Vidgimka/LoginTracking.git/models"
	"github.com/Vidgimka/LoginTracking.git/repository"
	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Login            string    `json:"login"`
	Session_id       int       `json:"session_id"`
	Lat              float64   `json:"lat"`
	Lon              float64   `json:"lon"`
	Station_distance float64   `json:"station_distance"`
	CreatedAt        time.Time `json:"сreated_at"`
}

func responseToGeojson(response ResponseData) (models.Feature, error) {
	return models.Feature{Type: "Feature",
		Geometry: models.Geometry{Type: "Point",
			Coordinates: []float64{response.Lon, response.Lat}},
		Properties: map[string]interface{}{
			"Login":      response.Login,
			"Session_id": response.Session_id,
			"CreatedA":   response.CreatedAt,
		}}, nil
}

func responseToCoord(response ResponseData) ([]float64, error) {
	return []float64{response.Lon, response.Lat}, nil
}

type handlers struct {
	db repository.PostgresGormRepoInterfase
}

func (repo *handlers) GetAllUsers(c *gin.Context) error {
	c.Header("Content-Type", "application/json")

	if err := repo.db.GetByAllUser(c.Request.Context(), c.Writer); err != nil {
		return err // обработать
	}

	//  переделать в интерфейс БД
	// rows, err := db.Raw("SELECT * FROM data").Rows()
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error:": "DB error"})
	// }
	// defer rows.Close()
	// inFirst := true
	// for rows.Next() {
	// 	var fD models.Data
	// 	if err := rows.Scan(&fD.Login, &fD.SessionId, &fD.Subnet, &fD.Mountpoint, &fD.Station, &fD.NtripAgent, &fD.ConnectTime,
	// 		&fD.TimeSpan, &fD.RecievedData, &fD.SentData, &fD.StatusCode, &fD.Latency, &fD.SvNum, &fD.Lat, &fD.Lon, &fD.Height,
	// 		&fD.StationDistance, &fD.CreatedAt); err != nil {
	// 		fmt.Printf("Scan error: %v", err)
	// 		continue
	// 	}
	// 	if !inFirst {
	// 		c.Writer.Write([]byte(","))
	// 	}
	// 	inJson, err := json.Marshal(fD)
	// 	if err != nil {
	// 		fmt.Printf("Serializationerror %v", err)
	// 		continue
	// 	}
	// 	inFirst = false
	// 	c.Writer.Write(inJson)
	// 	c.Writer.Flush()
	// }
	// c.Writer.Write([]byte("]"))

	return nil
}

// http://localhost:8080/UsersOnline2
