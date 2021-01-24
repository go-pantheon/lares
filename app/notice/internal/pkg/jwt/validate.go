package jwt

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var (
	defaultValidator = &Validator{certGetter: newCacheCertGetter(http.DefaultClient)}
	now              = time.Now
)

type Payload struct {
	Issuer   string                 `json:"iss"`
	Audience string                 `json:"aud"`
	Expires  int64                  `json:"exp"`
	IssuedAt int64                  `json:"iat"`
	Subject  string                 `json:"sub,omitempty"`
	Claims   map[string]interface{} `json:"-"`
}

type jwt struct {
	header    string
	payload   string
	signature string
}

type jwtHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
	KeyID     string `json:"kid"`
}

type certResponse struct {
	Keys []jwk `json:"keys"`
}

type jwk struct {
	Alg string `json:"alg"`
	Crv string `json:"crv"`
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Use string `json:"use"`
	E   string `json:"e"`
	N   string `json:"n"`
	X   string `json:"x"`
	Y   string `json:"y"`
}

type Validator struct {
	certGetter *cacheCertGetter
	saCertsURL string
}

func NewValidator(ctx context.Context, saCertURL string, client *http.Client) *Validator {
	return &Validator{saCertsURL: saCertURL, certGetter: newCacheCertGetter(client)}
}

func (v *Validator) Validate(ctx context.Context, idToken string, audience string) (*Payload, error) {
	jwt, err := parseJWT(idToken)
	if err != nil {
		return nil, err
	}
	header, err := jwt.parsedHeader()
	if err != nil {
		return nil, err
	}
	payload, err := jwt.parsedPayload()
	if err != nil {
		return nil, err
	}
	sig, err := jwt.decodedSignature()
	if err != nil {
		return nil, err
	}

	if audience != "" && payload.Audience != audience {
		return nil, errors.Errorf("idtoken: audience provided does not match aud claim in the JWT")
	}

	if now().Unix() > payload.Expires {
		return nil, errors.Errorf("idtoken: token expired")
	}

	switch header.Algorithm {
	case "RS256":
		if err := v.validateRS256(ctx, header.KeyID, jwt.hashedContent(), sig); err != nil {
			return nil, err
		}
	default:
		return nil, errors.Errorf("idtoken: expected JWT signed with RS256 but found %q", header.Algorithm)
	}

	return payload, nil
}

func (v *Validator) validateRS256(ctx context.Context, keyID string, hashedContent []byte, sig []byte) error {
	certResp, err := v.certGetter.getCert(ctx, v.saCertsURL)
	if err != nil {
		return err
	}
	j, err := findMatchingKey(certResp, keyID)
	if err != nil {
		return err
	}
	dn, err := decode(j.N)
	if err != nil {
		return err
	}
	de, err := decode(j.E)
	if err != nil {
		return err
	}

	pk := &rsa.PublicKey{
		N: new(big.Int).SetBytes(dn),
		E: int(new(big.Int).SetBytes(de).Int64()),
	}
	return rsa.VerifyPKCS1v15(pk, crypto.SHA256, hashedContent, sig)
}

func findMatchingKey(response *certResponse, keyID string) (*jwk, error) {
	if response == nil {
		return nil, errors.Errorf("idtoken: cert response is nil")
	}
	for _, v := range response.Keys {
		if v.Kid == keyID {
			return &v, nil
		}
	}
	return nil, errors.Errorf("idtoken: could not find matching cert keyId for the token provided")
}

func parseJWT(idToken string) (*jwt, error) {
	segments := strings.Split(idToken, ".")
	if len(segments) != 3 {
		return nil, errors.Errorf("idtoken: invalid token, token must have three segments; found %d", len(segments))
	}
	return &jwt{
		header:    segments[0],
		payload:   segments[1],
		signature: segments[2],
	}, nil
}

// decodedHeader base64 decodes the header segment.
func (j *jwt) decodedHeader() ([]byte, error) {
	dh, err := decode(j.header)
	if err != nil {
		return nil, errors.Errorf("idtoken: unable to decode JWT header: %v", err)
	}
	return dh, nil
}

// decodedPayload base64 payload the header segment.
func (j *jwt) decodedPayload() ([]byte, error) {
	p, err := decode(j.payload)
	if err != nil {
		return nil, errors.Errorf("idtoken: unable to decode JWT payload: %v", err)
	}
	return p, nil
}

// decodedPayload base64 payload the header segment.
func (j *jwt) decodedSignature() ([]byte, error) {
	p, err := decode(j.signature)
	if err != nil {
		return nil, errors.Errorf("idtoken: unable to decode JWT signature: %v", err)
	}
	return p, nil
}

// parsedHeader returns a struct representing a JWT header.
func (j *jwt) parsedHeader() (jwtHeader, error) {
	var h jwtHeader
	dh, err := j.decodedHeader()
	if err != nil {
		return h, err
	}
	err = json.Unmarshal(dh, &h)
	if err != nil {
		return h, errors.Errorf("idtoken: unable to unmarshal JWT header: %v", err)
	}
	return h, nil
}

// parsedPayload returns a struct representing a JWT payload.
func (j *jwt) parsedPayload() (*Payload, error) {
	var p Payload
	dp, err := j.decodedPayload()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(dp, &p); err != nil {
		return nil, errors.Errorf("idtoken: unable to unmarshal JWT payload: %v", err)
	}
	if err := json.Unmarshal(dp, &p.Claims); err != nil {
		return nil, errors.Errorf("idtoken: unable to unmarshal JWT payload claims: %v", err)
	}
	return &p, nil
}

// hashedContent gets the SHA256 checksum for verification of the JWT.
func (j *jwt) hashedContent() []byte {
	signedContent := j.header + "." + j.payload
	hashed := sha256.Sum256([]byte(signedContent))
	return hashed[:]
}

func (j *jwt) String() string {
	return fmt.Sprintf("%s.%s.%s", j.header, j.payload, j.signature)
}

func decode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}
