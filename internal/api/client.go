package api

type Client struct {
	BaseURL string
	Token   string
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
	}
}
