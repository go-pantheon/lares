package biz

import (
	"context"

	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/lares/app/account/internal/http/domain"
	"github.com/go-pantheon/lares/app/account/internal/pkg/security"
	"github.com/pkg/errors"
)

func (uc *AccountUseCase) CreateByUsername(ctx context.Context, username, pwd, ip string) (*domain.Account, error) {
	acc, err := uc.usernamedo.GetAccountByUsername(ctx, username)
	if err != nil && !errors.Is(err, xerrors.ErrDBRecordNotFound) {
		return nil, err
	}

	if acc != nil {
		// return 401 to blur the error for security
		return nil, xerrors.APIAuthFailed("username=%s is already registered", username)
	}

	hashPwd, err := security.HashPassword(pwd)
	if err != nil {
		return nil, err
	}

	return uc.accountdo.Create(ctx, &domain.Account{
		Username:     username,
		PasswordHash: hashPwd,
		RegisterIp:   ip,
		LastLoginIp:  ip,
	})
}

func (uc *AccountUseCase) GetByUsernameAndPassword(ctx context.Context, name, pwd string) (*domain.Account, error) {
	po, err := uc.usernamedo.GetAccountByUsername(ctx, name)
	if err != nil {
		if errors.Is(err, xerrors.ErrDBRecordNotFound) {
			// return 401 to blur the error for security
			return nil, xerrors.APIAuthFailed("username=%s is not registered", name)
		}

		return nil, err
	}

	valid, err := security.VerifyPassword(pwd, po.PasswordHash)
	if err != nil {
		return nil, xerrors.APIAuthFailed("password is invalid").WithCause(err)
	}

	if !valid {
		return nil, xerrors.APIAuthFailed("password is invalid")
	}

	return po, nil
}
