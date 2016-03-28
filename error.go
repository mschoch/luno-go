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

// ErrNotImplemented is returned for any API method not yet implemented
var ErrNotImplemented = fmt.Errorf("not implemented")

// Error Codes used by Luno - https://luno.io/docs/errors
const (
	ErrCodeIncorrectPassword = "incorrect_password"
	ErrCodeUserClosed        = "user_closed"
	ErrCodeSessionNotFound   = "session_not_found"
)

// Error represents all the information in a Luno Error
type Error struct {
	Code        string      `json:"code"`
	Message     string      `json:"message"`
	Description string      `json:"description"`
	Status      int         `json:"status"`
	Extra       interface{} `json:"extra"`
}

func (l *Error) Error() string {
	return fmt.Sprintf("LunoError code: %s, message: %s, description: %s, status: %d, extra: %v",
		l.Code, l.Message, l.Description, l.Status, l.Extra)
}

// ParseError parses a Luno error from an HTTP response
func ParseError(resp *http.Response) error {
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading luno error body: %v", err)
	}
	var rv Error
	err = json.Unmarshal(respBytes, &rv)
	if err != nil {
		return fmt.Errorf("error parsing luno error json: '%s' err: %v", respBytes, err)
	}
	return &rv
}

// IsErrorCode checks is an error was a luno error with the matching code
func IsErrorCode(err error, code string) bool {
	if err, ok := err.(*Error); ok {
		if err.Code == code {
			return true
		}
	}
	return false
}
