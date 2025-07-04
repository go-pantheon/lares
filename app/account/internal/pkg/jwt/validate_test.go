// Copyright 2020 Google LLC.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jwt

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	keyID        = "1234"
	testAudience = "test-audience"
)

func TestValidateRS256(t *testing.T) {
	t.Parallel()

	idToken, pk := createRS256JWT(t)
	tests := []struct {
		name    string
		keyID   string
		n       *big.Int
		e       int
		wantErr bool
	}{
		{
			name:    "works",
			keyID:   keyID,
			n:       pk.N,
			e:       pk.E,
			wantErr: false,
		},
		{
			name:    "no matching key",
			keyID:   "5678",
			n:       pk.N,
			e:       pk.E,
			wantErr: true,
		},
		{
			name:    "sig does not match",
			keyID:   keyID,
			n:       new(big.Int).SetBytes([]byte("42")),
			e:       42,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := &http.Client{
				Transport: RoundTripFn(func(req *http.Request) *http.Response {
					cr := certResponse{
						Keys: []jwk{
							{
								Kid: tt.keyID,
								N:   base64.RawURLEncoding.EncodeToString(tt.n.Bytes()),
								E:   base64.RawURLEncoding.EncodeToString(new(big.Int).SetInt64(int64(tt.e)).Bytes()),
							},
						},
					}

					b, err := json.Marshal(&cr)
					require.NoError(t, err)

					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewReader(b)),
						Header:     make(http.Header),
					}
				}),
			}

			v := NewValidator(context.Background(), testAppleSACertsURL, client)
			payload, err := v.Validate(context.Background(), idToken, testAudience)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, payload)
			assert.Equal(t, testAudience, payload.Audience)
			assert.NotEmpty(t, payload.Claims)

			got, ok := payload.Claims["aud"]
			require.True(t, ok)

			got, ok = got.(string)
			require.True(t, ok)
			assert.Equal(t, testAudience, got)
		})
	}
}

func createRS256JWT(t *testing.T) (string, rsa.PublicKey) {
	t.Helper()

	token := commonToken(t, "RS256")

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	sig, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, token.hashedContent())
	require.NoError(t, err)

	token.signature = base64.RawURLEncoding.EncodeToString(sig)

	return token.String(), privateKey.PublicKey
}

func commonToken(t *testing.T, alg string) *jwt {
	t.Helper()

	header := jwtHeader{
		KeyID:     keyID,
		Algorithm: alg,
		Type:      "JWT",
	}

	payload := Payload{
		Issuer:   "example.com",
		Audience: testAudience,
		Expires:  time.Now().Unix() + 1000,
	}

	hb, err := json.Marshal(&header)
	require.NoError(t, err)

	pb, err := json.Marshal(&payload)
	require.NoError(t, err)

	return &jwt{
		header:  base64.RawURLEncoding.EncodeToString(hb),
		payload: base64.RawURLEncoding.EncodeToString(pb),
	}
}

type RoundTripFn func(req *http.Request) *http.Response

func (f RoundTripFn) RoundTrip(req *http.Request) (*http.Response, error) { return f(req), nil }
