package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/fabrica-util/id"
	"github.com/go-pantheon/lares/app/account/internal/admin/biz"
	"github.com/go-pantheon/lares/app/account/internal/data"
	"gorm.io/gorm"
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

	if out.IdStr, err = id.EncodeId(out.Id); err != nil {
		err = xerrors.APICodecFailed("encode id=%d", out.Id)
		return
	}
	return
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
	result := d.data.Mdb.Debug().WithContext(ctx).First(&po)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, xerrors.APINotFound("id=%d", id)
		}
		return nil, xerrors.APIDBFailed("id=%d", id).WithCause(result.Error)
	}
	return accountPo2bo(&po)
}

func (d *accountData) GetList(ctx context.Context, index, size int64, cond *biz.Account) ([]*biz.Account, error) {
	var pos []Account

	if err := whereDB(d.data.Mdb.WithContext(ctx), cond).Offset(int(index)).Limit(int(size)).Find(&pos).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*biz.Account{}, nil
		}
		return nil, xerrors.APIDBFailed("cond=%+v", cond).WithCause(err)
	}

	bos := make([]*biz.Account, 0, len(pos))
	for _, po := range pos {
		bo, err := accountPo2bo(&po)
		if err != nil {
			return nil, err
		}
		bos = append(bos, bo)
	}
	return bos, nil
}

func (d *accountData) Count(ctx context.Context, cond *biz.Account) (count int64, err error) {
	if err = whereDB(d.data.Mdb.WithContext(ctx), cond).Model(&Account{}).Count(&count).Error; err != nil {
		err = xerrors.APIDBFailed("db query count failed").WithCause(err)
		return
	}
	return
}

func whereDB(db *gorm.DB, cond *biz.Account) *gorm.DB {
	db = db.Where("`id` < ? OR `id` > ?", minQaAccountId, maxQaAccountId)

	if cond == nil {
		return db
	}

	if len(cond.AppleId) > 0 {
		db = db.Where("`apple` = ?", cond.AppleId)
	}
	if len(cond.GoogleId) > 0 {
		db = db.Where("`google` = ?", cond.GoogleId)
	}
	if len(cond.Name) > 0 {
		db = db.Where("`name` LIKE ?", "%"+cond.Name+"%")
	}
	return db
}
