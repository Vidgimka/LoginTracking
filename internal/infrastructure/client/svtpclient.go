package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const defaultTimeout = 10 * time.Second

type client struct {
	client  *http.Client
	baseUrl string
}

func New(httpClient *http.Client, baseUrl string) *client {

	if baseUrl == "" {
		baseUrl = os.Getenv("URL")
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}

	return &client{
		client:  httpClient,
		baseUrl: baseUrl,
	}

}

func (c *client) GetUsersOnline(ctx context.Context) ([]Data, error) {

	var usersOnline GetUsersOnlineResponse
	resp, err := http.Get(c.baseUrl)
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
