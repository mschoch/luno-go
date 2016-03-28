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

func TestApiAuth(t *testing.T) {
	apiKey, secretKey, err := lunoKeysFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	lunoClient := NewClient(apiKey, secretKey)

	// get the list of api auths
	apiAuths, err := lunoClient.APIAuth.Recent(nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(apiAuths.List) != 0 {
		t.Errorf("expected 0 api auths, got %d", len(apiAuths.List))
	}

	// create a user
	newUser := &User{
		Name:     "API User",
		Email:    "api@user.com",
		Password: "luv2code",
	}
	createdUser, err := lunoClient.Users.Create(newUser, true, nil)
	if err != nil {
		t.Fatal(err)
	}

	// create an api auth
	apiAuth := &APIAuth{
		UserID: createdUser.ID,
	}
	createdAPIAuth, err := lunoClient.APIAuth.Create(apiAuth, nil)
	if err != nil {
		t.Fatal(err)
	}

	// get the api auth we created
	_, err = lunoClient.APIAuth.Get(createdAPIAuth.Key, nil)
	if err != nil {
		t.Fatal(err)
	}

	// update the api auth, add details
	createdAPIAuth.Details = map[string]interface{}{
		"hidden": "secret",
	}
	err = lunoClient.APIAuth.Update(createdAPIAuth, true)
	if err != nil {
		t.Fatal(err)
	}

	// get the api auth, and check the details
	lookupAPIAuth, err := lunoClient.APIAuth.Get(createdAPIAuth.Key, nil)
	if err != nil {
		t.Errorf("expected no error getting session, got %v", err)
	}
	if details, ok := lookupAPIAuth.Details.(map[string]interface{}); ok {
		if hiddenVal, ok := details["hidden"].(string); ok {
			if hiddenVal != "secret" {
				t.Errorf("expected session details for key 'mister' to be 'secret', got '%s'", hiddenVal)
			}
		} else {
			t.Errorf("expected valuer for key 'hidden' to be string, got %T", details["hidden"])
		}
	} else {
		t.Errorf("expected details to be a map, got %T", lookupAPIAuth.Details)
	}

	// update the api auth again, extending details
	createdAPIAuth.Details = map[string]interface{}{
		"secret": "stash",
	}
	err = lunoClient.APIAuth.Update(createdAPIAuth, false)
	if err != nil {
		t.Fatal(err)
	}

	// get the api auth, and check the details
	lookupAPIAuth, err = lunoClient.APIAuth.Get(createdAPIAuth.Key, nil)
	if err != nil {
		t.Errorf("expected no error getting session, got %v", err)
	}
	if details, ok := lookupAPIAuth.Details.(map[string]interface{}); ok {
		if hiddenVal, ok := details["hidden"].(string); ok {
			if hiddenVal != "secret" {
				t.Errorf("expected session details for key 'mister' to be 'secret', got '%s'", hiddenVal)
			}
		} else {
			t.Errorf("expected valuer for key 'hidden' to be string, got %T", details["hidden"])
		}
		if secretVal, ok := details["secret"].(string); ok {
			if secretVal != "stash" {
				t.Errorf("expected session details for key 'secret' to be 'stash', got '%s'", secretVal)
			}
		} else {
			t.Errorf("expected valuer for key 'secret' to be string, got %T", details["secret"])
		}
	} else {
		t.Errorf("expected details to be a map, got %T", lookupAPIAuth.Details)
	}

	// delete the api auth
	err = lunoClient.APIAuth.Delete(createdAPIAuth.Key)
	if err != nil {
		t.Fatal(err)
	}

	// delete the user
	err = lunoClient.Users.Delete(createdUser.ID)
	if err != nil {
		t.Fatal(err)
	}
}
