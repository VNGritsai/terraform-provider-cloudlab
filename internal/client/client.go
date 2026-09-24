package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var ErrNotFound = errors.New("resource not found")

// Instance represents a cloudlab instance.
type Instance struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Cpu    int    `json:"cpu"`
	Memory int    `json:"memory"`
}

// Client represents a client for interacting with the cloudlab API.
type Client struct {
	endpoint   string
	httpClient *http.Client
}

// NewClient creates a new client for interacting with the cloudlab API.
func NewClient(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doRequest sends a request to the cloudlab API and returns the response body.
func (c *Client) doRequest(ctx context.Context, method, path string, body any) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.endpoint+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	switch {
	case resp.StatusCode == http.StatusNotFound:
		return nil, ErrNotFound
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		return data, nil
	default:
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(data))
	}
}

// CreateInstance creates a new instance in cloudlab.
func (c *Client) CreateInstance(ctx context.Context, instance Instance) (Instance, error) {
	data, err := c.doRequest(ctx, http.MethodPost, "/api/v1/instances", instance)
	if err != nil {
		return Instance{}, err
	}

	var created Instance
	if err := json.Unmarshal(data, &created); err != nil {
		return Instance{}, fmt.Errorf("unmarshal created instance: %w", err)
	}
	return created, nil
}

// GetInstance retrieves an instance by its ID.
func (c *Client) GetInstance(ctx context.Context, id int) (Instance, error) {
	data, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/instances/%d", id), nil)
	if err != nil {
		return Instance{}, err
	}

	var instance Instance
	if err := json.Unmarshal(data, &instance); err != nil {
		return Instance{}, fmt.Errorf("unmarshal instance: %w", err)
	}
	return instance, nil
}

// UpdateInstance updates an existing instance with new data.
func (c *Client) UpdateInstance(ctx context.Context, instance Instance) (Instance, error) {
	data, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/api/v1/instances/%d", instance.Id), instance)
	if err != nil {
		return Instance{}, err
	}

	var updated Instance
	if err := json.Unmarshal(data, &updated); err != nil {
		return Instance{}, fmt.Errorf("unmarshal updated instance: %w", err)
	}
	return updated, nil
}

// DeleteInstance deletes an instance by its ID.
func (c *Client) DeleteInstance(ctx context.Context, id int) error {
	_, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/instances/%d", id), nil)
	return err
}
