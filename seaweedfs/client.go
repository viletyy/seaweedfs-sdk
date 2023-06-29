package seaweedfs

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var jsonHeader = http.Header{"content-type": []string{"application/json"}}

func Version() string {
	return "0.0.1"
}

type Client struct {
	url    string
	client *http.Client
}

func NewClient(url string) *Client {
	return &Client{
		url:    strings.TrimSuffix(url, "/"),
		client: &http.Client{},
	}
}

func NewClientWithHTTP(url string, httpClient *http.Client) {
	client := NewClient(url)
	client.client = httpClient
}

func (c *Client) SetHttpClient(client *http.Client) {
	c.client = client
}

func (c *Client) doRequest(method, path string, header http.Header, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, c.url+path, body)
	if err != nil {
		return nil, err
	}

	for k, v := range header {
		req.Header[k] = v
	}

	return c.client.Do(req)
}

func (c *Client) getResponse(method, path string, header http.Header, body io.Reader) ([]byte, error) {
	resp, err := c.doRequest(method, path, header, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) getParsedResponse(method, path string, header http.Header, body io.Reader, obj interface{}) error {
	data, err := c.getResponse(method, path, header, body)
	if err != nil {
		return nil
	}

	return json.Unmarshal(data, obj)
}

func (c *Client) getStatusCode(method, path string, header http.Header, body io.Reader) (int, error) {
	resp, err := c.doRequest(method, path, header, body)
	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()
	return resp.StatusCode, nil
}
