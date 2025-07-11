package jwt

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

const (
	maxResponseSize = 1024 * 1024 // 1MB
)

type cacheCertGetter struct {
	client *http.Client

	// clock optionally specifies a func to return the current time.
	// If nil, time.Now is used.
	clock func() time.Time

	mu    sync.Mutex
	certs map[string]*cachedResponse
}

func newCacheCertGetter(client *http.Client) *cacheCertGetter {
	return &cacheCertGetter{
		client: client,
		certs:  make(map[string]*cachedResponse, 2),
	}
}

type cachedResponse struct {
	resp *certResponse
	exp  time.Time
}

func (c *cacheCertGetter) requestCert(ctx context.Context, url string) (*certResponse, error) {
	if response, ok := c.get(url); ok {
		return response, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Errorf("close apple auth token response body failed: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("idtoken: unable to retrieve cert, got status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		return nil, err
	}

	certResp := &certResponse{}
	if err := json.Unmarshal(body, certResp); err != nil {
		return nil, err
	}

	c.set(url, certResp, resp.Header)

	return certResp, nil
}

func (c *cacheCertGetter) now() time.Time {
	if c.clock != nil {
		return c.clock()
	}

	return time.Now()
}

func (c *cacheCertGetter) get(url string) (*certResponse, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	cachedResp, ok := c.certs[url]
	if !ok {
		return nil, false
	}

	if c.now().After(cachedResp.exp) {
		return nil, false
	}

	return cachedResp.resp, true
}

func (c *cacheCertGetter) set(url string, resp *certResponse, headers http.Header) {
	exp := c.calculateExpireTime(headers)

	c.mu.Lock()
	defer c.mu.Unlock()

	c.certs[url] = &cachedResponse{resp: resp, exp: exp}
}

// calculateExpireTime will determine the expire time for the cache based on
// HTTP headers. If there is any difficulty reading the headers the fallback is
// to set the cache to expire now.
func (c *cacheCertGetter) calculateExpireTime(headers http.Header) time.Time {
	var maxAge int

	cc := strings.SplitSeq(headers.Get("cache-control"), ",")
	for v := range cc {
		if strings.Contains(v, "max-age") {
			ss := strings.Split(v, "=")
			if len(ss) < 2 {
				return c.now()
			}

			ma, err := strconv.Atoi(ss[1])
			if err != nil {
				return c.now()
			}

			maxAge = ma
		}
	}

	a := headers.Get("age")
	if a == "" {
		return c.now().Add(time.Duration(maxAge) * time.Second)
	}

	age, err := strconv.Atoi(a)
	if err != nil {
		return c.now()
	}

	return c.now().Add(time.Duration(maxAge-age) * time.Second)
}
