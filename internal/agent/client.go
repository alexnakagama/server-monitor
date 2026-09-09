package agent

import (
	"bytes"
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

func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		client:  &http.Client{},
	}
}

func (c *Client) SendMetric(metric model.Metric) error {
	body, err := json.Marshal(metric)
	if err != nil {
		return err
	}

	url := c.baseURL + "/agents/metrics"

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
