package v1

import (
	"context"
	"strings"

	"github.com/go-pantheon/fabrica-kit/ip"
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	v1 "github.com/go-pantheon/lares/gen/api/server/account/interface/account/v1"
)

func (s *AccountInterface) GoogleLogin(ctx context.Context, req *v1.GoogleLoginRequest) (*v1.GoogleLoginResponse, error) {
	if profile.IsDev() {
		return nil, xerrors.APIStatusIllegal("google login is disabled in dev mode")
	}

	token := strings.TrimSpace(req.Token)
	color := strings.TrimSpace(req.Color)

	if len(token) == 0 {
		return nil, xerrors.APIParamInvalid("token is empty")
	}

	account, isRegister, err := s.ac.LoginByGoogle(ctx, token, ip.GetClientIP(ctx))
	if err != nil {
		return nil, err
	}

	info, err := s.genLoginInfo(ctx, account, isRegister, "", color)
	if err != nil {
		return nil, err
	}

	return &v1.GoogleLoginResponse{
		Info: info,
	}, nil
}
