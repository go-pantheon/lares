package data

import (
	"context"
	"time"

	"github.com/luffy050596/rec-account/app/notice/internal/admin/biz"
	"github.com/luffy050596/rec-account/app/notice/internal/data"
	"github.com/luffy050596/rec-kit/xerrors"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ biz.NoticeRepo = (*noticeData)(nil)

type Notice struct {
	Id        int64     `gorm:"primarykey;autoIncrement"`
	Title     string    `gorm:"not null;default:''"`
	Content   string    `gorm:"not null;default:'';type:text"`
	Sort      int64     `gorm:"not null;default:0"`
	StartTime time.Time `gorm:"not null"`
	EndTime   time.Time `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type noticeData struct {
	data *data.Data
	log  *log.Helper
}

func NewNoticeData(data *data.Data, logger log.Logger) (biz.NoticeRepo, error) {
	r := &noticeData{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "notice/admin/data/notice")),
	}
	if err := data.Mdb.AutoMigrate(&Notice{}); err != nil {
		return nil, errors.Wrapf(err, "db mirgrate failed. table=notices")
	}
	return r, nil
}

func (d *noticeData) GetList(ctx context.Context, index, size int64) ([]*biz.Notice, error) {
	var pos []Notice

	if err := d.data.Mdb.WithContext(ctx).Offset(int(index)).Limit(int(size)).Find(&pos).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*biz.Notice{}, nil
		}
		return nil, err
	}

	bos := make([]*biz.Notice, 0, len(pos))
	for _, po := range pos {
		bos = append(bos, noticePo2bo(&po))
	}
	return bos, nil
}

func (d *noticeData) GetById(ctx context.Context, id int64) (*biz.Notice, error) {
	po := Notice{Id: id}
	result := d.data.Mdb.Debug().WithContext(ctx).First(&po)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, xerrors.APINotFound("id=%d", id)
		}
		return nil, xerrors.APIDBFailed("id=%d", id).WithCause(result.Error)
	}
	return noticePo2bo(&po), nil
}

func (d *noticeData) UpdateById(ctx context.Context, bo *biz.Notice) error {
	po := noticeBo2Po(bo)

	result := d.data.Mdb.WithContext(ctx).Updates(&po)
	if result.Error != nil {
		return xerrors.APIDBFailed("id=%d", bo.Id).WithCause(result.Error)
	}
	if result.RowsAffected != 1 {
		return xerrors.APINotFound("id=%d", bo.Id)
	}
	return nil
}

func (d *noticeData) Insert(ctx context.Context, bo *biz.Notice) error {
	po := noticeBo2Po(bo)

	result := d.data.Mdb.WithContext(ctx).Create(&po)
	if result.Error != nil {
		return xerrors.APIDBFailed("id=%d", bo.Id).WithCause(result.Error)
	}
	if result.RowsAffected != 1 {
		return xerrors.APIDBFailed("not affected. id=%d", bo.Id)
	}
	return nil
}

func (d *noticeData) DeleteById(ctx context.Context, id int64) error {
	result := d.data.Mdb.WithContext(ctx).Where("id", id).Delete(&Notice{})
	if result.Error != nil {
		return xerrors.APIDBFailed("id=%d", id).WithCause(result.Error)
	}
	if result.RowsAffected != 1 {
		return xerrors.APINotFound("id=%d", id)
	}
	return nil
}

func (d *noticeData) Count(ctx context.Context) (count int64, err error) {
	if err = d.data.Mdb.WithContext(ctx).Model(&Notice{}).Count(&count).Error; err != nil {
		err = xerrors.APIDBFailed("count failed").WithCause(err)
		return
	}
	return
}

func noticePo2bo(in *Notice) (out *biz.Notice) {
	out = &biz.Notice{
		Id:        in.Id,
		Title:     in.Title,
		Content:   in.Content,
		Sort:      in.Sort,
		StartTime: in.StartTime,
		EndTime:   in.EndTime,

		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}

	return
}

func noticeBo2Po(in *biz.Notice) (out *Notice) {
	out = &Notice{
		Id:        in.Id,
		Title:     in.Title,
		Content:   in.Content,
		Sort:      in.Sort,
		StartTime: in.StartTime,
		EndTime:   in.EndTime,

		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}

	return
}
