package domain

import (
	"context"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/lares/app/account/internal/conf"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
	goauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/proto"
)

type GoogleDomain struct {
	log  *log.Helper
	repo AccountRepo

	conf      *conf.Google
	validator *idtoken.Validator
}

func NewGoogleDomain(logger log.Logger, label *conf.Label, c *conf.Platform, repo AccountRepo) (do *GoogleDomain, err error) {
	do = &GoogleDomain{
		log:  log.NewHelper(log.With(logger, "module", "account/interface/domain/google")),
		conf: proto.Clone(c.Google).(*conf.Google),
		repo: repo,
	}

	if profile.IsDevStr(label.Profile) {
		do.log.Infof("use mock google login client")
		return
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
			TLSHandshakeTimeout: 5 * time.Second,
			DisableKeepAlives:   false,
		},
	})
	jwtConf, err := google.JWTConfigFromJSON([]byte(do.conf.Json), goauth2.OpenIDScope)
	if err != nil {
		return nil, xerrors.APIPlatformAuthFailed("build google jwt config failed")
	}

	val := jwtConf.Client(ctx).Transport.(*oauth2.Transport)
	_, err = val.Source.Token()
	if err != nil {
		return nil, xerrors.APIPlatformAuthFailed("get google token failed")
	}

	do.validator, err = idtoken.NewValidator(ctx, option.WithHTTPClient(jwtConf.Client(ctx)))
	if err != nil {
		return nil, xerrors.APIPlatformAuthFailed("build google validator failed")
	}
	return do, nil
}

func (do *GoogleDomain) RequestPlatformId(ctx context.Context, token string) (string, error) {
	if profile.IsDev() {
		return "", xerrors.APIPlatformAuthFailed("google token check is disabled in dev mode")
	}

	var info *idtoken.Payload

	info, err := do.validator.Validate(ctx, token, do.conf.Aud)
	if err != nil {
		return "", xerrors.APIPlatformAuthFailed("validate google token failed")
	}

	if info.Audience != do.conf.Aud {
		return "", xerrors.APIPlatformAuthFailed("aud error")
	}

	if info.Issuer != do.conf.Iss1 && info.Issuer != do.conf.Iss2 {
		return "", xerrors.APIPlatformAuthFailed("iss error")
	}

	if time.Now().Unix() >= info.Expires {
		return "", xerrors.APIPlatformAuthFailed("token expired")
	}

	return info.Subject, nil
}

func (do *GoogleDomain) GetOrCreateAccount(ctx context.Context, googleId string, ip string) (acc *Account, isRegister bool, err error) {
	acc, err = do.repo.GetByGoogle(ctx, googleId)
	if err == nil {
		return
	}

	if !errors.Is(err, xerrors.ErrDBRecordNotFound) {
		return
	}

	param := &Account{
		Channel:     int32(ChannelTypeGoogle),
		Google:      googleId,
		RegisterIp:  ip,
		LastLoginIp: ip,
	}
	acc, err = do.repo.Create(ctx, param)
	isRegister = true
	return
}

func (do *GoogleDomain) GetAccount(ctx context.Context, googleId string) (*Account, error) {
	return do.repo.GetByGoogle(ctx, googleId)
}
