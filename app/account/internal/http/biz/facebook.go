package biz

import (
	"context"

	"github.com/luffy050596/rec-account/app/account/internal/http/domain"
	"github.com/luffy050596/rec-kit/profile"
	"github.com/luffy050596/rec-kit/xerrors"
)

func (uc *AccountUseCase) LoginByFacebook(ctx context.Context, token string, ip string) (*domain.Account, bool, error) {
	if profile.IsDev() {
		return nil, false, xerrors.ErrAPIStatusIllegal
	}

	facebookId, err := uc.facebookdo.RequestPlatformId(ctx, token)
	if err != nil {
		return nil, false, xerrors.ErrAPIPlatformAuthFailed.WithCause(err)
	}

	return uc.facebookdo.GetOrCreateAccount(ctx, facebookId, ip)
}
