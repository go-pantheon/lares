package biz

import (
	"context"

	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/lares/app/account/internal/http/domain"
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
