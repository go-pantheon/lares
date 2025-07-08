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

func (d *accountData) GetByUsername(ctx context.Context, name string) (*domain.Account, error) {
	if name == "" {
		return nil, xerrors.APIDBFailed("username is empty")
	}

	po := Account{}

	fb := pg.NewSelectSQLFieldBuilder()
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

	sqlStr := fmt.Sprintf("SELECT %s FROM accounts WHERE username = $1", fieldSql)

	row := d.data.Pdb.QueryRowContext(ctx, sqlStr, name)
	if err := row.Scan(values...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, xerrors.APINotFound("name:%s", name)
		}

		return nil, xerrors.APIDBFailed("name:%s", name).WithCause(err)
	}

	return po2bo(&po), nil
}

func (d *accountData) UpdatePasswordHash(ctx context.Context, id int64, passwordHash string) (bool, error) {
	fb := pg.NewUpdateSQLFieldBuilder(2)

	fb.Append("password", passwordHash)

	fieldSql, values := fb.Build()
	values = pg.AppendValueFirst(values, id)
	sqlStr := fmt.Sprintf("UPDATE accounts SET %s WHERE account_id = $1", fieldSql)

	result, err := d.data.Pdb.ExecContext(ctx, sqlStr, values...)
	if err != nil {
		return false, xerrors.APIDBFailed("id:%d", id).WithCause(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, xerrors.APIDBFailed("id:%d", id).WithCause(err)
	}

	if rows == 0 {
		return false, xerrors.APIDBFailed("id:%d", id)
	}

	return true, nil
}
