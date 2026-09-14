package monitor

type Client struct {
	baseURL string
	token   string
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
	}
}
