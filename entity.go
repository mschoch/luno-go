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
	"net/url"
)

// Entity contains fields common to many Luno entities
type Entity struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
	URL  string `json:"url,omitempty"`
}

// Page contains fields in results related to paging
type Page struct {
	Next Entity `json:"next"`
	Prev Entity `json:"prev"`
}

// Paging contains request options related to paging
type Paging struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Limit int    `json:"limit"`
}

// Params converts Paging struct into HTTP request params
func (p *Paging) Params() url.Values {
	rv := make(url.Values)
	if p != nil {
		if p.From != "" {
			rv.Add("from", p.From)
		}
		if p.To != "" {
			rv.Add("to", p.To)
		}
		rv.Add("limit", fmt.Sprintf("%d", p.Limit))
	}
	return rv
}
