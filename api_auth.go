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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// APIAuthFilter allows you tou filter ApiAuth objects
type APIAuthFilter struct {
	UserID string `json:"user_id"`
}

// APIAuths represents a list of ApiAuth objects
type APIAuths struct {
	Entity
	List []*APIAuth `json:"list"`
	Page Page       `json:"page"`
}

// ParseAPIAuths extracts APIAuths from an HTTP response
func ParseAPIAuths(resp *http.Response) (*APIAuths, error) {
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading luno response body: %v", err)
	}
	var rv APIAuths
	err = json.Unmarshal(respBytes, &rv)
	if err != nil {
		return nil, fmt.Errorf("error parsing luno api auths json: '%s' err: %v", respBytes, err)
	}
	return &rv, nil
}

// APIAuth represents the Luno API Authetication - https://luno.io/docs#api_authentication
type APIAuth struct {
	Entity
	UserID  string      `json:"user_id,omitempty"`
	Key     string      `json:"key,omitempty"`
	Secret  string      `json:"secret,omitempty"`
	Created string      `json:"created,omitempty"`
	Details interface{} `json:"details,omitempty"`
	User    *User       `json:"user,omitempty"`
}

// MarshalForUpdate exports only the fields suitable for update
func (a *APIAuth) MarshalForUpdate() ([]byte, error) {
	tmp := map[string]interface{}{
		"details": a.Details,
	}
	return json.Marshal(tmp)
}

// ParseAPIAuth parses an APIAuth out of an HTTP response
func ParseAPIAuth(resp *http.Response) (*APIAuth, error) {
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading luno response body: %v", err)
	}
	var rv APIAuth
	err = json.Unmarshal(respBytes, &rv)
	if err != nil {
		return nil, fmt.Errorf("error parsing luno api auth json: '%s' err: %v", respBytes, err)
	}
	return &rv, nil
}
