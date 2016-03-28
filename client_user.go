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
	"net/http"
	"net/url"
)

type usersClient struct {
	*Client
}

func (c *usersClient) Recent(expand []string, paging *Paging) (*Users, error) {
	params := paging.Params()
	for _, item := range expand {
		params.Add("expand", item)
	}
	resp, err := c.request(http.MethodGet, "/users", params, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		return ParseUsers(resp)
	}
	return nil, ParseError(resp)
}

func (c *usersClient) Create(user *User, autoName bool, expand []string) (*User, error) {
	params := make(url.Values)
	for _, item := range expand {
		params.Add("expand", item)
	}
	params.Add("auto_name", fmt.Sprintf("%t", autoName))
	userJSON, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("error marshaling user json: %v", err)
	}
	resp, err := c.request(http.MethodPost, "/users", params, userJSON)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusCreated {
		return ParseUser(resp)
	}
	return nil, ParseError(resp)
}

func (c *usersClient) Update(user *User, autoName bool, overwriteProfile bool) error {
	params := make(url.Values)
	params.Add("auto_name", fmt.Sprintf("%t", autoName))
	method := http.MethodPatch
	if overwriteProfile {
		method = http.MethodPut
	}
	if autoName {
		user.FirstName = ""
		user.LastName = ""
	}
	userJSON, err := user.MarshalForUpdate()
	if err != nil {
		return fmt.Errorf("error marshaling user json: %v", err)
	}
	resp, err := c.request(method, "/users/"+user.ID, params, userJSON)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return ParseError(resp)
}

func (c *usersClient) delete(id string, permanent bool) error {
	params := make(url.Values)
	params.Add("permanent", fmt.Sprintf("%t", permanent))
	resp, err := c.request(http.MethodDelete, "/users/"+id, params, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return ParseError(resp)
}

func (c *usersClient) Delete(id string) error {
	return c.delete(id, true)
}

func (c *usersClient) Deactivate(id string) error {
	return c.delete(id, false)
}

func (c *usersClient) Reactivate(id string) error {
	resp, err := c.request(http.MethodPost, "/users/"+id+"/reactivate", nil, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return ParseError(resp)
}

func (c *usersClient) Get(id string) (*User, error) {
	resp, err := c.request(http.MethodGet, "/users/"+id, nil, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		return ParseUser(resp)
	}
	return nil, ParseError(resp)
}

func (c *usersClient) login(expand []string, login *Login) (*User, *Session, error) {
	params := make(url.Values)
	for _, item := range expand {
		params.Add("expand", item)
	}
	loginJSON, err := json.Marshal(login)
	if err != nil {
		return nil, nil, fmt.Errorf("error marshaling login json: %v", err)
	}
	resp, err := c.request(http.MethodPost, "/users/login", params, loginJSON)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode == http.StatusOK {
		return ParseLogin(resp)
	}
	return nil, nil, ParseError(resp)
}

func (c *usersClient) LoginWithID(id, password string, expand []string, session *Session) (*User, *Session, error) {
	return c.login(expand, &Login{
		ID:       id,
		Password: password,
		Session:  session,
	})
}

func (c *usersClient) LoginWithEmail(email, password string, expand []string, session *Session) (*User, *Session, error) {
	return c.login(expand, &Login{
		Email:    email,
		Password: password,
		Session:  session,
	})
}

func (c *usersClient) LoginWithUsername(username, password string, expand []string, session *Session) (*User, *Session, error) {
	return c.login(expand, &Login{
		Username: username,
		Password: password,
		Session:  session,
	})
}

func (c *usersClient) LoginWithAny(login, password string, expand []string, session *Session) (*User, *Session, error) {
	return c.login(expand, &Login{
		Login:    login,
		Password: password,
		Session:  session,
	})
}

func (c *usersClient) DeleteSessions(id string) error {
	resp, err := c.request(http.MethodDelete, "/users/"+id+"/sessions", nil, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return ParseError(resp)
}

func (c *usersClient) ValidatePassword(id, password string) error {
	validate := map[string]interface{}{
		"password": password,
	}
	validateJSON, err := json.Marshal(validate)
	resp, err := c.request(http.MethodPost, "/users/"+id+"/password/validate", nil, validateJSON)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return ParseError(resp)
}

func (c *usersClient) ChangePassword(id, newPassword, currentPassword string, requireCurrent bool) error {
	params := make(url.Values)
	params.Add("require_current_password", fmt.Sprintf("%t", requireCurrent))
	change := map[string]interface{}{
		"password": newPassword,
	}
	if requireCurrent {
		change["current_password"] = currentPassword
	}
	changeJSON, err := json.Marshal(change)
	resp, err := c.request(http.MethodPost, "/users/"+id+"/password/change", params, changeJSON)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return ParseError(resp)
}
