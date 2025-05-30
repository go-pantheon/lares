package data

import (
	"context"

	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/lares/app/account/internal/http/domain"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (d *accountData) GetByGoogle(ctx context.Context, google string) (*domain.Account, error) {
	if google == "" {
		return nil, xerrors.APIDBFailed("google is empty")
	}

	po := Account{}
	result := d.data.Mdb.Debug().WithContext(ctx).Where("`google`=?", google).First(&po)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, xerrors.APINotFound("google=%s", google)
		}
		return nil, xerrors.APIDBFailed("google=%s", google).WithCause(result.Error)
	}
	return po2bo(&po), nil
}
