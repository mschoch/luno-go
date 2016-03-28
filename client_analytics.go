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
	"net/http"
	"net/url"
)

type analyticsClient struct {
	*Client
}

func (c *analyticsClient) Users(days []string) (EntityAggregate, error) {
	params := make(url.Values)
	for _, day := range days {
		params.Add("days", day)
	}
	resp, err := c.request(http.MethodGet, "/analytics/users", params, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		return ParseEntityAggregate(resp)
	}
	return nil, ParseError(resp)
}

func (c *analyticsClient) Sessions(days []string) (EntityAggregate, error) {
	params := make(url.Values)
	for _, day := range days {
		params.Add("days", day)
	}
	resp, err := c.request(http.MethodGet, "/analytics/sessions", params, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		return ParseEntityAggregate(resp)
	}
	return nil, ParseError(resp)
}

func (c *analyticsClient) Events(days []string) (EntityAggregate, error) {
	params := make(url.Values)
	for _, day := range days {
		params.Add("days", day)
	}
	resp, err := c.request(http.MethodGet, "/analytics/events", params, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		return ParseEntityAggregate(resp)
	}
	return nil, ParseError(resp)
}

func (c *analyticsClient) EventsList() (*EventAggregates, error) {
	resp, err := c.request(http.MethodGet, "/analytics/events/list", nil, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		return ParseEventAggregates(resp)
	}
	return nil, ParseError(resp)
}

func (c *analyticsClient) EventsTimeline(filter *TimelineFilter) (*EventsTimeline, error) {
	params := filter.Params()
	resp, err := c.request(http.MethodGet, "/analytics/events/timeline", params, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		return ParseEventsTimeline(resp)
	}
	return nil, ParseError(resp)
}
