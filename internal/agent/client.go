package agent

import (
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

func (c *Client) SendMetric(metric model.Metric) error {}
