package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/fabrica-util/data/db/pg"
	"github.com/go-pantheon/lares/app/account/internal/http/domain"
	"github.com/pkg/errors"
)

func (d *accountData) GetByFacebook(ctx context.Context, facebook string) (*domain.Account, error) {
	if facebook == "" {
		return nil, xerrors.APIDBFailed("facebook is empty")
	}

	po := Account{}

	fb := pg.NewSelectSQLFieldBuilder()
	fb.Append("facebook", &po.Facebook)
	fb.Append("google", &po.Google)
	fb.Append("apple", &po.Apple)
	fb.Append("username", &po.Username)
	fb.Append("password", &po.PasswordHash)
	fb.Append("register_ip", &po.RegisterIp)
	fb.Append("last_login_ip", &po.LastLoginIp)
	fb.Append("default_color", &po.DefaultColor)
	fb.Append("channel", &po.Channel)
	fb.Append("state", &po.State)
	fb.Append("created_at", &po.CreatedAt)
	fb.Append("updated_at", &po.UpdatedAt)

	fieldSql, values := fb.Build()

	sqlStr := fmt.Sprintf("SELECT %s FROM accounts WHERE facebook = $1", fieldSql)

	row := d.data.Pdb.QueryRowContext(ctx, sqlStr, facebook)
	if err := row.Scan(values...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, xerrors.APINotFound("facebook=%s", facebook)
		}

		return nil, xerrors.APIDBFailed("facebook=%s", facebook).WithCause(err)
	}

	return po2bo(&po), nil
}
