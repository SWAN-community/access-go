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
	"testing"
	"time"
)

type testCache struct {
	Cache     *Cache // Access cache under test
	evictions int    // Number of evictions
	calls     int    // Number of calls to refresh
}

// TestCache tests the cached implementation for allowed state.
func TestCache(t *testing.T) {
	t.Run("key allowed", func(t *testing.T) {
		testCacheChecks(t, testNewCache(t, 1), "A", true, 1, 0, 1)
	})
	t.Run("key not allowed", func(t *testing.T) {
		testCacheChecks(t, testNewCache(t, 1), "B", false, 1, 0, 1)
	})
	t.Run("evicted", func(t *testing.T) {
		a := testNewCache(t, 2)
		testCacheChecks(t, a, "A", true, 1, 0, 1)
		testCacheChecks(t, a, "B", false, 2, 0, 2)
		testCacheChecks(t, a, "C", false, 3, 1, 2)
	})
}

// testCache performs a GetAllowed operation on the access cache and
// checks the result, calls, and evictions match the parameters provided.
// a test access cache which tracks calls and evictions
// key to be checked
// allowed the expected result
// calls the expected calls to the cache
// evictions the expected number of evictions
// len the expected length of the cache
func testCacheChecks(
	t *testing.T,
	a *testCache,
	key string,
	allowed bool,
	calls int,
	evictions int,
	len int) {
	v, _ := a.Cache.GetAllowed(key)
	if v != allowed {
		t.Fatalf("Excepted '%v', got '%v'", allowed, v)
	}
	if a.calls != calls {
		t.Fatalf("Excepted '%v' call, got '%v'", calls, a.calls)
	}
	if a.evictions != evictions {
		t.Fatalf("Excepted '%v' evictions, got '%v'", evictions, a.evictions)
	}
	if a.Cache.Len() != len {
		t.Fatalf("Excepted '%v' cache length, got '%v'",
			len,
			a.Cache.Len())
	}
}

func testNewCache(t *testing.T, size int) *testCache {
	var err error
	r := &testCache{}
	r.Cache, err = NewCache(
		size,
		func(key string, state *State) error {
			state.Allowed = key == "A"
			state.NextRefresh = time.Now().Add(time.Minute)
			state.Count = 0
			r.calls++
			return nil
		},
		func(key string, state *State) {
			r.evictions++
		})
	if err != nil {
		t.Error(err)
	}
	return r
}
