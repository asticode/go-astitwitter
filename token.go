package astitwitter

import (
	"bytes"
	"net/http"

	"github.com/pkg/errors"
)

type TokenBody struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func (c *Client) bearerToken() (t string, err error) {
	// Send
	var b TokenBody
	if err = c.send(
		http.MethodPost,
		"/oauth2/token",
		bytes.NewBufferString("grant_type=client_credentials"),
		func(r *http.Request) {
			r.SetBasicAuth(c.apiKey, c.apiSecretKey)
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		},
		&b,
	); err != nil {
		err = errors.Wrap(err, "astitwitter: sending failed")
		return
	}

	// Set token
	t = b.AccessToken
	return
}
