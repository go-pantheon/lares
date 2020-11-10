package biz

import (
	"context"

	"github.com/luffy050596/rec-account/app/account/internal/http/domain"
)

type AccountUseCase struct {
	log *log.Helper

	accountdo  *domain.AccountDomain
	usernamedo *domain.UsernameDomain
	googledo   *domain.GoogleDomain
	appledo    *domain.AppleDomain
	facebookdo *domain.FacebookDomain
}

func NewAccountUseCase(logger log.Logger, acc *domain.AccountDomain,
	name *domain.UsernameDomain,
	gg *domain.GoogleDomain,
	ap *domain.AppleDomain,
	fb *domain.FacebookDomain,
) (uc *AccountUseCase, err error) {
	uc = &AccountUseCase{
		log:        log.NewHelper(log.With(logger, "module", "account/http/biz/account")),
		accountdo:  acc,
		usernamedo: name,
		googledo:   gg,
		appledo:    ap,
		facebookdo: fb,
	}
	return
}

func (uc *AccountUseCase) GetById(ctx context.Context, id int64) (*domain.Account, error) {
	return uc.accountdo.GetById(ctx, id)
}
