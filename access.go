/* ****************************************************************************
 * Copyright 2020 51 Degrees Mobile Experts Limited (51degrees.com)
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

	"github.com/SWAN-community/common-go"
)

// Access interface for validating entitlement to access the network.
type Access interface {

	// GetAllowed returns true if the accessKey is allowed access to the service
	// network, otherwise false. If false is returned then the error will
	// provide the reason.
	GetAllowed(accessKey string) (bool, error)

	// Returns true if the request is allowed to access the handler, otherwise
	// false. Removes the accessKey parameter from the form to prevent it being
	// used by other methods. If false is returned then no further action is
	// needed as the method will have responded to the request already.
	GetAllowedHttp(w http.ResponseWriter, r *http.Request) bool
}

// GetAllowedHttp of the Access.GetAllowedHttp method for use my implementing
// structures.
func GetAllowedHttp(
	w http.ResponseWriter,
	r *http.Request,
	a Access) bool {
	err := r.ParseForm()
	if err != nil {
		common.ReturnServerError(w, err)
		return false
	}
	k, ok := r.Form["accessKey"]
	if !ok || len(k) != 1 {
		common.ReturnApplicationError(w, &common.HttpError{
			Log:     false,
			Request: r,
			Code:    http.StatusBadRequest,
			Message: "accessKey missing",
			Error:   err})
		return false
	}
	v, err := a.GetAllowed(k[0])
	if !v || err != nil {
		common.ReturnApplicationError(w, &common.HttpError{
			Log:     false,
			Request: r,
			Code:    http.StatusNetworkAuthenticationRequired,
			Message: "accessKey invalid",
			Error:   err})
		return false
	}
	r.Form.Del("accessKey")
	return true
}
