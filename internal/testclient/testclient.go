package testclient

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	httpClient http.Client
	URL        *url.URL
	Opts
}

type Opts struct {
	insecureSkipVerify bool
}

type OptFunc func(*Opts)

func defaultOpts() Opts {
	return Opts{
		insecureSkipVerify: false,
	}
}

func WithInsecureSkipVerify(b bool) OptFunc {
	return func(opts *Opts) {
		opts.insecureSkipVerify = b
	}
}

func NewClient(URL *url.URL, opts ...OptFunc) Client {
	o := defaultOpts()

	for _, fn := range opts {
		fn(&o)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: o.insecureSkipVerify},
	}
	httpClient := http.Client{Transport: tr}
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
