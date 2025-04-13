package data

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/lares/app/notice/internal/data"
	"github.com/go-pantheon/lares/app/notice/internal/http/biz"
	"gorm.io/gorm"
)

var _ biz.NoticeRepo = (*noticeData)(nil)

type Notice struct {
	Title     string
	Content   string
	StartTime time.Time
	EndTime   time.Time
}

func noticePo2bo(in *Notice) (out *biz.Notice) {
	out = &biz.Notice{
		Title:     in.Title,
		Content:   in.Content,
		StartTime: in.StartTime,
		EndTime:   in.EndTime,
	}
	return
}

type noticeData struct {
	data *data.Data
	log  *log.Helper
}

func NewNoticeData(data *data.Data, logger log.Logger) (biz.NoticeRepo, error) {
	r := &noticeData{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "account/interface/data/notice")),
	}
	return r, nil
}

func (d *noticeData) NoticeList(ctx context.Context, now time.Time) ([]*biz.Notice, error) {
	var pos []Notice

	if err := d.data.Mdb.WithContext(ctx).Limit(10).Order("sort desc").Order("end_time").Order("start_time").Where("start_time<=? AND end_time>=?", now, now).Find(&pos).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*biz.Notice{}, nil
		}
		return nil, errors.InternalServer("notice.NoticeList", fmt.Sprintf("now<%s>", now)).WithCause(err)
	}

	bos := make([]*biz.Notice, 0, len(pos))
	for _, po := range pos {
		po := po
		bos = append(bos, noticePo2bo(&po))
	}
	return bos, nil
}
