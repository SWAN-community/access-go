/* ****************************************************************************
 * Copyright 2022 51 Degrees Mobile Experts Limited (51degrees.com)
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 * ***************************************************************************/

package access

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/SWAN-community/common-go"
)

func TestAccessAllowed(t *testing.T) {
	a := NewFixed([]string{"A"})
	t.Run("allowed", func(t *testing.T) {
		rr := testAccessHttp(t, a, "A")
		if rr.Code != http.StatusOK {
			t.Fatal("access service should allow")
		}
	})
	t.Run("disallowed", func(t *testing.T) {
		rr := testAccessHttp(t, a, "B")
		if rr.Code != http.StatusNetworkAuthenticationRequired {
			t.Fatal("access service should disallow")
		}
	})
}

// testAccessHttp calls the handler with the key provided and returns the result
// from the service.
func testAccessHttp(
	t *testing.T,
	service Access,
	key string) *httptest.ResponseRecorder {
	v := url.Values{}
	v.Add("accessKey", key)
	return common.HTTPTest(
		t,
		"GET",
		"",
		"/test",
		v,
		func(w http.ResponseWriter, r *http.Request) {
			// Send an empty string to generate a 200 status good if the access
			// key is allowed. If the result is false then the method will
			// already have responded.
			if service.GetAllowedHttp(w, r) {
				common.SendString(w, "")
			}
		})
}
