package data

import (
	"context"

	"github.com/luffy050596/rec-account/app/account/internal/http/domain"
	"github.com/luffy050596/rec-kit/xerrors"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (d *accountData) GetByUsername(ctx context.Context, name string) (*domain.Account, error) {
	if name == "" {
		return nil, xerrors.APIDBFailed("username is empty")
	}

	po := Account{}
	result := d.data.Mdb.Debug().WithContext(ctx).Where("`username`=?", name).First(&po)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, xerrors.APINotFound("name:%s", name)
		}
		return nil, xerrors.APIDBFailed("name:%s", name).WithCause(result.Error)
	}
	return po2bo(&po), nil
}

func (d *accountData) UpdatePasswordHash(ctx context.Context, id int64, passwordHash string) (bool, error) {
	po := Account{}
	result := d.data.Mdb.Debug().WithContext(ctx).Where("`id`=?", id).First(&po)
	if result.Error != nil {
		return false, xerrors.APIDBFailed("id:%d", id).WithCause(result.Error)
	}
	po.PasswordHash = passwordHash
	result = d.data.Mdb.Debug().WithContext(ctx).Where("`id`=?", id).Updates(&po)
	if result.Error != nil {
		return false, xerrors.APIDBFailed("id:%d", id).WithCause(result.Error)
	}
	return true, nil
}
