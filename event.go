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

// Events represents a list of Event objects
type Events struct {
	Entity
	List []*Event `json:"list"`
	Page Page     `json:"page"`
}

// ParseEvents parses Events out of an HTTP response
func ParseEvents(resp *http.Response) (*Events, error) {
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading luno response body: %v", err)
	}
	var rv Events
	err = json.Unmarshal(respBytes, &rv)
	if err != nil {
		return nil, fmt.Errorf("error parsing luno events json: '%s' err: %v", respBytes, err)
	}
	return &rv, nil
}

// EventFilter lets you filter events
type EventFilter struct {
	Name   string `json:"name"`
	UserID string `json:"user_id"`
}

// Event represents a Luno Event - https://luno.io/docs#events
type Event struct {
	Entity
	UserID    string      `json:"user_id,omitempty"`
	Timestamp string      `json:"timestamp,omitempty"`
	Name      string      `json:"name,omitempty"`
	Details   interface{} `json:"details,omitempty"`
	User      *User       `json:"user,omitempty"`
}

// MarshalForUpdate exports the event fields suitable for an update operation
func (e *Event) MarshalForUpdate() ([]byte, error) {
	tmp := map[string]interface{}{
		"details": e.Details,
	}
	return json.Marshal(tmp)
}

// ParseEvent parses an Event out of an HTTP response
func ParseEvent(resp *http.Response) (*Event, error) {
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading luno response body: %v", err)
	}
	var rv Event
	err = json.Unmarshal(respBytes, &rv)
	if err != nil {
		return nil, fmt.Errorf("error parsing luno event json: '%s' err: %v", respBytes, err)
	}
	return &rv, nil
}
