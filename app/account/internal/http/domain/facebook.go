package domain

import (
	"context"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/lares/app/account/internal/conf"
	"github.com/huandu/facebook/v2"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

type FacebookDomain struct {
	log  *log.Helper
	repo AccountRepo

	conf *conf.Facebook
	fb   *facebook.App
	cli  *http.Client
}

func NewFacebookDomain(logger log.Logger, label *conf.Label, c *conf.Platform, repo AccountRepo) (do *FacebookDomain, err error) {
	do = &FacebookDomain{
		log:  log.NewHelper(log.With(logger, "module", "account/interface/domain/facebook")),
		conf: proto.Clone(c.Facebook).(*conf.Facebook),
		repo: repo,
	}

	if profile.IsDevStr(label.Profile) {
		do.log.Infof("use mock facebook login client")
		return
	}

	do.fb = facebook.New(do.conf.AppId, do.conf.AppSecret)
	do.fb.EnableAppsecretProof = true
	facebook.RFC3339Timestamps = true

	do.cli = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
			TLSHandshakeTimeout: 5 * time.Second,
			DisableKeepAlives:   false,
		},
	}

	return do, nil
}

func (do *FacebookDomain) RequestPlatformId(ctx context.Context, token string) (string, error) {
	if profile.IsDev() {
		return "", xerrors.APIPlatformAuthFailed("facebook token check is disabled in dev mode")
	}

	session := do.fb.Session(token)
	session.HttpClient = do.cli

	facebookId, err := session.WithContext(ctx).User()
	if err != nil {
		return "", xerrors.APIPlatformAuthFailed("check facebook token failed")
	}
	return facebookId, nil
}

func (do *FacebookDomain) GetOrCreateAccount(ctx context.Context, facebookId string, ip string) (acc *Account, isCreated bool, err error) {
	acc, err = do.repo.GetByFacebook(ctx, facebookId)
	if err == nil {
		return
	}

	if !errors.Is(err, xerrors.ErrDBRecordNotFound) {
		return
	}

	param := &Account{
		Channel:     int32(ChannelTypeFacebook),
		Facebook:    facebookId,
		RegisterIp:  ip,
		LastLoginIp: ip,
	}
	acc, err = do.repo.Create(ctx, param)
	isCreated = true
	return
}

func (do *FacebookDomain) GetAccount(ctx context.Context, facebookId string) (*Account, error) {
	return do.repo.GetByFacebook(ctx, facebookId)
}
