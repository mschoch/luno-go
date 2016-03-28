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
	"net/url"
)

// EntityAggregate represents a response containing aggregte information
// about an entitites, Users, Sessions, Events
type EntityAggregate map[string]int

// ParseEntityAggregate parses an EntityAggregate out of an HTTP response
func ParseEntityAggregate(resp *http.Response) (EntityAggregate, error) {
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading luno response body: %v", err)
	}
	var rv EntityAggregate
	err = json.Unmarshal(respBytes, &rv)
	if err != nil {
		return nil, fmt.Errorf("error parsing luno analytic aggregate json: '%s' err: %v", respBytes, err)
	}

	return rv, nil
}

// EventAggregates contains a list of EventAggregate
type EventAggregates struct {
	List []*EventAggregate `json:"list"`
}

// ParseEventAggregates parses EventAggregates out of an HTTP response
func ParseEventAggregates(resp *http.Response) (*EventAggregates, error) {
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading luno response body: %v", err)
	}
	var rv EventAggregates
	err = json.Unmarshal(respBytes, &rv)
	if err != nil {
		return nil, fmt.Errorf("error parsing luno event aggregates json: '%s' err: %v", respBytes, err)
	}
	return &rv, nil
}

// EventAggregate represents aggregate information about an Event
type EventAggregate struct {
	Entity
	Name  string `json:"name"`
	Count int    `json:"count"`
	Last  string `json:"last"`
}

// TimelineFilter allows you to filter items returned from an Event Timeline
type TimelineFilter struct {
	Distinct   bool   `json:"distinct"`
	From       string `json:"from"`
	To         string `json:"to"`
	Group      string `json:"group"`
	Name       string `json:"name"`
	RoundRange bool   `json:"round_range"`
	UserID     string `json:"user_id"`
}

// Params converts a TimelineFilter into HTTP URL parameters
func (t *TimelineFilter) Params() url.Values {
	rv := make(url.Values)
	if t != nil {
		rv.Add("distinct", fmt.Sprintf("%t", t.Distinct))
		rv.Add("round_range", fmt.Sprintf("%t", t.RoundRange))
		if t.From != "" {
			rv.Add("from", t.From)
		}
		if t.To != "" {
			rv.Add("to", t.To)
		}
		if t.Group != "" {
			rv.Add("group", t.Group)
		}
		if t.Name != "" {
			rv.Add("name", t.Name)
		}
		if t.UserID != "" {
			rv.Add("user_id", t.UserID)
		}
	}
	return rv
}

// Range represents a date range
type Range struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// TimelineEntry represents an point on a Timeline
type TimelineEntry struct {
	Timestamp string `json:"timestamp"`
	Range     *Range `json:"range"`
	Count     int
}

// EventsTimeline represents a list of TimelineEntry objects
type EventsTimeline struct {
	Timeline []*TimelineEntry `json:"timeline"`
	Total    int              `json:"total"`
}

// ParseEventsTimeline parses an EventsTimeline out of an HTTP response
func ParseEventsTimeline(resp *http.Response) (*EventsTimeline, error) {
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading luno response body: %v", err)
	}
	var rv EventsTimeline
	err = json.Unmarshal(respBytes, &rv)
	if err != nil {
		return nil, fmt.Errorf("error parsing luno events timeline json: '%s' err: %v", respBytes, err)
	}
	return &rv, nil
}
