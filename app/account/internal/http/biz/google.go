package biz

import (
	"context"

	"github.com/luffy050596/rec-account/app/account/internal/http/domain"
	"github.com/luffy050596/rec-kit/profile"
	"github.com/luffy050596/rec-kit/xerrors"
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
