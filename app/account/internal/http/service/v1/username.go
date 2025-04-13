package v1

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/go-pantheon/fabrica-kit/ip"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	v1 "github.com/go-pantheon/lares/gen/api/server/account/interface/account/v1"
)

const (
	minUsernameLength = 6
	maxUsernameLength = 20
	minPasswordLength = 8
	maxPasswordLength = 32
)

func (s *AccountInterface) UsernameRegister(ctx context.Context, req *v1.UsernameRegisterRequest) (*v1.UsernameRegisterResponse, error) {
	name := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)
	color := strings.TrimSpace(req.Color)

	if l := utf8.RuneCountInString(name); l < minUsernameLength || l > maxUsernameLength {
		return nil, xerrors.APIParamInvalid("length of username=%s is not valid", name)
	}

	if l := utf8.RuneCountInString(password); l < minPasswordLength || l > maxPasswordLength {
		return nil, xerrors.APIParamInvalid("length of password=%s is not valid", password)
	}

	account, err := s.ac.CreateByUsername(ctx, name, password, ip.GetClientIP(ctx))
	if err != nil {
		return nil, err
	}

	info, err := s.genLoginInfo(ctx, account, false, "", color)
	if err != nil {
		return nil, err
	}

	return &v1.UsernameRegisterResponse{
		Info: info,
	}, nil
}

func (s *AccountInterface) UsernameLogin(ctx context.Context, req *v1.UsernameLoginRequest) (*v1.UsernameLoginResponse, error) {
	name := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)
	color := strings.TrimSpace(req.Color)

	account, err := s.ac.GetByUsernameAndPassword(ctx, name, password)
	if err != nil {
		return nil, err
	}

	info, err := s.genLoginInfo(ctx, account, false, "", color)
	if err != nil {
		return nil, err
	}

	return &v1.UsernameLoginResponse{
		Info: info,
	}, nil
}
