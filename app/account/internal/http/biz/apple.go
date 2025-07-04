package biz

import (
	"context"

	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/lares/app/account/internal/http/domain"
	v1 "github.com/go-pantheon/lares/gen/api/server/account/interface/account/v1"
)

func (uc *AccountUseCase) LoginByAppleToken(ctx context.Context, token, secret string, ip string) (acc *domain.Account, isRegister bool, state string, err error) {
	var (
		appleId string
	)

	if appleId, err = uc.appledo.RequestPlatformIdByAppToken(ctx, token); err != nil {
		return nil, false, "", xerrors.ErrAPIPlatformAuthFailed.WithCause(err)
	}

	if len(secret) > 0 {
		var idInfo *v1.AppleId

		if idInfo, err = uc.appledo.CheckAppleId(ctx, secret); err != nil {
			return nil, false, "", xerrors.ErrAPIPlatformAuthFailed.WithCause(err)
		}

		if idInfo.AppleId != appleId {
			err = xerrors.ErrAPIPlatformAuthFailed.WithMetadata(map[string]string{
				"secretAppleId": idInfo.AppleId,
				"tokenAppleId":  appleId,
			})

			return nil, false, "", err
		}

		state = idInfo.State
	}

	acc, isRegister, err = uc.appledo.GetOrCreateAccount(ctx, appleId, ip)
	if err != nil {
		return nil, false, "", err
	}

	return acc, isRegister, state, nil
}

func (uc *AccountUseCase) LoginByAppleCode(ctx context.Context, code string, ip string) (*domain.Account, bool, error) {
	token, err := uc.appledo.RequestToken(ctx, code)
	if err != nil {
		return nil, false, err
	}

	appleId, err := uc.appledo.RequestPlatformIdByWebToken(ctx, token)
	if err != nil {
		return nil, false, xerrors.ErrAPIPlatformAuthFailed.WithCause(err)
	}

	return uc.appledo.GetOrCreateAccount(ctx, appleId, ip)
}

func (uc *AccountUseCase) LoginByAppleId(ctx context.Context, secret string) (*domain.Account, string, error) {
	idInfo, err := uc.appledo.CheckAppleId(ctx, secret)
	if err != nil {
		return nil, "", xerrors.ErrAPIPlatformAuthFailed.WithCause(err)
	}

	acc, err := uc.appledo.GetAccount(ctx, idInfo.AppleId)
	if err != nil {
		return nil, "", xerrors.ErrAPIPlatformAuthFailed.WithCause(err)
	}

	return acc, idInfo.State, nil
}
