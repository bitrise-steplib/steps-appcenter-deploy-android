package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

const (
	baseURL = `https://api.appcenter.ms`
)

type roundTripper struct {
	token string
}

// RoundTrip ...
func (rt roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add(
		"x-api-token", rt.token,
	)
	req.Header.Add(
		"content-type", "application/json; charset=utf-8",
	)
	return http.DefaultTransport.RoundTrip(req)
}

// Client ...
type Client struct {
	httpClient *retryablehttp.Client
}

// NewClient returns an AppCenter authenticated client
func NewClient(token string) Client {
	retClient := retryablehttp.NewClient()

	retClient.RetryMax = 5
	retClient.RetryWaitMin = 5 * time.Second
	retClient.RetryWaitMax = 10 * time.Second

	retClient.HTTPClient.Transport = &roundTripper{
		token: token,
	}

	return Client{
		httpClient: retClient,
	}
}

func (c Client) jsonRequest(method, url string, body []byte, response interface{}) (int, error) {
	var reader io.Reader

	if body != nil {
		reader = bytes.NewReader(body)
	}

	req, err := retryablehttp.NewRequest(method, url, reader)

	if err != nil {
		return -1, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return -1, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
		}
	}()

	if response != nil {
		rb, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return -1, err
		}

		if err := json.Unmarshal(rb, response); err != nil {
			return resp.StatusCode, fmt.Errorf("error: %s, response: %s", err, string(rb))
		}
	}

	return resp.StatusCode, nil
}

// MarshallContent ...
func (c Client) MarshallContent(content interface{}) ([]byte, error) {
	b, err := json.Marshal(content)
	if err != nil {
		return []byte{}, err
	}

	return b, err
}

func (c Client) uploadFile(url string, filePath string) (int, error) {
	fb, err := ioutil.ReadFile(filePath)
	if err != nil {
		return -1, err
	}

	uploadReq, err := retryablehttp.NewRequest("PUT", url, bytes.NewReader(fb))
	if err != nil {
		return -1, err
	}

	uploadReq.Header.Set("x-ms-blob-type", "BlockBlob")
	uploadReq.Header.Set("content-length", strconv.Itoa(len(fb)))

	resp, err := c.httpClient.Do(uploadReq)
	if err != nil {
		return -1, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
		}
	}()

	return resp.StatusCode, nil
}
