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
	"fmt"
	"net/http"
	"net/url"
)

type accountClient struct {
	*Client
}

func (c *accountClient) Get() (*Account, error) {
	resp, err := c.request(http.MethodGet, "/account", nil, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		return ParseAccount(resp)
	}
	return nil, ParseError(resp)
}

func (c *accountClient) Update(account *Account, autoName bool) error {
	params := make(url.Values)
	params.Add("auto_name", fmt.Sprintf("%t", autoName))
	if autoName {
		account.FirstName = ""
		account.LastName = ""
	}
	accountJSON, err := account.MarshalForUpdate()
	if err != nil {
		return fmt.Errorf("error marshaling user json: %v", err)
	}
	resp, err := c.request(http.MethodPut, "/account", params, accountJSON)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return ParseError(resp)
}

func (c *accountClient) Delete(token string) error {
	return ErrNotImplemented
}
