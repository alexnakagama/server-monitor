package monitor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alexnakagama/server-monitor/internal/model"
)

type Client struct {
	baseURL string
	token   string
	client  *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (c *Client) Login(ctx context.Context, username string, password string) error {
	reqBody := loginRequest{
		Username: username,
		Password: password,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	url := c.baseURL + "/users/login"

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status code: %d", resp.StatusCode)
	}

	var loginResp loginResponse

	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	if err != nil {
		return err
	}

	c.token = loginResp.Token

	return nil
}

func (c *Client) GetServers(ctx context.Context) ([]model.Server, error) {
	url := c.baseURL + "/servers"

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get servers failed with status code: %d", resp.StatusCode)
	}

	var servers []model.Server

	err = json.NewDecoder(resp.Body).Decode(&servers)
	if err != nil {
		return nil, err
	}

	return servers, nil
}

func (c *Client) GetServerMetrics(ctx context.Context, serverID int) ([]model.Metric, error) {
	url := fmt.Sprintf("%s/metrics/server/%d", c.baseURL, serverID)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"get server metrics failed with status code: %d",
			resp.StatusCode,
		)
	}

	var metrics []model.Metric

	err = json.NewDecoder(resp.Body).Decode(&metrics)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}
