package domain

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/fabrica-util/xtime"
	"github.com/go-pantheon/lares/app/account/internal/conf"
	xjwt "github.com/go-pantheon/lares/app/account/internal/pkg/jwt"
	"github.com/go-pantheon/lares/app/account/internal/pkg/security"
	v1 "github.com/go-pantheon/lares/gen/api/server/account/interface/account/v1"
	"github.com/golang-jwt/jwt/v4"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/proto"
)

const (
	maxResponseSize  int64         = 1024 * 1024 // 1MB
	appleTokenExpire time.Duration = 5 * time.Minute
)

type AppleDomain struct {
	log  *log.Helper
	repo AccountRepo

	conf      *conf.Apple
	validator *xjwt.Validator

	cli    *http.Client
	priKey *ecdsa.PrivateKey
}

func NewAppleDomain(logger log.Logger, label *conf.Label, c *conf.Platform, repo AccountRepo) (do *AppleDomain, err error) {
	do = &AppleDomain{
		log:  log.NewHelper(log.With(logger, "module", "account/http/domain/apple")),
		conf: proto.Clone(c.Apple).(*conf.Apple),
		repo: repo,
	}

	do.cli = &http.Client{
		Timeout: profile.ClientTimeout,
		Transport: &http.Transport{
			MaxIdleConns:        profile.ClientMaxIdleConns,
			MaxIdleConnsPerHost: profile.ClientMaxIdleConnsPerHost,
			IdleConnTimeout:     profile.ClientIdleConnTimeout,
			TLSHandshakeTimeout: profile.ClientTLSHandshakeTimeout,
			DisableKeepAlives:   false,
		},
	}
	do.validator = xjwt.NewValidator(context.Background(), c.Apple.AppleSaCertsUrl, do.cli)

	if profile.IsDevStr(label.Profile) {
		do.priKey, _ = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		do.log.Infof("use mock apple login client")

		return do, nil
	}

	do.priKey, err = jwt.ParseECPrivateKeyFromPEM([]byte(do.conf.Secret))
	if err != nil {
		return nil, err
	}

	return do, nil
}

func (do *AppleDomain) GetOrCreateAccount(ctx context.Context, appleId string, ip string) (acc *Account, isCreated bool, err error) {
	acc, err = do.repo.GetByApple(ctx, appleId)
	if err == nil {
		return acc, false, nil
	}

	if !errors.IsNotFound(err) {
		return nil, false, err
	}

	param := &Account{
		Channel:     int32(ChannelTypeApple),
		Apple:       appleId,
		RegisterIp:  ip,
		LastLoginIp: ip,
	}

	acc, err = do.repo.Create(ctx, param)
	if err != nil {
		return nil, false, err
	}

	return acc, true, nil
}

func (do *AppleDomain) GetAccount(ctx context.Context, appleId string) (*Account, error) {
	return do.repo.GetByApple(ctx, appleId)
}

func (do *AppleDomain) CreateAccount(ctx context.Context, appleId string, ip string) (*Account, error) {
	param := &Account{
		Channel:     int32(ChannelTypeApple),
		Apple:       appleId,
		RegisterIp:  ip,
		LastLoginIp: ip,
	}

	return do.repo.Create(ctx, param)
}

const (
	AppleLoginTypeApp = iota
	AppleLoginTypeWeb
)

func (do *AppleDomain) RequestPlatformIdByAppToken(ctx context.Context, token string) (appleId string, err error) {
	return do.requestPlatformId(ctx, token, do.conf.AudApp)
}

func (do *AppleDomain) RequestPlatformIdByWebToken(ctx context.Context, token string) (appleId string, err error) {
	return do.requestPlatformId(ctx, token, do.conf.AudWeb)
}

func (do *AppleDomain) requestPlatformId(ctx context.Context, token string, aud string) (appleId string, err error) {
	var info *xjwt.Payload

	info, err = do.validator.Validate(ctx, token, aud)
	if err != nil {
		return "", xerrors.APIPlatformAuthFailed("validate apple token failed")
	}

	return info.Subject, nil
}

func (do *AppleDomain) RequestToken(ctx context.Context, code string) (string, error) {
	secret, err := do.buildSecret()
	if err != nil {
		return "", xerrors.APIPlatformAuthFailed("build apple auth token secret failed")
	}

	form := url.Values{}
	form.Set("client_id", do.conf.ClientId)
	form.Set("client_secret", secret)
	form.Set("code", code)
	form.Set("grant_type", "authorization_code")
	form.Set("redirect_uri", do.conf.RedirectUri)

	body := bytes.NewBufferString(form.Encode())

	req, err := http.NewRequestWithContext(ctx, "POST", do.conf.AppleAuthTokenUrl, body)
	if err != nil {
		return "", xerrors.APIPlatformAuthFailed("create apple auth token request failed")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := do.cli.Do(req)
	if err != nil {
		return "", xerrors.APIPlatformAuthFailed("do apple auth token request failed")
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			do.log.Errorf("close apple auth token response body failed: %v", err)
		}
	}()

	data, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		return "", xerrors.APIPlatformAuthFailed("read apple auth token response body failed")
	}

	if resp.StatusCode != http.StatusOK {
		return "", xerrors.APIPlatformAuthFailed("apple auth token request failed")
	}

	tr := &tokenResp{}
	if err = jsoniter.Unmarshal(data, tr); err != nil {
		return "", xerrors.APIPlatformAuthFailed("unmarshal apple auth token response body failed")
	}

	return tr.IdToken, nil
}

type tokenResp struct {
	IdToken string `json:"id_token"`
}

func (do *AppleDomain) buildSecret() (k string, err error) {
	token := &jwt.Token{
		Header: map[string]interface{}{
			"alg": "ES256",
			"kid": do.conf.KeyId,
		},
		Claims: jwt.MapClaims{
			"iss": do.conf.TeamId,
			"iat": time.Now().Unix() - 10, // ahead of 10 seconds
			"exp": time.Now().Add(time.Hour).Unix(),
			"aud": do.conf.Validator,
			"sub": do.conf.ClientId,
		},
		Method: jwt.SigningMethodES256,
	}

	return token.SignedString(do.priKey)
}

func (do *AppleDomain) CheckAppleId(ctx context.Context, secret string) (*v1.AppleId, error) {
	data, err := security.DecryptPlatform(secret)
	if err != nil {
		return nil, err
	}

	idInfo := &v1.AppleId{}

	if err = jsoniter.Unmarshal(data, &idInfo); err != nil {
		return nil, xerrors.APIPlatformAuthFailed("unmarshal apple id info failed")
	}

	now := time.Now()
	if now.After(xtime.Time(idInfo.T).Add(appleTokenExpire)) {
		return nil, xerrors.APIPlatformAuthFailed("apple id info expired")
	}

	return idInfo, nil
}
