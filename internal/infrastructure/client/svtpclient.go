package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Vidgimka/LoginTracking/internal/domain"
)

const defaultTimeout = 10 * time.Second

type client struct {
	client  *http.Client
	baseUrl string
}

func New(httpClient *http.Client, baseUrl string) *client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}
	return &client{
		client:  httpClient,
		baseUrl: baseUrl,
	}
}

func NewUrlEndpoint(baseUrl string) (string, error) {
	path := os.Getenv("QUERY")
	return fmt.Sprintf("%s%s", baseUrl, path), nil
}

func (c *client) GetUsersOnline(ctx context.Context) ([]domain.Data, error) {
	var usersOnline getUsersOnlineResponse
	url, err := NewUrlEndpoint(c.baseUrl)
	if err != nil {
		return nil, fmt.Errorf("url build error: %w", err)
	}
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
	return usersOnline.ToEntities(), nil
}
