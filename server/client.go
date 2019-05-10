// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mattermost/mattermost-server/model"
)

const (
	STEEDOS_API_URI  = "192.168.0.86"
	ZOOM_API_VERSION = "/v2"
	JWT_ALG          = "HS256"
)

type Client struct {
	ApiKey     string
	ApiSecret  string
	HttpClient *http.Client
	BaseUrl    string
}

// NewClient returns a new Zoom API client. An empty url will default to https://api.zoom.us/v2.
func NewClient(zoomUrl, apiKey, apiSecret string) *Client {
	if zoomUrl == "" {
		zoomUrl = (&url.URL{
			Scheme: "https",
			Host:   STEEDOS_API_URI,
			Path:   ZOOM_API_VERSION,
		}).String()
	}

	return &Client{
		ApiKey:     apiKey,
		ApiSecret:  apiSecret,
		HttpClient: &http.Client{},
		BaseUrl:    zoomUrl,
	}
}

func (c *Client) generateJWT(username string, sessionId string) (string, error) {
	claims := jwt.MapClaims{}

	claims["iss"] = c.ApiKey
	claims["exp"] = model.GetMillis() + (10 * 1000) // expire after 10s
	claims["username"] = username
	claims["sessionId"] = sessionId

	alg := jwt.GetSigningMethod(JWT_ALG)
	if alg == nil {
		return "", fmt.Errorf("Couldn't find signing method")
	}

	token := jwt.NewWithClaims(alg, claims)

	out, err := token.SignedString([]byte(c.ApiSecret))
	if err != nil {
		return "", err
	}

	return out, nil
}

func closeBody(r *http.Response) {
	if r.Body != nil {
		ioutil.ReadAll(r.Body)
		r.Body.Close()
	}
}

type ClientError struct {
	StatusCode int
	Err        string
}

func (ce *ClientError) Error() string {
	return ce.Err
}

func (c *Client) request(method string, path string, data interface{}, ret interface{}, username string, sessionId string) *ClientError {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return &ClientError{http.StatusInternalServerError, err.Error()}
	}

	rq, err := http.NewRequest(method, c.BaseUrl+path, bytes.NewReader(jsonData))
	if err != nil {
		return &ClientError{http.StatusInternalServerError, err.Error()}
	}
	rq.Header.Set("Content-Type", "application/json")
	rq.Close = true

	token, err := c.generateJWT(username, sessionId)
	if err != nil {
		return &ClientError{http.StatusInternalServerError, err.Error()}
	}
	rq.Header.Set("Authorization", "BEARER "+token)

	if rp, err := c.HttpClient.Do(rq); err != nil {
		return &ClientError{http.StatusInternalServerError, fmt.Sprintf("Unable to make request to %v: %v", c.BaseUrl+path, err.Error())}
	} else if rp == nil {
		return &ClientError{http.StatusInternalServerError, fmt.Sprintf("Received nil response when making request to %v", c.BaseUrl+path)}
	} else if rp.StatusCode >= 300 {
		defer closeBody(rp)
		buf := new(bytes.Buffer)
		buf.ReadFrom(rp.Body)
		return &ClientError{rp.StatusCode, buf.String()}
	} else {
		defer closeBody(rp)
		buf := new(bytes.Buffer)
		buf.ReadFrom(rp.Body)
		if err := json.Unmarshal(buf.Bytes(), &ret); err != nil {
			return &ClientError{rp.StatusCode, err.Error()}
		}
		return nil
	}
}
