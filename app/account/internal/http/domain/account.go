package domain

import (
	"context"
	"time"
)

type ChannelType int32

const (
	ChannelTypeDefault ChannelType = iota
	ChannelTypeApple
	ChannelTypeGoogle
	ChannelTypeFacebook
)

type Account struct {
	Id           int64
	Zone         uint8
	IdStr        string
	Apple        string
	Google       string
	Facebook     string
	Username     string
	PasswordHash string
	Salt         string
	RegisterIp   string
	LastLoginIp  string
	DefaultColor string
	Channel      int32
	State        int32
	Device       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type AccountRepo interface {
	Create(ctx context.Context, account *Account) (*Account, error)
	GetById(ctx context.Context, id int64) (*Account, error)
	GetByApple(ctx context.Context, apple string) (*Account, error)
	GetByGoogle(ctx context.Context, google string) (*Account, error)
	GetByFacebook(ctx context.Context, facebook string) (*Account, error)
	GetByUsername(ctx context.Context, name string) (*Account, error)
	UpdatePasswordHash(ctx context.Context, id int64, passwordHash string) (bool, error)
}

type AccountDomain struct {
	log  *log.Helper
	repo AccountRepo
}

func NewAccountDomain(logger log.Logger, repo AccountRepo) (do *AccountDomain) {
	do = &AccountDomain{
		log:  log.NewHelper(log.With(logger, "module", "account/http/domain/account")),
		repo: repo,
	}
	return do
}

func (do *AccountDomain) GetById(ctx context.Context, id int64) (*Account, error) {
	return do.repo.GetById(ctx, id)
}

func (do *AccountDomain) Create(ctx context.Context, acc *Account) (*Account, error) {
	return do.repo.Create(ctx, acc)
}
