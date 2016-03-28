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

func TestAccount(t *testing.T) {

	apiKey, secretKey, err := lunoKeysFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	lunoClient := NewClient(apiKey, secretKey)

	// get account info
	account, err := lunoClient.Account.Get()
	if err != nil {
		t.Fatal(err)
	}

	originalName := account.Name

	// try to update it
	account.Name += "-updated"
	err = lunoClient.Account.Update(account, false)
	if err != nil {
		t.Fatal(err)
	}

	// get account info again
	account, err = lunoClient.Account.Get()
	if err != nil {
		t.Fatal(err)
	}
	if account.Name != originalName+"-updated" {
		t.Errorf("expected '%s', got '%s'", originalName+"-updated", account.Name)
	}

	// try to reset it
	account.Name = originalName
	err = lunoClient.Account.Update(account, false)
	if err != nil {
		t.Fatal(err)
	}

}
