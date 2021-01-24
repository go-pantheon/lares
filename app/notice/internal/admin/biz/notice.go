package biz

import (
	"context"
	"time"
)

type Notice struct {
	Id        int64
	Title     string
	Content   string
	Sort      int64
	StartTime time.Time
	EndTime   time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

type NoticeRepo interface {
	GetList(ctx context.Context, index, size int64) ([]*Notice, error)
	Count(ctx context.Context) (int64, error)
	GetById(ctx context.Context, id int64) (*Notice, error)
	UpdateById(ctx context.Context, notice *Notice) error
	Insert(ctx context.Context, notice *Notice) error
	DeleteById(ctx context.Context, id int64) error
}

type NoticeUseCase struct {
	log *log.Helper

	repo NoticeRepo
}

func NewNoticeUseCase(repo NoticeRepo, logger log.Logger) *NoticeUseCase {
	return &NoticeUseCase{
		log:  log.NewHelper(log.With(logger, "module", "notice/admin/biz/notice")),
		repo: repo,
	}
}

func (uc *NoticeUseCase) GetList(ctx context.Context, index, size int64) (list []*Notice, count int64, err error) {
	list, err = uc.repo.GetList(ctx, index, size)
	if err != nil {
		return
	}
	count, err = uc.repo.Count(ctx)
	if err != nil {
		return
	}
	return
}

func (uc *NoticeUseCase) Count(ctx context.Context) (int64, error) {
	return uc.repo.Count(ctx)
}

func (uc *NoticeUseCase) GetById(ctx context.Context, id int64) (*Notice, error) {
	return uc.repo.GetById(ctx, id)
}

func (uc *NoticeUseCase) Update(ctx context.Context, notice *Notice) error {
	return uc.repo.UpdateById(ctx, notice)
}

func (uc *NoticeUseCase) Create(ctx context.Context, notice *Notice) error {
	return uc.repo.Insert(ctx, notice)
}

func (uc *NoticeUseCase) Delete(ctx context.Context, id int64) error {
	return uc.repo.DeleteById(ctx, id)
}
