package v1

import (
	"context"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/fabrica-util/id"
	"github.com/go-pantheon/lares/app/account/internal/http/biz"
	v1 "github.com/go-pantheon/lares/gen/api/server/account/interface/account/v1"
)

const (
	RandLength             = 8
	TokenExpiredDuration   = time.Minute * 10
	SessionExpiredDuration = time.Hour * 24 * 15
)

type AccountInterface struct {
	v1.UnimplementedAccountInterfaceServer

	log *log.Helper
	ac  *biz.AccountUseCase
}

func NewAccountInterface(logger log.Logger, ac *biz.AccountUseCase) *AccountInterface {
	return &AccountInterface{
		log: log.NewHelper(log.With(logger, "module", "account/interface/account/v1")),
		ac:  ac,
	}
}

func (s *AccountInterface) ClientNetHealthy(ctx context.Context, req *v1.DevPingRequest) (*v1.DevPingResponse, error) {

	return &v1.DevPingResponse{
		Message: req.Message,
		Time:    time.Now().Format(time.RFC3339),
	}, nil
}

func (s *AccountInterface) Refresh(ctx context.Context, req *v1.RefreshRequest) (*v1.RefreshResponse, error) {

	accountId, err := id.DecodeId(req.AccountId)
	if err != nil {
		return nil, xerrors.APIParamInvalid("id=%s", req.AccountId).WithCause(err)
	}

	if len(req.Session) == 0 {
		return nil, xerrors.APIParamInvalid("session is empty")
	}
	ss, err := unmarshalSession(req.Session)
	if err != nil {
		return nil, xerrors.APIParamInvalid("session=%s", req.Session).WithCause(err)
	}
	if time.Now().Unix() > ss.Timeout {
		return nil, xerrors.APISessionTimeout("session timeout")
	}

	if ss.AccountId != accountId {
		return nil, xerrors.APISessionIllegal("ssAccountId=%d, accountId=%d", ss.AccountId, accountId)
	}

	account, err := s.ac.GetById(ctx, accountId)
	if err != nil {
		return nil, err
	}

	session, sessionTimeout, err := genSession(account)
	if err != nil {
		return nil, err
	}

	return &v1.RefreshResponse{
		Session:        session,
		SessionTimeout: sessionTimeout,
	}, nil
}

func (s *AccountInterface) Token(ctx context.Context, req *v1.TokenRequest) (*v1.TokenResponse, error) {
	accountId, err := id.DecodeId(req.AccountId)
	if err != nil {
		return nil, xerrors.APIParamInvalid("id=%s", req.AccountId).WithCause(err)
	}

	if len(req.Session) == 0 {
		return nil, xerrors.APIParamInvalid("session is empty")
	}
	ss, err := unmarshalSession(req.Session)
	if err != nil {
		return nil, xerrors.APIParamInvalid("session=%s", req.Session).WithCause(err)
	}
	if time.Now().Unix() > ss.Timeout {
		return nil, xerrors.APISessionTimeout("session timeout")
	}

	if ss.AccountId != accountId {
		return nil, xerrors.APISessionIllegal("session=%s, accountId=%s", req.Session, req.AccountId)
	}

	color := strings.TrimSpace(req.Color)

	account, err := s.ac.GetById(ctx, accountId)
	if err != nil {
		return nil, err
	}

	token, tokenTimeout, err := genToken(account, color)
	if err != nil {
		return nil, err
	}

	session, sessionTimeout, err := genSession(account)
	if err != nil {
		return nil, err
	}

	return &v1.TokenResponse{
		Token:          token,
		Session:        session,
		TokenTimeout:   tokenTimeout,
		SessionTimeout: sessionTimeout,
	}, nil
}
