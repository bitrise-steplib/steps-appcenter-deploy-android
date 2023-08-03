package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/retry"
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
	req.Header.Set(
		"x-api-token", rt.token,
	)
	req.Header.Set(
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
	retClient := retry.NewHTTPClient()
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
		if resp != nil {
			if err := resp.Body.Close(); err != nil {
				log.Warnf("failed to close body: %s", err)
			}
		}
	}()

	if resp != nil && response != nil {
		rb, err := io.ReadAll(resp.Body)
		if err != nil {
			return -1, err
		}

		if err := json.Unmarshal(rb, response); err != nil {
			reqDump, err := httputil.DumpRequestOut(resp.Request, true)
			if err != nil {
				log.Warnf("failed to dump request: %v", err)
			}

			respDump, err := httputil.DumpResponse(resp, false)
			if err != nil {
				log.TWarnf("failed to dump response: %s", err)
			}

			return resp.StatusCode, fmt.Errorf("failed to unmarshal response: %s, request: %s, response headers: %s response body: %s", err, reqDump, respDump, string(rb))
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
	fb, err := os.ReadFile(filePath)
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
			log.Warnf("failed to close body: %s", err)
		}
	}()

	return resp.StatusCode, nil
}
