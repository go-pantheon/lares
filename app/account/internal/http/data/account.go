package data

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	xid "github.com/go-pantheon/fabrica-util/id"
	"github.com/go-pantheon/lares/app/account/internal/data"
	"github.com/go-pantheon/lares/app/account/internal/http/domain"
	"github.com/go-pantheon/lares/app/account/internal/pkg/security"
	"gorm.io/gorm"
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
	out = &domain.Account{
		Id:           xid.CombineZoneId(in.AccountId, uint8(profile.Zone())),
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
	return
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
	if err := data.Mdb.AutoMigrate(&Account{}); err != nil {
		return nil, err
	}

	if err := d.createQaAccounts(context.Background(), data, 100); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *accountData) createQaAccounts(ctx context.Context, data *data.Data, num int) error {
	if !data.Mdb.Migrator().HasTable(&Account{}) {

		return xerrors.APIDBFailed("table=accounts not exists")
	}

	var count int64
	if err := d.data.Mdb.WithContext(ctx).Model(&Account{}).Count(&count).Error; err != nil {
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
	po, fields, err := d.buildCreatePo(acc)
	if err != nil {
		return nil, err
	}

	result := d.data.Mdb.WithContext(ctx).Select(fields).Create(&po)
	if result.Error != nil {
		return nil, xerrors.APIDBFailed("username=%s apple=%s google=%s facebook=%s", acc.Username, acc.Apple, acc.Google, acc.Facebook).WithCause(result.Error)
	}
	return po2bo(&po), nil
}

func (d *accountData) buildCreatePo(acc *domain.Account) (po Account, fields []string, err error) {
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

	fields = make([]string, 0, 5)
	if acc.Username != "" {
		fields = append(fields, "Username")
		fields = append(fields, "PasswordHash")
	}
	if acc.Device != "" {
		fields = append(fields, "Device")
	}
	if acc.Apple != "" {
		fields = append(fields, "Apple")
	}
	if acc.Google != "" {
		fields = append(fields, "Google")
	}
	if acc.Facebook != "" {
		fields = append(fields, "Facebook")
	}

	if len(fields) == 0 {
		err = xerrors.APIDBFailed("no username or token")
		return
	}
	fields = append(fields, "RegisterIp", "LastLoginIp", "Channel")

	return po, fields, nil
}

func (d *accountData) GetById(ctx context.Context, id int64) (*domain.Account, error) {
	zoneId, zone := xid.SplitId(id)
	po := Account{AccountId: zoneId, Zone: zone}
	result := d.data.Mdb.Debug().WithContext(ctx).First(&po)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, xerrors.APINotFound("id=%d", id)
		}
		return nil, xerrors.APIDBFailed("id=%d", id).WithCause(result.Error)
	}
	return po2bo(&po), nil
}
