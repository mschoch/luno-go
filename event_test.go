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

import "testing"

func TestEvents(t *testing.T) {
	apiKey, secretKey, err := lunoKeysFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	lunoClient := NewClient(apiKey, secretKey)

	// now create a user
	user := &User{
		Name:     "Charles Winchester",
		Email:    "chuck@af.com",
		Password: "imrich",
	}
	createdUser, err := lunoClient.Users.Create(user, true, nil)
	if err != nil {
		t.Fatal(err)
	}

	events, err := lunoClient.Events.Recent(nil, &EventFilter{UserID: createdUser.ID}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(events.List) != 1 {
		t.Errorf("expected 1 event, got %d", len(events.List))
	} else {
		if events.List[0].Name != "User Created" {
			t.Errorf("expected event named User Created, got %s", events.List[0].Name)
		}
	}

	// now create a custom event
	customEvent := &Event{
		UserID: createdUser.ID,
		Name:   "bad_hat",
	}
	createdEvent, err := lunoClient.Events.Create(customEvent, nil)
	if err != nil {
		t.Fatal(err)
	}

	// lookup the event we created
	lookupEvent, err := lunoClient.Events.Get(createdEvent.ID)
	if err != nil {
		t.Fatal(err)
	}

	// now lets update the event details
	lookupEvent.Details = map[string]interface{}{
		"color": "neon-green",
	}
	err = lunoClient.Events.Update(lookupEvent, true)
	if err != nil {
		t.Fatal(err)
	}

	// lookup the event after the change
	lookupEvent, err = lunoClient.Events.Get(createdEvent.ID)
	if err != nil {
		t.Fatal(err)
	}

	// verify the changed
	if details, ok := lookupEvent.Details.(map[string]interface{}); ok {
		if colorVal, ok := details["color"].(string); ok {
			if colorVal != "neon-green" {
				t.Errorf("expected session details for key 'color' to be 'neon-green', got '%s'", colorVal)
			}
		} else {
			t.Errorf("expected valuer for key 'color' to be string, got %T", details["color"])
		}
	} else {
		t.Errorf("expected details to be a map, got %T", lookupEvent.Details)
	}

	// now lets extend the event details
	lookupEvent.Details = map[string]interface{}{
		"size": "small",
	}
	err = lunoClient.Events.Update(lookupEvent, false)
	if err != nil {
		t.Fatal(err)
	}

	// lookup the event after the change
	lookupEvent, err = lunoClient.Events.Get(createdEvent.ID)
	if err != nil {
		t.Fatal(err)
	}

	// verify the changed
	if details, ok := lookupEvent.Details.(map[string]interface{}); ok {
		if colorVal, ok := details["color"].(string); ok {
			if colorVal != "neon-green" {
				t.Errorf("expected session details for key 'color' to be 'neon-green', got '%s'", colorVal)
			}
		} else {
			t.Errorf("expected valuer for key 'color' to be string, got %T", details["color"])
		}
		if sizeVal, ok := details["size"].(string); ok {
			if sizeVal != "small" {
				t.Errorf("expected session details for key 'size' to be 'small', got '%s'", sizeVal)
			}
		} else {
			t.Errorf("expected valuer for key 'size' to be string, got %T", details["size"])
		}
	} else {
		t.Errorf("expected details to be a map, got %T", lookupEvent.Details)
	}

	// delete the custom event
	err = lunoClient.Events.Delete(lookupEvent.ID)
	if err != nil {
		t.Fatal(err)
	}

	// delete the user
	err = lunoClient.Users.Delete(createdUser.ID)
	if err != nil {
		t.Fatal(err)
	}
}
