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

// Users represents a list of User objects
type Users struct {
	Entity
	List []*User `json:"list"`
	Page Page    `json:"page"`
}

// ParseUsers parses Users out of an HTTP response
func ParseUsers(resp *http.Response) (*Users, error) {
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading luno response body: %v", err)
	}
	var rv Users
	err = json.Unmarshal(respBytes, &rv)
	if err != nil {
		return nil, fmt.Errorf("error parsing luno users json: '%s' err: %v", respBytes, err)
	}
	return &rv, nil
}

// User represents a Luno User - https://luno.io/docs#users
type User struct {
	Entity
	Email     string      `json:"email,omitempty"`
	UserName  string      `json:"username,omitempty"`
	Name      string      `json:"name,omitempty"`
	FirstName string      `json:"first_name,omitempty"`
	LastName  string      `json:"last_name,omitempty"`
	Created   string      `json:"created,omitempty"`
	Closed    string      `json:"closed,omitempty"`
	Password  string      `json:"password,omitempty"`
	Profile   interface{} `json:"profile,omitempty"`
}

// ParseUser parses a User out of an HTTP response
func ParseUser(resp *http.Response) (*User, error) {
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading luno response body: %v", err)
	}
	var rv User
	err = json.Unmarshal(respBytes, &rv)
	if err != nil {
		return nil, fmt.Errorf("error parsing luno user json: '%s' err: %v", respBytes, err)
	}
	return &rv, nil
}

// ParseLogin parses a login response (User and Session) from an HTTP response
func ParseLogin(resp *http.Response) (*User, *Session, error) {
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading luno response body: %v", err)
	}
	var rv struct {
		User    *User    `json:"user"`
		Session *Session `json:"session"`
	}
	err = json.Unmarshal(respBytes, &rv)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing luno user json: '%s' err: %v", respBytes, err)
	}

	return rv.User, rv.Session, nil
}

// MarshalForUpdate exports only the User fields suitable for an update operation
func (u *User) MarshalForUpdate() ([]byte, error) {
	tmp := map[string]interface{}{
		"email":      u.Email,
		"username":   u.UserName,
		"name":       u.Name,
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"profile":    u.Profile,
	}
	return json.Marshal(tmp)
}
