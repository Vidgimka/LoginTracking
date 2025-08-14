package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type client struct {
	client  *http.Client
	baseUrl string
}

func New(httpClient *http.Client, baseUrl string) *client {

	// url := os.Getenv("URL")

	return &client{
		client:  httpClient,
		baseUrl: baseUrl,
	}

}

func (c *client) GetUsersOnline(ctx context.Context, url string) ([]Data, error) {

	var usersOnline GetUsersOnlineResponse
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("sratus code error: %d", resp.StatusCode)
	}
	d, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read error %w", err)
	}
	if err := json.Unmarshal(d, &usersOnline); err != nil {
		return nil, fmt.Errorf("unmarshal error %w", err)
	}
	return usersOnline.Data, nil
}
