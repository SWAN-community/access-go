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
	"time"

	lru "github.com/hashicorp/golang-lru"
)

// State for access associated with a key
type State struct {
	Allowed     bool      // True if the key was allowed access when last checked, otherwise false
	Count       int       // The number of times the key has been used
	NextRefresh time.Time // The time after which the state needs to be checked again
}

type Cache struct {
	refresh   func(key string, state *State) error // Function to refresh the state
	onEvicted func(key string, state *State)       // Function called when an access key is evicted
	cache     *lru.Cache                           // LRU cache of keys and states
}

// NewAccessProxy creates a new access proxy instance.
func NewCache(
	size int,
	refresh func(string, *State) error,
	onEvicted func(key string, state *State)) (*Cache, error) {
	var err error
	a := &Cache{refresh: refresh, onEvicted: onEvicted}
	if onEvicted != nil {
		a.cache, err = lru.NewWithEvict(size, a.onEvictedInternal)
	} else {
		a.cache, err = lru.New(size)
	}
	if err != nil {
		return nil, err
	}
	return a, nil
}

// GetAllowed returns true if the accessKey is allowed access to the service
// network, otherwise false. If false is returned then the error will
// provide the reason.
func (a *Cache) GetAllowed(accessKey string) (bool, error) {
	var s *State
	t := time.Now()
	i, p := a.cache.Get(accessKey)
	if p {
		// The key has a state instance in the cache. If it needs to be
		// refreshed then update it from the source. If not do nothing.
		s = i.(*State)
		if s.NextRefresh.Before(t) {
			a.refresh(accessKey, s)
		}
	} else {
		// There is no response already in the cache. Create a new one and
		// update it from the source before adding it to the cache.
		s = &State{}
		a.refresh(accessKey, s)
		a.cache.Add(accessKey, s)
	}

	// Increase the number of
	s.Count++
	return s.Allowed, nil
}

// GetAllowedHttp validates the access key in the request and handles any
// failures.
func (a *Cache) GetAllowedHttp(w http.ResponseWriter, r *http.Request) bool {
	return GetAllowedHttp(w, r, a)
}

func (a *Cache) Len() int {
	return a.cache.Len()
}

func (a *Cache) onEvictedInternal(key interface{}, value interface{}) {
	a.onEvicted(key.(string), value.(*State))
}
