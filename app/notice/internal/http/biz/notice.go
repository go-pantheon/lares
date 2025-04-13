package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type Notice struct {
	Title     string
	Content   string
	StartTime time.Time
	EndTime   time.Time
}

type NoticeRepo interface {
	NoticeList(ctx context.Context, now time.Time) ([]*Notice, error)
}

type NoticeUseCase struct {
	log *log.Helper

	repo NoticeRepo
}

func NewNoticeUseCase(repo NoticeRepo, logger log.Logger) *NoticeUseCase {
	return &NoticeUseCase{
		log:  log.NewHelper(log.With(logger, "module", "account/biz/notice")),
		repo: repo,
	}
}

func (uc *NoticeUseCase) NoticeList(ctx context.Context) ([]*Notice, error) {
	return uc.repo.NoticeList(ctx, time.Now())
}
