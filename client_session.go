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

type sessionsClient struct {
	*Client
}

func (c *sessionsClient) Recent(expand []string, filter *SessionFilter, paging *Paging) (*Sessions, error) {
	params := paging.Params()
	for _, item := range expand {
		params.Add("expand", item)
	}
	if filter != nil && filter.UserID != "" {
		params.Add("user_id", filter.UserID)
	}
	resp, err := c.request(http.MethodGet, "/sessions", params, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		return ParseSessions(resp)
	}
	return nil, ParseError(resp)
}

func (c *sessionsClient) Create(session *Session, expand []string) (*Session, error) {
	params := make(url.Values)
	for _, item := range expand {
		params.Add("expand", item)
	}
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return nil, fmt.Errorf("error marshaling session json: %v", err)
	}
	resp, err := c.request(http.MethodPost, "/sessions", params, sessionJSON)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusCreated {
		return ParseSession(resp)
	}
	return nil, ParseError(resp)
}

func (c *sessionsClient) Delete(id string) error {
	resp, err := c.request(http.MethodDelete, "/sessions/"+id, nil, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return ParseError(resp)
}

func (c *sessionsClient) Get(id string) (*Session, error) {
	resp, err := c.request(http.MethodGet, "/sessions/"+id, nil, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		return ParseSession(resp)
	}
	return nil, ParseError(resp)
}

func (c *sessionsClient) Update(session *Session, overwriteDetails bool) error {
	method := http.MethodPatch
	if overwriteDetails {
		method = http.MethodPut
	}
	sessionJSON, err := session.MarshalForUpdate(false)
	if err != nil {
		return fmt.Errorf("error marshaling session json: %v", err)
	}
	resp, err := c.request(method, "/sessions/"+session.ID, nil, sessionJSON)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return ParseError(resp)
}

func (c *sessionsClient) Access(session *Session, expand []string) (*Session, error) {
	params := make(url.Values)
	for _, item := range expand {
		params.Add("expand", item)
	}
	sessionJSON, err := session.MarshalForUpdate(true)
	if err != nil {
		return nil, fmt.Errorf("error marshaling session json: %v", err)
	}
	resp, err := c.request(http.MethodPost, "/sessions/access", nil, sessionJSON)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		return ParseSession(resp)
	}
	return nil, ParseError(resp)
}
