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

func TestSessionCrud(t *testing.T) {
	apiKey, secretKey, err := lunoKeysFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	lunoClient := NewClient(apiKey, secretKey)

	// try to get a non-existant session
	_, err = lunoClient.Sessions.Get("sess_xxxxxxxxxxxxxxxxxxxxxxxx")
	if !IsErrorCode(err, ErrCodeSessionNotFound) {
		t.Errorf("expected error session not found, got %v", err)
	}

	// get the list of sessions, expecting 0
	recentSessions, err := lunoClient.Sessions.Recent(nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(recentSessions.List) != 0 {
		t.Errorf("expected 0 sessions initially, got %d", len(recentSessions.List))
	}

	// try to create an anonymous session
	anonSession := &Session{
		UserID: "",
		Details: map[string]interface{}{
			"robot": "yes",
		},
	}
	createdSession, err := lunoClient.Sessions.Create(anonSession, nil)
	if err != nil {
		t.Fatal(err)
	}

	// get the session, should not error
	createdSession, err = lunoClient.Sessions.Get(createdSession.ID)
	if err != nil {
		t.Errorf("expected no error getting session, got %v", err)
	}

	currentAccessCount := createdSession.AccessCount

	// access the sessions
	createdSession, err = lunoClient.Sessions.Access(createdSession, nil)
	if createdSession.AccessCount != currentAccessCount+1 {
		t.Errorf("expected access count to go up by 1, was %d now %d", currentAccessCount, createdSession.AccessCount)
	}

	// update the sessions, overwriting details
	createdSession.Details = map[string]interface{}{
		"mister": "robot",
	}
	err = lunoClient.Sessions.Update(createdSession, true)
	if err != nil {
		t.Fatal(err)
	}

	// get the session, and check the details
	createdSession, err = lunoClient.Sessions.Get(createdSession.ID)
	if err != nil {
		t.Errorf("expected no error getting session, got %v", err)
	}
	if details, ok := createdSession.Details.(map[string]interface{}); ok {
		if misterVal, ok := details["mister"].(string); ok {
			if misterVal != "robot" {
				t.Errorf("expected session details for key 'mister' to be 'robot', got '%s'", misterVal)
			}
		} else {
			t.Errorf("expected valuer for key 'mister' to be string, got %T", details["mister"])
		}
	} else {
		t.Errorf("expected details to be a map, got %T", createdSession.Details)
	}

	// update the session again, extending details
	createdSession.Details = map[string]interface{}{
		"crazy": "maybe",
	}
	err = lunoClient.Sessions.Update(createdSession, false)
	if err != nil {
		t.Fatal(err)
	}

	// get the session, and check the details
	createdSession, err = lunoClient.Sessions.Get(createdSession.ID)
	if err != nil {
		t.Errorf("expected no error getting session, got %v", err)
	}
	if details, ok := createdSession.Details.(map[string]interface{}); ok {
		if misterVal, ok := details["mister"].(string); ok {
			if misterVal != "robot" {
				t.Errorf("expected session details for key 'mister' to be 'robot', got '%s'", misterVal)
			}
		} else {
			t.Errorf("expected valuer for key 'mister' to be string, got %T", details["mister"])
		}
		if crazyVal, ok := details["crazy"].(string); ok {
			if crazyVal != "maybe" {
				t.Errorf("expected session details for key 'crazy' to be 'maybe', got '%s'", crazyVal)
			}
		} else {
			t.Errorf("expected valuer for key 'crazy' to be string, got %T", details["crazy"])
		}
	} else {
		t.Errorf("expected details to be a map, got %T", createdSession.Details)
	}

	// get the list of sessions, expecting 1
	recentSessions, err = lunoClient.Sessions.Recent(nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(recentSessions.List) != 1 {
		t.Errorf("expected 1 sessions initially, got %d", len(recentSessions.List))
	}

	// delete the sessions we created
	err = lunoClient.Sessions.Delete(createdSession.ID)
	if err != nil {
		t.Fatal(err)
	}

	// get the list of sessions again, expecting 0
	recentSessions, err = lunoClient.Sessions.Recent(nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(recentSessions.List) != 0 {
		t.Errorf("expected 0 sessions at end, got %d", len(recentSessions.List))
	}
}
