package data

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	upg "github.com/go-pantheon/fabrica-util/data/db/postgresql"
	"github.com/go-pantheon/fabrica-util/xid"
	"github.com/go-pantheon/lares/app/account/internal/admin/biz"
	"github.com/go-pantheon/lares/app/account/internal/data"
)

const (
	minQaAccountId = 1
	maxQaAccountId = 1000
)

type Account struct {
	Id          int64
	Name        string
	Apple       string
	Google      string
	RegisterIp  string
	LastLoginIp string
	Channel     int32
	State       int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func accountPo2bo(in *Account) (out *biz.Account, err error) {
	out = &biz.Account{
		Id:          in.Id,
		Name:        in.Name,
		AppleId:     in.Apple,
		GoogleId:    in.Google,
		RegisterIp:  in.RegisterIp,
		LastLoginIp: in.LastLoginIp,
		Channel:     in.Channel,
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
	}

	if out.IdStr, err = xid.EncodeID(out.Id); err != nil {
		return nil, xerrors.APICodecFailed("encode id=%d", out.Id)
	}

	return out, nil
}

var _ biz.AccountRepo = (*accountData)(nil)

type accountData struct {
	data *data.Data
	log  *log.Helper
}

func NewAccountData(data *data.Data, logger log.Logger) biz.AccountRepo {
	return &accountData{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "account/admin/data/account")),
	}
}

func (d *accountData) GetById(ctx context.Context, id int64) (*biz.Account, error) {
	po := Account{Id: id}

	fb := upg.NewSelectSQLFieldBuilder()
	fb.Append("id", &po.Id)
	fb.Append("name", &po.Name)
	fb.Append("apple", &po.Apple)
	fb.Append("google", &po.Google)
	fb.Append("register_ip", &po.RegisterIp)
	fb.Append("last_login_ip", &po.LastLoginIp)
	fb.Append("channel", &po.Channel)
	fb.Append("state", &po.State)
	fb.Append("created_at", &po.CreatedAt)
	fb.Append("updated_at", &po.UpdatedAt)

	fieldSql, values := fb.Build()

	sqlStr := fmt.Sprintf("SELECT %s FROM accounts WHERE id = $1", fieldSql)

	row := d.data.Pdb.QueryRowContext(ctx, sqlStr, id)
	if err := row.Scan(values...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, xerrors.APINotFound("id=%d", id)
		}

		return nil, xerrors.APIDBFailed("id=%d", id).WithCause(err)
	}

	return accountPo2bo(&po)
}

func (d *accountData) GetList(ctx context.Context, index, size int64, cond *biz.Account) ([]*biz.Account, error) {
	fb := upg.NewSelectSQLFieldBuilder()
	fb.Append("id", nil)
	fb.Append("name", nil)
	fb.Append("apple", nil)
	fb.Append("google", nil)
	fb.Append("register_ip", nil)
	fb.Append("last_login_ip", nil)
	fb.Append("channel", nil)
	fb.Append("state", nil)
	fb.Append("created_at", nil)
	fb.Append("updated_at", nil)

	fieldSql, _ := fb.Build()

	sqlStr := fmt.Sprintf("SELECT %s FROM accounts LIMIT $1 OFFSET $2", fieldSql)

	rows, err := d.data.Pdb.QueryContext(ctx, sqlStr, size, index)
	if err != nil {
		return nil, xerrors.APIDBFailed("cond=%+v", cond).WithCause(err)
	}

	defer rows.Close()

	var bos []*biz.Account

	for rows.Next() {
		po := &Account{}
		if err := rows.Scan(
			&po.Id, &po.Name, &po.Apple,
			&po.Google, &po.RegisterIp, &po.LastLoginIp,
			&po.Channel, &po.State, &po.CreatedAt,
			&po.UpdatedAt); err != nil {
			return nil, xerrors.APIDBFailed("cond=%+v", cond).WithCause(err)
		}

		bo, err := accountPo2bo(po)
		if err != nil {
			return nil, err
		}

		bos = append(bos, bo)
	}

	return bos, nil
}

func (d *accountData) Count(ctx context.Context, cond *biz.Account) (count int64, err error) {
	whereSql, values := whereDB(1, cond)

	sqlStr := fmt.Sprintf("SELECT COUNT(*) FROM accounts WHERE %s", whereSql)

	row := d.data.Pdb.QueryRowContext(ctx, sqlStr, values...)
	if err := row.Scan(&count); err != nil {
		return 0, xerrors.APIDBFailed("cond=%+v", cond).WithCause(err)
	}

	return count, nil
}

func whereDB(index int64, cond *biz.Account) (string, []any) {
	b := strings.Builder{}

	b.WriteString(fmt.Sprintf("`id` < $%d OR `id` > $%d", index, index+1))

	values := []any{minQaAccountId, maxQaAccountId}

	if cond == nil {
		return b.String(), values
	}

	if len(cond.AppleId) > 0 {
		b.WriteString(fmt.Sprintf(" AND `apple` = $%d", index+2))

		values = append(values, cond.AppleId)
	}

	if len(cond.GoogleId) > 0 {
		b.WriteString(fmt.Sprintf(" AND `google` = $%d", index+3))

		values = append(values, cond.GoogleId)
	}

	if len(cond.Name) > 0 {
		b.WriteString(fmt.Sprintf(" AND `name` LIKE $%d", index+4))

		values = append(values, cond.Name)
	}

	return b.String(), values
}
