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

type eventsClient struct {
	*Client
}

func (c *eventsClient) Recent(expand []string, filter *EventFilter, paging *Paging) (*Events, error) {
	params := paging.Params()
	for _, item := range expand {
		params.Add("expand", item)
	}
	if filter != nil {
		if filter.UserID != "" {
			params.Add("user_id", filter.UserID)
		}
		if filter.Name != "" {
			params.Add("name", filter.Name)
		}
	}
	resp, err := c.request(http.MethodGet, "/events", params, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		return ParseEvents(resp)
	}
	return nil, ParseError(resp)
}

func (c *eventsClient) Create(event *Event, expand []string) (*Event, error) {
	params := make(url.Values)
	for _, item := range expand {
		params.Add("expand", item)
	}
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("error marshaling event json: %v", err)
	}
	resp, err := c.request(http.MethodPost, "/events", params, eventJSON)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusCreated {
		return ParseEvent(resp)
	}
	return nil, ParseError(resp)
}

func (c *eventsClient) Get(id string) (*Event, error) {
	resp, err := c.request(http.MethodGet, "/events/"+id, nil, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		return ParseEvent(resp)
	}
	return nil, ParseError(resp)
}

func (c *eventsClient) Update(event *Event, overwriteDetails bool) error {
	method := http.MethodPatch
	if overwriteDetails {
		method = http.MethodPut
	}
	eventJSON, err := event.MarshalForUpdate()
	if err != nil {
		return fmt.Errorf("error marshaling event json: %v", err)
	}
	resp, err := c.request(method, "/events/"+event.ID, nil, eventJSON)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return ParseError(resp)
}

func (c *eventsClient) Delete(id string) error {
	resp, err := c.request(http.MethodDelete, "/events/"+id, nil, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return ParseError(resp)
}
