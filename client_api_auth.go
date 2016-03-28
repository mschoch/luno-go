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
	"net/http"
	"net/url"
)

type apiAuthClient struct {
	*Client
}

func (c *apiAuthClient) Recent(expand []string, filter *APIAuthFilter, paging *Paging) (*APIAuths, error) {
	params := paging.Params()
	for _, item := range expand {
		params.Add("expand", item)
	}
	if filter != nil && filter.UserID != "" {
		params.Add("user_id", filter.UserID)
	}
	resp, err := c.request(http.MethodGet, "/api_authentication", params, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		return ParseAPIAuths(resp)
	}
	return nil, ParseError(resp)
}

func (c *apiAuthClient) Create(apiAuth *APIAuth, expand []string) (*APIAuth, error) {
	params := make(url.Values)
	for _, item := range expand {
		params.Add("expand", item)
	}
	apiAuthJSON, err := json.Marshal(apiAuth)
	if err != nil {
		return nil, fmt.Errorf("error marshaling session json: %v", err)
	}
	resp, err := c.request(http.MethodPost, "/api_authentication", params, apiAuthJSON)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusCreated {
		return ParseAPIAuth(resp)
	}
	return nil, ParseError(resp)
}

func (c *apiAuthClient) Get(id string, expand []string) (*APIAuth, error) {
	params := make(url.Values)
	for _, item := range expand {
		params.Add("expand", item)
	}
	resp, err := c.request(http.MethodGet, "/api_authentication/"+id, params, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		return ParseAPIAuth(resp)
	}
	return nil, ParseError(resp)
}

func (c *apiAuthClient) Update(apiAuth *APIAuth, overwriteProfile bool) error {
	method := http.MethodPatch
	if overwriteProfile {
		method = http.MethodPut
	}
	apiAuthJSON, err := apiAuth.MarshalForUpdate()
	if err != nil {
		return fmt.Errorf("error marshaling api auth json: %v", err)
	}
	resp, err := c.request(method, "/api_authentication/"+apiAuth.Key, nil, apiAuthJSON)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return ParseError(resp)
}

func (c *apiAuthClient) Delete(id string) error {
	resp, err := c.request(http.MethodDelete, "/api_authentication/"+id, nil, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return ParseError(resp)
}
