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

// Account represents the Luno Account - https://luno.io/docs#account
type Account struct {
	Entity
	Email     string `json:"email"`
	Name      string `json:"name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
	Created   string `json:"created"`
	Closed    string `json:"closed"`
}

// MarshalForUpdate exports only those fields suitable for an update operation
func (a *Account) MarshalForUpdate() ([]byte, error) {
	tmp := map[string]interface{}{
		"email":      a.Email,
		"name":       a.Name,
		"first_name": a.FirstName,
		"last_name":  a.LastName,
	}
	return json.Marshal(tmp)
}

// ParseAccount parses an Account out of an HTTP response
func ParseAccount(resp *http.Response) (*Account, error) {
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading luno response body: %v", err)
	}
	var rv Account
	err = json.Unmarshal(respBytes, &rv)
	if err != nil {
		return nil, fmt.Errorf("error parsing luno account json: '%s' err: %v", respBytes, err)
	}
	return &rv, nil
}
