package biz

import (
	"context"

	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/lares/app/account/internal/http/domain"
)

func (uc *AccountUseCase) LoginByGoogle(ctx context.Context, token string, ip string) (*domain.Account, bool, error) {
	if profile.IsDev() {
		return nil, false, xerrors.ErrAPIStatusIllegal
	}

	googleId, err := uc.googledo.RequestPlatformId(ctx, token)
	if err != nil {
		return nil, false, err
	}

	return uc.googledo.GetOrCreateAccount(ctx, googleId, ip)
}
