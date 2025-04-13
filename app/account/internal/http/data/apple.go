package data

import (
	"context"

	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/lares/app/account/internal/http/domain"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (d *accountData) GetByApple(ctx context.Context, apple string) (*domain.Account, error) {
	if apple == "" {
		return nil, xerrors.APIDBFailed("apple is empty")
	}

	po := Account{}
	result := d.data.Mdb.Debug().WithContext(ctx).Where("`apple`=?", apple).First(&po)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, xerrors.APINotFound("apple=%s", apple)
		}
		return nil, xerrors.APIDBFailed("apple=%s", apple).WithCause(result.Error)
	}
	return po2bo(&po), nil
}
