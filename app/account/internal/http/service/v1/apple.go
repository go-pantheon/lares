package v1

import (
	"context"
	"fmt"
	gohttp "net/http"
	"strings"

	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-pantheon/fabrica-kit/ip"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/lares/app/account/internal/http/domain"
	"github.com/go-pantheon/lares/app/account/internal/pkg/security"
	v1 "github.com/go-pantheon/lares/gen/api/server/account/interface/account/v1"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
)

const _callback = "lares://go-pantheon.io/auth"

func (s *AccountInterface) AppleLogin(ctx context.Context, req *v1.AppleLoginRequest) (*v1.AppleLoginResponse, error) {
	var (
		token   = strings.TrimSpace(req.Token)
		appleId = strings.TrimSpace(req.AppleId)
		color   = strings.TrimSpace(req.Color)
	)

	if len(token) > 0 {
		info, err := s.tokenLogin(ctx, token, appleId, color)
		if err != nil {
			return nil, err
		}

		return &v1.AppleLoginResponse{Info: info}, nil
	}

	if len(appleId) > 0 {
		info, err := s.appleIdLogin(ctx, appleId, color)
		if err != nil {
			return nil, err
		}

		return &v1.AppleLoginResponse{Info: info}, nil
	}

	return nil, xerrors.APIParamInvalid("token or appleId is empty")
}

func (s *AccountInterface) tokenLogin(ctx context.Context, token, appleId, color string) (*v1.LoginInfo, error) {
	account, isRegister, state, err := s.ac.LoginByAppleToken(ctx, token, appleId, ip.GetClientIP(ctx))
	if err != nil {
		return nil, err
	}

	info, err := s.genLoginInfo(ctx, account, isRegister, state, color)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (s *AccountInterface) appleIdLogin(ctx context.Context, appleId string, color string) (*v1.LoginInfo, error) {
	account, state, err := s.ac.LoginByAppleId(ctx, appleId)
	if err != nil {
		return nil, err
	}

	return s.genLoginInfo(ctx, account, false, state, color)
}

func (s *AccountInterface) AppleLoginCallback(ctx context.Context, req *v1.AppleLoginCallbackRequest) (*v1.AppleLoginCallbackResponse, error) {
	var (
		reply = &v1.AppleLoginCallbackResponse{}
		err   error
	)

	kctx, ok := ctx.(http.Context)
	if !ok {
		s.appleCallbackFail(kctx, v1.AppleLoginCallbackResponse_CODE_ERR_UNSPECIFIED, errors.New("ctx is not kratos.Context"))
		return reply, nil
	}

	if len(req.Error) > 0 {
		s.appleCallbackFail(kctx, v1.AppleLoginCallbackResponse_CODE_ERR_APPLE, errors.New(req.Error))
		return reply, nil
	}

	var (
		code  = strings.TrimSpace(req.Code)
		state = strings.TrimSpace(req.State)
	)

	var (
		account    *domain.Account
		isRegister bool
		info       *v1.LoginInfo
		k          string
	)

	if len(code) == 0 {
		s.appleCallbackFail(kctx, v1.AppleLoginCallbackResponse_CODE_ERR_PARAM, errors.New("code is empty"))
		return reply, nil
	}

	if account, isRegister, err = s.ac.LoginByAppleCode(ctx, code, ip.GetClientIP(ctx)); err != nil {
		if errors.Is(err, xerrors.ErrAPIPlatformAuthFailed) {
			s.appleCallbackFail(kctx, v1.AppleLoginCallbackResponse_CODE_ERR_TOKEN, err)
		} else {
			s.appleCallbackFail(kctx, v1.AppleLoginCallbackResponse_CODE_ERR_UNSPECIFIED, err)
		}

		return reply, nil
	}

	if info, err = s.genLoginInfo(ctx, account, isRegister, state, ""); err != nil {
		s.appleCallbackFail(kctx, v1.AppleLoginCallbackResponse_CODE_ERR_UNSPECIFIED, err)
		return reply, nil
	}

	if k, err = encryptLoginInfo(info); err != nil {
		s.appleCallbackFail(kctx, v1.AppleLoginCallbackResponse_CODE_ERR_UNSPECIFIED, err)
		return reply, nil
	}

	s.appleCallbackSuccess(kctx, k)

	return reply, nil
}

func (s *AccountInterface) appleCallbackSuccess(kctx http.Context, k string) {
	cb := fmt.Sprintf("%s?k=%s", _callback, k)
	gohttp.Redirect(kctx.Response(), kctx.Request(), cb, gohttp.StatusSeeOther)
}

func (s *AccountInterface) appleCallbackFail(kctx http.Context, errCode v1.AppleLoginCallbackResponse_Code, err error) {
	s.log.WithContext(kctx).Errorf("apple failed callback. code=%s, %+v", v1.AppleLoginCallbackResponse_Code_name[int32(errCode)], err)

	cb := fmt.Sprintf("%s?e=%s", _callback, v1.AppleLoginCallbackResponse_Code_name[int32(errCode)])
	gohttp.Redirect(kctx.Response(), kctx.Request(), cb, gohttp.StatusSeeOther)
}

func encryptLoginInfo(info *v1.LoginInfo) (string, error) {
	data, err := protojson.Marshal(info)
	if err != nil {
		return "", errors.Wrapf(err, "login info encode failed")
	}

	ser, err := security.EncryptPlatform(data)
	if err != nil {
		return "", errors.Wrapf(err, "login info encrypt failed")
	}

	return ser, nil
}
