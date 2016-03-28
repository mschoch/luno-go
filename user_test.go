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

func TestUserCrud(t *testing.T) {
	apiKey, secretKey, err := lunoKeysFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	lunoClient := NewClient(apiKey, secretKey)

	// get the list of users, expecting 0
	recentUsers, err := lunoClient.Users.Recent(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(recentUsers.List) != 0 {
		t.Errorf("expected 0 users initially, got %d", len(recentUsers.List))
	}

	// now create a user
	newUser := &User{
		Name:     "Bozo Clown",
		Email:    "bozo@clown.com",
		Password: "h8clownz",
		Profile: map[string]interface{}{
			"title": "clown king",
		},
	}
	createdUser, err := lunoClient.Users.Create(newUser, true, nil)
	if err != nil {
		t.Fatal(err)
	}

	// validate wrong password
	err = lunoClient.Users.ValidatePassword(createdUser.ID, "wrongpassword")
	if !IsErrorCode(err, ErrCodeIncorrectPassword) {
		t.Errorf("expected incorrect password, got %v", err)
	}

	// valid correct password
	err = lunoClient.Users.ValidatePassword(createdUser.ID, newUser.Password)
	if err != nil {
		t.Errorf("validate password failed")
	}

	// get the list of users again, expecting 1
	recentUsers, err = lunoClient.Users.Recent(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(recentUsers.List) != 1 {
		t.Errorf("expected 0 users initially, got %d", len(recentUsers.List))
	}

	// try to get the user we just created by id
	findUserbyID, err := lunoClient.Users.Get(createdUser.ID)
	if err != nil {
		t.Fatal(err)
	}
	if findUserbyID.Name != createdUser.Name {
		t.Errorf("expected username to match, got %s != %s", findUserbyID.Name, createdUser.Name)
	}

	// let's change our  name and replace the profile
	findUserbyID.Name = "Bill Clown"
	findUserbyID.Profile = map[string]interface{}{
		"instrument": "trumpet",
	}
	err = lunoClient.Users.Update(findUserbyID, true, true)
	if err != nil {
		t.Fatal(err)
	}

	// get the user we just updated by id
	findUserbyIDAgain, err := lunoClient.Users.Get(createdUser.ID)
	if err != nil {
		t.Fatal(err)
	}
	if findUserbyIDAgain.Name != findUserbyID.Name {
		t.Errorf("expected username to match, got %s != %s", findUserbyIDAgain.Name, findUserbyID.Name)
	}
	// since auto_name was true, first and last name should also be updated
	if findUserbyIDAgain.FirstName != "Bill" {
		t.Errorf("expected auto_name=true to set first name to Bill, got %s", findUserbyIDAgain.FirstName)
	}
	if findUserbyIDAgain.LastName != "Clown" {
		t.Errorf("expected auto_name=true to set last name to Clown, got %s", findUserbyIDAgain.LastName)
	}
	// since we set overwrite to true, old profile info should be gone
	if profileMap, ok := findUserbyIDAgain.Profile.(map[string]interface{}); ok {
		if profileMap["title"] != nil {
			t.Errorf("expected no profile key for title, got %v", profileMap["title"])
		}
		if instrument, ok := profileMap["instrument"].(string); ok {
			if instrument != "trumpet" {
				t.Errorf("expected instrument to be trumpet, got %s", instrument)
			}
		} else {
			t.Errorf("expected instrument to be string, got %T", profileMap["instrument"])
		}
	} else {
		t.Errorf("expected profile to be map, got %T", findUserbyIDAgain.Profile)
	}

	// now lets just try to extend the profile
	findUserbyIDAgain.Profile = map[string]interface{}{
		"secure": "einstein",
	}
	err = lunoClient.Users.Update(findUserbyIDAgain, false, false)
	if err != nil {
		t.Fatal(err)
	}

	// get the user we just updated by id
	findUserbyIDAgainAgain, err := lunoClient.Users.Get(createdUser.ID)
	if err != nil {
		t.Fatal(err)
	}
	if findUserbyIDAgainAgain.Name != findUserbyID.Name {
		t.Errorf("expected username to match, got %s != %s", findUserbyIDAgainAgain.Name, findUserbyID.Name)
	}
	// since we set overwrite to false, old profile info should still be there
	if profileMap, ok := findUserbyIDAgainAgain.Profile.(map[string]interface{}); ok {
		if instrument, ok := profileMap["instrument"].(string); ok {
			if instrument != "trumpet" {
				t.Errorf("expected instrument to be trumpet, got %s", instrument)
			}
		} else {
			t.Errorf("expected instrument to be string, got %T", profileMap["instrument"])
		}
		if secure, ok := profileMap["secure"].(string); ok {
			if secure != "einstein" {
				t.Errorf("expected secure to be einstein, got %s", secure)
			}
		} else {
			t.Errorf("exepcted secure to be a string, got %T", profileMap["secure"])
		}
	} else {
		t.Errorf("expected profile to be map, got %T", findUserbyIDAgainAgain.Profile)
	}

	// now try to deactivate this user
	err = lunoClient.Users.Deactivate(createdUser.ID)
	if err != nil {
		t.Fatal(err)
	}

	// now try to find the user by email
	findUserbyEmail, err := lunoClient.Users.Get("email:" + newUser.Email)
	if err != nil {
		t.Fatal(err)
	}
	// make sure its the right user
	if findUserbyEmail.Name != findUserbyID.Name {
		t.Errorf("expected username to match, got %s != %s", findUserbyEmail.Name, findUserbyID.Name)
	}
	// check that its actually deactivated
	if findUserbyEmail.Closed == "" {
		t.Fatalf("expected deactivated user to have non-empty closed")
	}

	// get the list of users again, should still be 1, even if deactivated
	recentUsers, err = lunoClient.Users.Recent(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(recentUsers.List) != 1 {
		t.Errorf("expected 0 users initially, got %d", len(recentUsers.List))
	}

	// now try to reactivate this user
	err = lunoClient.Users.Reactivate(createdUser.ID)
	if err != nil {
		t.Fatal(err)
	}

	// now try to find the user again by email
	findUserbyEmail, err = lunoClient.Users.Get("email:" + newUser.Email)
	if err != nil {
		t.Fatal(err)
	}
	// make sure its the right user
	if findUserbyEmail.Name != findUserbyID.Name {
		t.Errorf("expected username to match, got %s != %s", findUserbyEmail.Name, findUserbyID.Name)
	}
	// check that its actually deactivated
	if findUserbyEmail.Closed != "" {
		t.Fatalf("expected deactivated user to have empty closed")
	}

	// permanently delete
	err = lunoClient.Users.Delete(createdUser.ID)
	if err != nil {
		t.Fatal(err)
	}

	// get the list of users final time, expecting 0
	recentUsers, err = lunoClient.Users.Recent(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(recentUsers.List) != 0 {
		t.Errorf("expected 0 users initially, got %d", len(recentUsers.List))
	}
}

func TestUserLogin(t *testing.T) {
	apiKey, secretKey, err := lunoKeysFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	lunoClient := NewClient(apiKey, secretKey)

	// now create a user
	user := &User{
		Name:     "Bob Wood",
		Email:    "bob@wood.com",
		Password: "splinterz",
	}
	createdUser, err := lunoClient.Users.Create(user, true, nil)
	if err != nil {
		t.Fatal(err)
	}

	// try to login with wrong password
	_, _, err = lunoClient.Users.LoginWithEmail("bob@wood.com", "wrong", nil, nil)
	if !IsErrorCode(err, ErrCodeIncorrectPassword) {
		t.Errorf("expected error with wrong password, got %v", err)
	}

	// try to login with correct password
	_, _, err = lunoClient.Users.LoginWithEmail("bob@wood.com", "splinterz", nil, nil)
	if err != nil {
		t.Errorf("expected login to succeed, failed with: %v", err)
	}

	// deactivate the user
	err = lunoClient.Users.Deactivate(createdUser.ID)
	if err != nil {
		t.Fatal(err)
	}

	// try to login again, expect fail because deactivated
	_, _, err = lunoClient.Users.LoginWithEmail("bob@wood.com", "splinterz", nil, nil)
	if !IsErrorCode(err, ErrCodeUserClosed) {
		t.Errorf("expected error with wrong password, got %v", err)
	}

	// delete the user
	err = lunoClient.Users.Delete(createdUser.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserDeleteSessions(t *testing.T) {
	apiKey, secretKey, err := lunoKeysFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	lunoClient := NewClient(apiKey, secretKey)

	// now create a user
	user := &User{
		Name:     "Mittens Steve",
		Email:    "mitt@mittens.com",
		Password: "iweargloves",
	}
	createdUser, err := lunoClient.Users.Create(user, true, nil)
	if err != nil {
		t.Fatal(err)
	}

	// try to login with correct password
	_, _, err = lunoClient.Users.LoginWithEmail(user.Email, user.Password, nil, nil)
	if err != nil {
		t.Errorf("expected login to succeed, failed with: %v", err)
	}

	// get sessions for this user
	sessions, err := lunoClient.Sessions.Recent(nil, &SessionFilter{UserID: createdUser.ID}, nil)
	if err != nil {
		t.Fatal(err)
	}

	// check to ensure its 1
	if len(sessions.List) != 1 {
		t.Errorf("expected 1 session, got %d", len(sessions.List))
	}

	// delete sesions for this user
	err = lunoClient.Users.DeleteSessions(createdUser.ID)
	if err != nil {
		t.Fatal(err)
	}

	// get sessions for this user again
	sessions, err = lunoClient.Sessions.Recent(nil, &SessionFilter{UserID: createdUser.ID}, nil)
	if err != nil {
		t.Fatal(err)
	}

	// check to ensure its 0
	if len(sessions.List) != 0 {
		t.Errorf("expected 0 session, got %d", len(sessions.List))
	}

	// cleanup the user
	// delete the user
	err = lunoClient.Users.Delete(createdUser.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserChangePassword(t *testing.T) {
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

	// try to login with correct password
	_, _, err = lunoClient.Users.LoginWithEmail(user.Email, user.Password, nil, nil)
	if err != nil {
		t.Errorf("expected login to succeed, failed with: %v", err)
	}

	// try to change password, with incorrect current password
	err = lunoClient.Users.ChangePassword(createdUser.ID, "imbroke", "notmypassword", true)
	if !IsErrorCode(err, ErrCodeIncorrectPassword) {
		t.Errorf("expected incorrect password, got %v", err)
	}

	// try to change password, with correct current password
	err = lunoClient.Users.ChangePassword(createdUser.ID, "imbroke", user.Password, true)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// try to login with new password
	_, _, err = lunoClient.Users.LoginWithEmail(user.Email, "imbroke", nil, nil)
	if err != nil {
		t.Errorf("expected login to succeed, failed with: %v", err)
	}

	// delete the user
	err = lunoClient.Users.Delete(createdUser.ID)
	if err != nil {
		t.Fatal(err)
	}
}
