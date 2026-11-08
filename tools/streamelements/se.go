// Package streamelements is the main interface for Streamelements' API.
package streamelements

import (
	"bytes"
	"net/http"
	"strings"
)

// New returns a new Streamelements API connection with the given bearer token.
func New(token string) *Streamelements {
	if !strings.HasPrefix(token, "Bearer ") {
		token = "Bearer " + token
	}
	se := &Streamelements{
		c:     &http.Client{},
		token: token,
	}
	return se
}

// doReq makes a new request with the given properties and parameters.
func (se *Streamelements) doReq(method, path string, body []byte, header map[string]string) (*http.Response, error) {
	const baseURL = "https://api.streamelements.com/kappa/v2"
	req, err := http.NewRequest(method, baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}
	req.Header.Set("Authorization", se.token)

	return se.c.Do(req)
}
