package data

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	upg "github.com/go-pantheon/fabrica-util/data/db/postgresql"
	"github.com/go-pantheon/fabrica-util/data/db/postgresql/migrate"
	"github.com/go-pantheon/lares/app/notice/internal/data"
	"github.com/go-pantheon/lares/app/notice/internal/http/biz"
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

	if err := migrate.Migrate(context.Background(), data.Pdb.GetPool(), "notices", &Notice{}, map[string]string{}); err != nil {
		return nil, errors.InternalServer("notice.NewNoticeData", "db mirgrate failed. table=notices").WithCause(err)
	}

	return r, nil
}

func (d *noticeData) NoticeList(ctx context.Context, now time.Time) ([]*biz.Notice, error) {
	fb := upg.NewSelectSQLFieldBuilder()
	fb.Append("title", nil)
	fb.Append("content", nil)
	fb.Append("start_time", nil)
	fb.Append("end_time", nil)

	fieldSql, _ := fb.Build()

	sqlStr := fmt.Sprintf("SELECT %s FROM notices WHERE start_time<=$	1 AND end_time>=$2 ORDER BY sort DESC, end_time DESC, start_time DESC", fieldSql)

	rows, err := d.data.Pdb.QueryContext(ctx, sqlStr, now, now)
	if err != nil {
		return nil, errors.InternalServer("notice.NoticeList", fmt.Sprintf("now<%s>", now)).WithCause(err)
	}

	defer rows.Close()

	var pos []*Notice

	for rows.Next() {
		po := &Notice{}
		if err := rows.Scan(&po.Title, &po.Content, &po.StartTime, &po.EndTime); err != nil {
			return nil, errors.InternalServer("notice.NoticeList", fmt.Sprintf("now<%s>", now)).WithCause(err)
		}

		pos = append(pos, po)
	}

	bos := make([]*biz.Notice, 0, len(pos))
	for _, po := range pos {
		bos = append(bos, noticePo2bo(po))
	}

	return bos, nil
}
