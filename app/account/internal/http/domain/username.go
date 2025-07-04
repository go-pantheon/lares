package domain

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type UsernameDomain struct {
	log  *log.Helper
	repo AccountRepo
}

func NewUsernameDomain(logger log.Logger, repo AccountRepo) (do *UsernameDomain) {
	do = &UsernameDomain{
		log:  log.NewHelper(log.With(logger, "module", "account/http/domain/username")),
		repo: repo,
	}

	return do
}

func (do *UsernameDomain) GetAccountByUsername(ctx context.Context, name string) (*Account, error) {
	return do.repo.GetByUsername(ctx, name)
}

func (do *UsernameDomain) UpdatePasswordHash(ctx context.Context, id int64, password string) (bool, error) {
	return do.repo.UpdatePasswordHash(ctx, id, password)
}
