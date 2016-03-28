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

// Sessions represents a list of Session objects
type Sessions struct {
	Entity
	List []*Session `json:"list"`
	Page Page       `json:"page"`
}

// ParseSessions parses Sessions out of an HTTP response
func ParseSessions(resp *http.Response) (*Sessions, error) {
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading luno response body: %v", err)
	}
	var rv Sessions
	err = json.Unmarshal(respBytes, &rv)
	if err != nil {
		return nil, fmt.Errorf("error parsing luno sessions json: '%s' err: %v", respBytes, err)
	}
	return &rv, nil
}

// SessionFilter lets you filter session
type SessionFilter struct {
	UserID string `json:"user_id"`
}

// Session represents a Luno session - https://luno.io/docs#sessions
type Session struct {
	Entity
	UserID      string      `json:"user_id,omitempty"`
	Key         string      `json:"key,omitempty"`
	Created     string      `json:"created,omitempty"`
	Expires     string      `json:"expires,omitempty"`
	LastAccess  string      `json:"last_access,omitempty"`
	AccessCount int         `json:"access_count,omitempty"`
	IP          string      `json:"ip,omitempty"`
	UserAgent   string      `json:"user_agent,omitempty"`
	Details     interface{} `json:"details,omitempty"`
	User        *User       `json:"user,omitempty"`
}

// MarshalForUpdate exports only the Session fields suitable for an update operation
func (s *Session) MarshalForUpdate(includeKey bool) ([]byte, error) {
	tmp := map[string]interface{}{}
	if includeKey && s.Key != "" {
		tmp["key"] = s.Key
	}
	if s.UserID != "" {
		tmp["user_id"] = s.UserID
	}
	if s.Expires != "" {
		tmp["expires"] = s.Expires
	}
	if s.IP != "" {
		tmp["ip"] = s.IP
	}
	if s.UserAgent != "" {
		tmp["user_agent"] = s.UserAgent
	}
	if s.Details != nil {
		tmp["details"] = s.Details
	}

	return json.Marshal(tmp)
}

// ParseSession parses a Session out of an HTTP response
func ParseSession(resp *http.Response) (*Session, error) {
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading luno response body: %v", err)
	}
	var rv Session
	err = json.Unmarshal(respBytes, &rv)
	if err != nil {
		return nil, fmt.Errorf("error parsing luno session json: '%s' err: %v", respBytes, err)
	}
	return &rv, nil
}

// Login represents information used to log in to the system
type Login struct {
	ID       string   `json:"id,omitempty"`
	Email    string   `json:"email,omitempty"`
	Username string   `json:"username,omitempty"`
	Login    string   `json:"username,omitempty"`
	Password string   `json:"password,omitempty"`
	Session  *Session `json:"session,omitempty"`
}

// MarshalJSON converts a Login to JSON
func (l *Login) MarshalJSON() ([]byte, error) {
	session := make(map[string]interface{})
	if l.Session != nil {
		if l.Session.Expires != "" {
			session["expires"] = l.Session.Expires
		}
		if l.Session.IP != "" {
			session["ip"] = l.Session.IP
		}
		if l.Session.UserAgent != "" {
			session["user_agent"] = l.Session.UserAgent
		}
		if l.Session.Details != nil {
			session["details"] = l.Session.Details
		}
	}
	tmp := make(map[string]interface{})
	if l.ID != "" {
		tmp["id"] = l.ID
	}
	if l.Email != "" {
		tmp["email"] = l.Email
	}
	if l.Username != "" {
		tmp["username"] = l.Username
	}
	if l.Login != "" {
		tmp["login"] = l.Login
	}
	if l.Password != "" {
		tmp["password"] = l.Password
	}
	tmp["session"] = session
	return json.Marshal(tmp)
}
