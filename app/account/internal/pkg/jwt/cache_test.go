// Copyright 2020 Google LLC.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jwt

import (
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testAppleSACertsURL string = "https://appleid.apple.com/auth/keys"
)

type fakeClock struct {
	mu sync.Mutex
	t  time.Time
}

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.t
}

func (c *fakeClock) Sleep(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.t = c.t.Add(d)
}

func TestCacheHit(t *testing.T) {
	t.Parallel()

	value := "123"
	clock := &fakeClock{t: time.Now()}
	dummyResp := &certResponse{
		Keys: []jwk{
			{
				Kid: value,
			},
		},
	}

	cache := newCacheCertGetter(nil)
	cache.clock = clock.Now

	// Cache should be empty
	cert, ok := cache.get(testAppleSACertsURL)
	assert.False(t, ok)
	assert.Nil(t, cert)

	// Add an item, but make it expire now
	cache.set(testAppleSACertsURL, dummyResp, make(http.Header))
	clock.Sleep(time.Nanosecond) // it expires when current time is > expiration, not >=

	cert, ok = cache.get(testAppleSACertsURL)
	assert.False(t, ok)
	assert.Nil(t, cert)

	// Add an item that expires in 1 seconds
	h := make(http.Header)
	h.Set("age", "0")
	h.Set("cache-control", "public, max-age=1, must-revalidate, no-transform")

	cache.set(testAppleSACertsURL, dummyResp, h)

	cert, ok = cache.get(testAppleSACertsURL)
	assert.True(t, ok)
	assert.NotNil(t, cert)
	assert.Equal(t, value, cert.Keys[0].Kid)

	// Wait
	clock.Sleep(2 * time.Second)

	cert, ok = cache.get(testAppleSACertsURL)
	assert.False(t, ok)
	assert.Nil(t, cert)
}
