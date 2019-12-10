package appcenter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
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
	httpClient *http.Client
	debug      bool
}

// Apps ...
func (c Client) Apps(owner, name string) App {
	return App{client: c, owner: owner, name: name}
}

// NewClient returns an AppCenter authenticated client
func NewClient(token string, debug bool) Client {
	return Client{
		httpClient: &http.Client{
			Transport: &roundTripper{
				token: token,
			},
		},
		debug: debug,
	}
}

func (c Client) jsonRequest(method, url string, body interface{}, response interface{}) (int, error) {
	var reader io.Reader

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return -1, err
		}
		reader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return -1, err
	}

	if c.debug {
		b, err := httputil.DumpRequest(req, true)
		if err != nil {
			return -1, err
		}
		fmt.Println(string(b))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return -1, err
	}

	if c.debug {
		b, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return -1, err
		}
		fmt.Println(string(b))
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

func (c Client) uploadForm(url string, files map[string]string) (int, error) {
	var (
		b bytes.Buffer
		w = multipart.NewWriter(&b)
	)

	for fileName, filePath := range files {
		f, err := os.Open(filePath)
		if err != nil {
			return -1, err
		}

		fw, err := w.CreateFormFile(fileName, filePath)
		if err != nil {
			return -1, err
		}

		if _, err = io.Copy(fw, f); err != nil {
			return -1, err
		}

		if err := f.Close(); err != nil {
			return -1, nil
		}
	}

	if err := w.Close(); err != nil {
		return -1, err
	}

	uploadReq, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return -1, err
	}

	uploadReq.Header.Set("Content-Type", w.FormDataContentType())

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

func (c Client) uploadFile(url string, filePath string) (int, error) {
	fb, err := ioutil.ReadFile(filePath)
	if err != nil {
		return -1, err
	}

	uploadReq, err := http.NewRequest("PUT", url, bytes.NewReader(fb))
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
