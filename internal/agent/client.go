package agent

import "net/http"

type Client struct {
	baseURL string
	token   string
	client  *http.Client
}
