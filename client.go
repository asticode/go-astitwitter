package astitwitter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/asticode/go-astikit"
)

const (
	baseURL                        = "https://api.twitter.com"
	errorCodeInvalidOrExpiredToken = 89
)

// Client represents the client
type Client struct {
	apiKey       string
	apiSecretKey string
	s            *astikit.HTTPSender
	t            string
}

// New creates a new client
func New(c Configuration) *Client {
	return &Client{
		apiKey:       c.APIKey,
		apiSecretKey: c.APISecretKey,
		s:            astikit.NewHTTPSender(c.Sender),
	}
}

type ErrorBody struct {
	Errors Errors `json:"errors"`
}

type Errors []Error

func (e Errors) Error() string {
	var ss []string
	for _, e := range e {
		ss = append(ss, fmt.Sprintf("code: %d - message: %s", e.Code, e.Message))
	}
	return strings.Join(ss, " | ")
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (c *Client) send(method, url string, body io.Reader, reqFunc func(r *http.Request), respPayload interface{}) (err error) {
	// Create request
	var req *http.Request
	if req, err = http.NewRequest(method, baseURL+url, body); err != nil {
		err = fmt.Errorf("astitwitter: creating %s request to %s failed: %w", method, url, err)
		return
	}

	// Adapt request
	if reqFunc != nil {
		reqFunc(req)
	}

	// Send
	var resp *http.Response
	if resp, err = c.s.Send(req); err != nil {
		err = fmt.Errorf("astitwitter: sending %s request to %s failed: %w", req.Method, req.URL.Path, err)
		return
	}
	defer resp.Body.Close()

	// Process error
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		// Unmarshal
		var b ErrorBody
		if err = json.NewDecoder(resp.Body).Decode(&b); err != nil {
			err = fmt.Errorf("astitwitter: unmarshaling errors failed: %w", err)
			return
		}

		// Set error
		err = b.Errors
		return
	}

	// Unmarshal
	if respPayload != nil {
		if err = json.NewDecoder(resp.Body).Decode(respPayload); err != nil {
			err = fmt.Errorf("astitwitter: unmarshaling failed: %w", err)
			return
		}
	}
	return
}

func (c *Client) sendAuthenticated(method, url string, body io.Reader, reqFunc func(r *http.Request), respPayload interface{}) (err error) {
	// No token
	if c.t == "" {
		// Get bearer token
		if c.t, err = c.bearerToken(); err != nil {
			err = fmt.Errorf("astitwitter: getting bearer token failed: %w", err)
			return
		}
	}

	// Bearer authorization
	authenticatedReqFunc := func(r *http.Request) {
		if reqFunc != nil {
			reqFunc(r)
		}
		r.Header.Set("Authorization", "Bearer "+c.t)
	}

	// Send
	if err = c.send(method, url, body, authenticatedReqFunc, respPayload); err != nil {
		// Assert error
		var retry bool
		if es, ok := err.(Errors); ok {
			// Loop through errors
			for _, e := range es {
				// Token is either invalid or expired
				if e.Code == errorCodeInvalidOrExpiredToken {
					// Get bearer token
					if c.t, err = c.bearerToken(); err != nil {
						err = fmt.Errorf("astitwitter: getting bearer token failed: %w", err)
						return
					}

					// Make sure to retry
					retry = true
					break
				}
			}
		}

		// Retry
		if retry {
			err = c.send(method, url, body, authenticatedReqFunc, respPayload)
		}

		// Process error
		if err != nil {
			err = fmt.Errorf("astitwitter: sending failed: %w", err)
			return
		}
	}
	return
}
