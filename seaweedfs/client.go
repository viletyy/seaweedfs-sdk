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
	masterUrl string
	volumnUrl string
	client    *http.Client
}

func NewClient(murl, vUrl string) *Client {
	return &Client{
		masterUrl: strings.TrimSuffix(murl, "/"),
		volumnUrl: strings.TrimSuffix(vUrl, "/"),

		client: &http.Client{},
	}
}

func NewClientWithHTTP(murl, vUrl string, httpClient *http.Client) {
	client := NewClient(murl, vUrl)
	client.client = httpClient
}

func (c *Client) SetHttpClient(client *http.Client) {
	c.client = client
}

func (c *Client) doMasterRequest(method, path string, header http.Header, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, c.masterUrl+path, body)
	if err != nil {
		return nil, err
	}

	for k, v := range header {
		req.Header[k] = v
	}

	return c.client.Do(req)
}

func (c *Client) doVolumnRequest(method, path string, header http.Header, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, c.volumnUrl+path, body)
	if err != nil {
		return nil, err
	}

	for k, v := range header {
		req.Header[k] = v
	}

	return c.client.Do(req)
}

func (c *Client) getMasterResponse(method, path string, header http.Header, body io.Reader) ([]byte, error) {
	resp, err := c.doMasterRequest(method, path, header, body)
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

func (c *Client) getVolumnResponse(method, path string, header http.Header, body io.Reader) ([]byte, error) {
	resp, err := c.doVolumnRequest(method, path, header, body)
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

func (c *Client) getMasterParsedResponse(method, path string, header http.Header, body io.Reader, obj interface{}) error {
	data, err := c.getMasterResponse(method, path, header, body)
	if err != nil {
		return nil
	}

	return json.Unmarshal(data, obj)
}

func (c *Client) getVolumnParsedResponse(method, path string, header http.Header, body io.Reader, obj interface{}) error {
	data, err := c.getVolumnResponse(method, path, header, body)
	if err != nil {
		return nil
	}

	return json.Unmarshal(data, obj)
}

func (c *Client) getMasterStatusCode(method, path string, header http.Header, body io.Reader) (int, error) {
	resp, err := c.doMasterRequest(method, path, header, body)
	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func (c *Client) getVolumnStatusCode(method, path string, header http.Header, body io.Reader) (int, error) {
	resp, err := c.doVolumnRequest(method, path, header, body)
	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()
	return resp.StatusCode, nil
}
