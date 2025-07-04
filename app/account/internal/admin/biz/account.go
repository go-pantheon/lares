package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type Account struct {
	Id          int64
	IdStr       string
	AppleId     string
	GoogleId    string
	Name        string
	RegisterIp  string
	LastLoginIp string
	Channel     int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type AccountRepo interface {
	GetById(ctx context.Context, id int64) (*Account, error)
	GetList(ctx context.Context, index, size int64, cond *Account) ([]*Account, error)
	Count(ctx context.Context, cond *Account) (int64, error)
}

type AccountUseCase struct {
	log *log.Helper

	repo AccountRepo
}

func NewAccountUseCase(accountRepo AccountRepo, logger log.Logger) *AccountUseCase {
	return &AccountUseCase{
		log:  log.NewHelper(log.With(logger, "module", "account/admin/biz/account")),
		repo: accountRepo,
	}
}

func (uc *AccountUseCase) GetById(ctx context.Context, id int64) (*Account, error) {
	return uc.repo.GetById(ctx, id)
}

func (uc *AccountUseCase) GetList(ctx context.Context, index, size int64, cond *Account) (accounts []*Account, count int64, err error) {
	accounts, err = uc.repo.GetList(ctx, index, size, cond)
	if err != nil {
		return nil, 0, err
	}

	count, err = uc.repo.Count(ctx, cond)
	if err != nil {
		return nil, 0, err
	}

	return accounts, count, err
}
