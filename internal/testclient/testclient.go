package testclient

import (
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	httpClient http.Client
	URL        *url.URL
}

func NewClient(URL *url.URL) Client {
	httpClient := http.Client{}
	return Client{httpClient: httpClient}
}

func (c *Client) NewRequest(method string, url url.URL, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

type Data struct {
	Memo string `json:"memo"`
}
