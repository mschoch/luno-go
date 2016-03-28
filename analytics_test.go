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
	"reflect"
	"testing"
)

func TestAnalytics(t *testing.T) {
	apiKey, secretKey, err := lunoKeysFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	lunoClient := NewClient(apiKey, secretKey)

	// create a user and login, so there are some things to look at
	user := &User{
		Name:     "Ducker Cup",
		Email:    "d@c.com",
		Password: "quack",
	}
	createdUser, err := lunoClient.Users.Create(user, true, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = lunoClient.Users.LoginWithEmail(user.Email, user.Password, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	// get default user analytics
	analytics, err := lunoClient.Analytics.Users(nil)
	if err != nil {
		t.Fatal(err)
	}
	expected := EntityAggregate{
		"total":   1,
		"7_days":  1,
		"28_days": 1,
	}
	if !reflect.DeepEqual(expected, analytics) {
		t.Errorf("expected %v, got %v", expected, analytics)
	}

	// get user analytics with custom days
	analytics, err = lunoClient.Analytics.Users([]string{"3", "9"})
	if err != nil {
		t.Fatal(err)
	}
	expected = EntityAggregate{
		"3_days": 1,
		"9_days": 1,
	}
	if !reflect.DeepEqual(expected, analytics) {
		t.Errorf("expected %v, got %v", expected, analytics)
	}

	// get default session analytics
	analytics, err = lunoClient.Analytics.Sessions(nil)
	if err != nil {
		t.Fatal(err)
	}
	expected = EntityAggregate{
		"total":   1,
		"7_days":  1,
		"28_days": 1,
	}
	if !reflect.DeepEqual(expected, analytics) {
		t.Errorf("expected %v, got %v", expected, analytics)
	}

	// get session analytics with custom days
	analytics, err = lunoClient.Analytics.Sessions([]string{"3", "9"})
	if err != nil {
		t.Fatal(err)
	}
	expected = EntityAggregate{
		"3_days": 1,
		"9_days": 1,
	}
	if !reflect.DeepEqual(expected, analytics) {
		t.Errorf("expected %v, got %v", expected, analytics)
	}

	// get default event analytics
	analytics, err = lunoClient.Analytics.Events(nil)
	if err != nil {
		t.Fatal(err)
	}
	expected = EntityAggregate{
		"total":   4,
		"7_days":  4,
		"28_days": 4,
	}
	if !reflect.DeepEqual(expected, analytics) {
		t.Errorf("expected %v, got %v", expected, analytics)
	}

	// get event analytics with custom days
	analytics, err = lunoClient.Analytics.Events([]string{"3", "9"})
	if err != nil {
		t.Fatal(err)
	}
	expected = EntityAggregate{
		"3_days": 4,
		"9_days": 4,
	}
	if !reflect.DeepEqual(expected, analytics) {
		t.Errorf("expected %v, got %v", expected, analytics)
	}

	eventAggregates, err := lunoClient.Analytics.EventsList()
	if err != nil {
		t.Fatal(err)
	}
	expectedAggregateCounts := map[string]int{
		"User Created":     1,
		"Correct Password": 1,
		"Session Created":  1,
		"Logged In":        1,
	}

	actualAggregateCounts := make(map[string]int)
	if len(eventAggregates.List) > 0 {
		for _, eventAgg := range eventAggregates.List {
			actualAggregateCounts[eventAgg.Name] = eventAgg.Count
		}
	}

	if !reflect.DeepEqual(expectedAggregateCounts, actualAggregateCounts) {
		t.Errorf("expected aggregate counts: %v, got %v", expectedAggregateCounts, actualAggregateCounts)
	}

	eventsTimeline, err := lunoClient.Analytics.EventsTimeline(nil)
	if err != nil {
		t.Fatal(err)
	}
	if eventsTimeline.Total != 4 {
		t.Errorf("expected 4 events in timeline, got %d", eventsTimeline.Total)
	}

	// delete the user
	err = lunoClient.Users.Delete(createdUser.ID)
	if err != nil {
		t.Fatal(err)
	}
}
