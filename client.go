//  Copyright (c) 2016 Marty Schoch
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the
//  License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing,
//  software distributed under the License is distributed on an "AS
//  IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
//  express or implied. See the License for the specific language
//  governing permissions and limitations under the License.

package luno

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Log allows an application to provie a custom logger for logging
// care must be taken, as logs may contain sensitive information
var Log *log.Logger

// LogRequestURL allows you to log the URL of every HTTP request
var LogRequestURL = false

// LogRequestBody allows you to log the body of every HTTP request
var LogRequestBody = false

// LogResponseCode allows you to log every HTTP response code
var LogResponseCode = false

// LogResponseBody allows you to log every HTTP response body
var LogResponseBody = false

// Client is a Luno Client - https://luno.io/docs/libraries
type Client struct {
	host       string
	version    string
	apiKey     string
	secretKey  string
	httpClient *http.Client

	Users     *usersClient
	Events    *eventsClient
	Sessions  *sessionsClient
	APIAuth   *apiAuthClient
	Analytics *analyticsClient
	Account   *accountClient
}

// NewClient builds a new client with the provided API key and secret key
func NewClient(apiKey, secretKey string) *Client {
	rv := &Client{
		host:      "api.luno.io",
		version:   "v1",
		apiKey:    apiKey,
		secretKey: secretKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	rv.Users = &usersClient{rv}
	rv.Events = &eventsClient{rv}
	rv.Sessions = &sessionsClient{rv}
	rv.APIAuth = &apiAuthClient{rv}
	rv.Analytics = &analyticsClient{rv}
	rv.Account = &accountClient{rv}
	return rv
}

func (c *Client) request(method, endpoint string, params url.Values, body []byte) (*http.Response, error) {
	if params == nil {
		params = make(url.Values)
	}
	params.Add("key", c.apiKey)
	params.Add("timestamp", c.timestamp())
	req := &http.Request{
		Method: method,
		URL: &url.URL{
			Host:   c.host,
			Scheme: "https",
			Opaque: "/" + c.version + endpoint + "?" + params.Encode(),
		},
		Body:          ioutil.NopCloser(bytes.NewReader(body)),
		ContentLength: int64(len(body)),
		Header:        make(http.Header),
	}
	req.Header.Set("Content-Type", "application/json")

	// sign and add signature to request
	sign := c.signRequest(req, body)
	req.URL.Opaque += "&sign=" + sign
	if Log != nil && LogRequestURL {
		Log.Printf("%s request: %s", req.Method, req.URL.Opaque)
	}
	if len(body) > 0 && Log != nil && LogRequestBody {
		Log.Print(string(body))
	}
	resp, err := c.httpClient.Do(req)
	if Log != nil && LogResponseCode {
		Log.Printf("response %d", resp.StatusCode)
	}
	if resp.Body != nil && Log != nil && LogResponseBody {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body for logging")
		}
		Log.Print(string(respBody))
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(respBody))
	}

	return resp, err
}

func (c *Client) timestamp() string {
	now := time.Now()
	return now.Format(time.RFC3339)
}

func (c *Client) signRequest(req *http.Request, body []byte) string {
	msg := req.Method + ":" + req.URL.Opaque
	if len(body) > 0 {
		msg += ":" + string(body)
	}
	mac := hmac.New(sha512.New, []byte(c.secretKey))
	mac.Write([]byte(msg))
	return hex.EncodeToString(mac.Sum(nil))
}
