package v1

import (
	"context"
	"strings"

	v1 "github.com/luffy050596/rec-account/gen/api/server/account/interface/account/v1"
	"github.com/luffy050596/rec-kit/ip"
	"github.com/luffy050596/rec-kit/profile"
	"github.com/luffy050596/rec-kit/xerrors"
)

func (s *AccountInterface) FacebookLogin(ctx context.Context, req *v1.FacebookLoginRequest) (*v1.FacebookLoginResponse, error) {
	if profile.IsDev() {
		return nil, xerrors.APIStatusIllegal("facebook login is disabled in dev mode")
	}

	token := strings.TrimSpace(req.Token)
	color := strings.TrimSpace(req.Color)

	if len(token) == 0 {
		return nil, xerrors.APIParamInvalid("token is empty")
	}

	account, isRegister, err := s.ac.LoginByFacebook(ctx, token, ip.GetClientIP(ctx))
	if err != nil {
		return nil, err
	}

	info, err := s.genLoginInfo(ctx, account, isRegister, "", color)
	if err != nil {
		return nil, err
	}
	return &v1.FacebookLoginResponse{
		Info: info,
	}, nil
}
