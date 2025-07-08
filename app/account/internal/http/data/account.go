package data

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/fabrica-util/data/db/pg"
	"github.com/go-pantheon/fabrica-util/data/db/pg/migrate"
	"github.com/go-pantheon/fabrica-util/xid"
	"github.com/go-pantheon/lares/app/account/internal/data"
	"github.com/go-pantheon/lares/app/account/internal/http/domain"
	"github.com/go-pantheon/lares/app/account/internal/pkg/security"
)

const (
	defaultQAColor  = "qa"
	defaultQAPrefix = "qa"
	minQaAccountId  = 1
	maxQaAccountId  = 1000
	idZoneBit       = 8
)

type Account struct {
	AccountId    int64  `gorm:"primarykey;autoIncrement"`
	Zone         uint8  `gorm:"default:0"`
	Apple        string `gorm:"uniqueIndex"`
	Google       string `gorm:"uniqueIndex"`
	Facebook     string `gorm:"uniqueIndex"`
	Username     string `gorm:"uniqueIndex"`
	PasswordHash string `gorm:"not null;default:'';column:password"`
	RegisterIp   string `gorm:"default:''"`
	LastLoginIp  string `gorm:"default:''"`
	DefaultColor string `gorm:"default:''"`
	Channel      int32  `gorm:"default:0"`
	State        int32  `gorm:"default:0"`
	Device       string `gorm:"default:''"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func po2bo(in *Account) (out *domain.Account) {
	uid, err := xid.BuildUID(in.AccountId, uint8(profile.Zone()))
	if err != nil {
		return nil
	}

	out = &domain.Account{
		Id:           uid,
		Zone:         in.Zone,
		Apple:        in.Apple,
		Google:       in.Google,
		Facebook:     in.Facebook,
		Username:     in.Username,
		PasswordHash: in.PasswordHash,
		RegisterIp:   in.RegisterIp,
		LastLoginIp:  in.LastLoginIp,
		DefaultColor: in.DefaultColor,
		Channel:      in.Channel,
		Device:       in.Device,
		State:        in.State,
		CreatedAt:    in.CreatedAt,
		UpdatedAt:    in.UpdatedAt,
	}

	return out
}

var _ domain.AccountRepo = (*accountData)(nil)

type accountData struct {
	once *sync.Once

	data *data.Data
	log  *log.Helper
}

func NewAccountData(data *data.Data, logger log.Logger) (domain.AccountRepo, error) {
	d := &accountData{
		once: &sync.Once{},
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "account/http/data/account")),
	}

	if err := migrate.Migrate(context.Background(), data.Pdb, "accounts", &Account{}, nil); err != nil {
		return nil, err
	}

	if err := d.createQaAccounts(context.Background(), 100); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *accountData) createQaAccounts(ctx context.Context, num int) error {
	var count int64
	if err := d.data.Pdb.QueryRowContext(ctx, "SELECT COUNT(*) FROM accounts").Scan(&count); err != nil {
		return xerrors.APIDBFailed("count failed").WithCause(err)
	}

	if count > 0 {
		return nil
	}

	d.once.Do(func() {
		d.log.WithContext(ctx).Infof("create qa accounts. id=1~%d", num)

		for i := minQaAccountId; i <= maxQaAccountId; i++ {
			username := fmt.Sprintf("%s%d", defaultQAPrefix, i)

			password, err := security.HashPassword(username)
			if err != nil {
				d.log.WithContext(ctx).Errorf("hash password failed. err=%v", err)
				return
			}

			if _, err := d.Create(ctx, &domain.Account{
				Id:           int64(i),
				Username:     username,
				PasswordHash: password,
				DefaultColor: defaultQAColor,
			}); err != nil {
				d.log.WithContext(ctx).Errorf("create qa account failed. err=%v", err)
				return
			}
		}
	})

	return nil
}

func (d *accountData) Create(ctx context.Context, acc *domain.Account) (*domain.Account, error) {
	po := d.buildCreatePo(acc)

	fb := pg.NewInsertSQLFieldBuilder()

	fb.Append("zone", po.Zone)
	fb.Append("device", po.Device)
	fb.Append("apple", po.Apple)
	fb.Append("google", po.Google)
	fb.Append("facebook", po.Facebook)
	fb.Append("username", po.Username)
	fb.Append("password", po.PasswordHash)
	fb.Append("register_ip", po.RegisterIp)
	fb.Append("last_login_ip", po.LastLoginIp)
	fb.Append("default_color", po.DefaultColor)
	fb.Append("channel", po.Channel)

	colSql, argSql, values := fb.Build()

	sqlStr := fmt.Sprintf("INSERT INTO accounts (%s) VALUES (%s)", colSql, argSql)

	result, err := d.data.Pdb.ExecContext(ctx, sqlStr, values...)
	if err != nil {
		return nil, xerrors.APIDBFailed("%s", acc.LogInfo()).WithCause(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, xerrors.APIDBFailed("%s", acc.LogInfo()).WithCause(err)
	}

	if rows == 0 {
		return nil, xerrors.APIDBFailed("%s", acc.LogInfo())
	}

	return po2bo(&po), nil
}

func (d *accountData) GetById(ctx context.Context, id int64) (*domain.Account, error) {
	gameID, zone := xid.SplitUID(id)
	po := Account{AccountId: gameID, Zone: zone}

	fb := pg.NewSelectSQLFieldBuilder()
	fb.Append("account_id", &po.AccountId)
	fb.Append("zone", &po.Zone)
	fb.Append("device", &po.Device)
	fb.Append("apple", &po.Apple)
	fb.Append("google", &po.Google)
	fb.Append("facebook", &po.Facebook)
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

	sqlStr := fmt.Sprintf("SELECT %s FROM accounts WHERE account_id = $1", fieldSql)

	row := d.data.Pdb.QueryRowContext(ctx, sqlStr, id)
	if err := row.Scan(values...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, xerrors.APINotFound("id=%d", id)
		}

		return nil, xerrors.APIDBFailed("id=%d", id).WithCause(err)
	}

	return po2bo(&po), nil
}

func (d *accountData) buildCreatePo(acc *domain.Account) (po Account) {
	po = Account{
		Zone:         uint8(profile.Zone()),
		Device:       acc.Device,
		Apple:        acc.Apple,
		Google:       acc.Google,
		Facebook:     acc.Facebook,
		Username:     acc.Username,
		PasswordHash: acc.PasswordHash,
		RegisterIp:   acc.RegisterIp,
		LastLoginIp:  acc.LastLoginIp,
		DefaultColor: acc.DefaultColor,
		Channel:      acc.Channel,
	}

	return po
}
