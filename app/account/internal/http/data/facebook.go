package data

import (
	"context"

	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/lares/app/account/internal/http/domain"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (d *accountData) GetByFacebook(ctx context.Context, facebook string) (*domain.Account, error) {
	if facebook == "" {
		return nil, xerrors.APIDBFailed("facebook is empty")
	}

	po := Account{}
	result := d.data.Mdb.Debug().WithContext(ctx).Where("`facebook`=?", facebook).First(&po)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, xerrors.APINotFound("facebook=%s", facebook)
		}
		return nil, xerrors.APIDBFailed("facebook=%s", facebook).WithCause(result.Error)
	}
	return po2bo(&po), nil
}
