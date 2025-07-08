package data

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/fabrica-util/data/db/pg"
	"github.com/go-pantheon/lares/app/notice/internal/admin/biz"
	"github.com/go-pantheon/lares/app/notice/internal/data"
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

	return r, nil
}

func (d *noticeData) GetList(ctx context.Context, index, size int64) ([]*biz.Notice, error) {
	fb := pg.NewSelectSQLFieldBuilder()
	fb.Append("title", nil)
	fb.Append("content", nil)
	fb.Append("start_time", nil)
	fb.Append("end_time", nil)

	fieldSql, _ := fb.Build()

	sqlStr := fmt.Sprintf("SELECT %s FROM notices LIMIT $1 OFFSET $2 ORDER BY sort DESC, end_time DESC, start_time DESC", fieldSql)

	rows, err := d.data.Pdb.QueryContext(ctx, sqlStr, size, index)
	if err != nil {
		return nil, xerrors.APIDBFailed("index=%d, size=%d", index, size).WithCause(err)
	}

	defer rows.Close()

	var pos []*Notice

	for rows.Next() {
		po := &Notice{}
		if err := rows.Scan(&po.Title, &po.Content, &po.StartTime, &po.EndTime); err != nil {
			return nil, xerrors.APIDBFailed("index=%d, size=%d", index, size).WithCause(err)
		}

		pos = append(pos, po)
	}

	bos := make([]*biz.Notice, 0, len(pos))
	for _, po := range pos {
		bos = append(bos, noticePo2bo(po))
	}

	return bos, nil
}

func (d *noticeData) GetById(ctx context.Context, id int64) (*biz.Notice, error) {
	po := &Notice{}

	fb := pg.NewSelectSQLFieldBuilder()
	fb.Append("id", &po.Id)
	fb.Append("title", &po.Title)
	fb.Append("content", &po.Content)
	fb.Append("start_time", &po.StartTime)
	fb.Append("end_time", &po.EndTime)

	fieldSql, _ := fb.Build()

	sqlStr := fmt.Sprintf("SELECT %s FROM notices WHERE id=$1", fieldSql)

	row := d.data.Pdb.QueryRowContext(ctx, sqlStr, id)
	if err := row.Scan(&po.Id, &po.Title, &po.Content, &po.StartTime, &po.EndTime); err != nil {
		return nil, xerrors.APIDBFailed("id=%d", id).WithCause(err)
	}

	return noticePo2bo(po), nil
}

func (d *noticeData) UpdateById(ctx context.Context, bo *biz.Notice) error {
	po := noticeBo2Po(bo)

	fb := pg.NewUpdateSQLFieldBuilder(2)
	fb.Append("title", &po.Title)
	fb.Append("content", &po.Content)
	fb.Append("start_time", &po.StartTime)
	fb.Append("end_time", &po.EndTime)

	fieldSql, values := fb.Build()
	values = pg.AppendValueFirst(values, &po.Id)

	sqlStr := fmt.Sprintf("UPDATE notices SET %s WHERE id=$1", fieldSql)

	result, err := d.data.Pdb.ExecContext(ctx, sqlStr, values...)
	if err != nil {
		return xerrors.APIDBFailed("id=%d", bo.Id).WithCause(err)
	}

	row, err := result.RowsAffected()
	if err != nil {
		return xerrors.APIDBFailed("id=%d", bo.Id).WithCause(err)
	}

	if row != 1 {
		return xerrors.APINotFound("id=%d", bo.Id)
	}

	return nil
}

func (d *noticeData) Insert(ctx context.Context, bo *biz.Notice) error {
	po := noticeBo2Po(bo)

	fb := pg.NewInsertSQLFieldBuilder()
	fb.Append("title", &po.Title)
	fb.Append("content", &po.Content)
	fb.Append("start_time", &po.StartTime)
	fb.Append("end_time", &po.EndTime)

	fieldSql, argSql, values := fb.Build()

	sqlStr := fmt.Sprintf("INSERT INTO notices (%s) VALUES (%s)", fieldSql, argSql)

	result, err := d.data.Pdb.ExecContext(ctx, sqlStr, values...)
	if err != nil {
		return xerrors.APIDBFailed("id=%d", bo.Id).WithCause(err)
	}

	row, err := result.RowsAffected()
	if err != nil {
		return xerrors.APIDBFailed("id=%d", bo.Id).WithCause(err)
	}

	if row != 1 {
		return xerrors.APIDBFailed("not affected. id=%d", bo.Id)
	}

	return nil
}

func (d *noticeData) DeleteById(ctx context.Context, id int64) error {
	sqlStr := "DELETE FROM notices WHERE id=$1"

	result, err := d.data.Pdb.ExecContext(ctx, sqlStr, id)
	if err != nil {
		return xerrors.APIDBFailed("id=%d", id).WithCause(err)
	}

	row, err := result.RowsAffected()
	if err != nil {
		return xerrors.APIDBFailed("id=%d", id).WithCause(err)
	}

	if row != 1 {
		return xerrors.APINotFound("id=%d", id)
	}

	return nil
}

func (d *noticeData) Count(ctx context.Context) (count int64, err error) {
	sqlStr := "SELECT COUNT(*) FROM notices"

	row := d.data.Pdb.QueryRowContext(ctx, sqlStr)
	if err := row.Scan(&count); err != nil {
		return 0, xerrors.APIDBFailed("count failed").WithCause(err)
	}

	return count, nil
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

	return out
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

	return out
}
