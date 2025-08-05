package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Vidgimka/LoginTracking.git/internal/models"
)

// Чтобы внешний код зависел от интерфейса, а не от конкретной структуры.
type HttpClientInterface interface {
	ReadDataFromAPI() ([]models.Data, error)
}

type ApiClient struct {
	client *http.Client
}

func NewHttpClient() HttpClientInterface {
	return &ApiClient{
		client: &http.Client{},
	}
}

func (c *ApiClient) ReadDataFromAPI() ([]models.Data, error) {
	url := os.Getenv("URL")
	var usersOnline models.GeoData
	resp, err := http.Get(url)
	if err != nil {
		return []models.Data{}, fmt.Errorf("HTTP request error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return []models.Data{}, fmt.Errorf("sratus code error: %d", resp.StatusCode)
	}
	d, err := io.ReadAll(resp.Body)
	if err != nil {
		return []models.Data{}, fmt.Errorf("read error %w", err)
	}
	if err := json.Unmarshal(d, &usersOnline); err != nil {
		return []models.Data{}, fmt.Errorf("unmarshal error %w", err)
	}
	return usersOnline.Data, nil
}
