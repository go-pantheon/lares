package v1

import (
	"context"
	"net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/fabrica-util/xtime"
	"github.com/go-pantheon/lares/app/notice/internal/admin/biz"
	adminv1 "github.com/go-pantheon/lares/gen/api/server/notice/admin/notice/v1"
)

type NoticeAdmin struct {
	adminv1.UnimplementedNoticeAdminServer

	log *log.Helper
	uc  *biz.NoticeUseCase
}

func NewNoticeAdmin(logger log.Logger, uc *biz.NoticeUseCase) *NoticeAdmin {
	return &NoticeAdmin{
		log: log.NewHelper(log.With(logger, "module", "notice/admin/notice/v1")),
		uc:  uc,
	}
}

func (s *NoticeAdmin) GetById(ctx context.Context, req *adminv1.GetByIdRequest) (*adminv1.GetByIdResponse, error) {
	item, err := s.uc.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	reply := &adminv1.GetByIdResponse{
		Code: http.StatusOK,
		Item: noticeBo2dto(item),
	}
	return reply, nil
}

func (s *NoticeAdmin) GetList(ctx context.Context, req *adminv1.GetListRequest) (*adminv1.GetListResponse, error) {
	start, limit, err := buildGetListCond(req)
	if err != nil {
		return nil, err
	}

	bos, count, err := s.uc.GetList(ctx, start, limit)
	if err != nil {
		return nil, err
	}

	reply := &adminv1.GetListResponse{
		Code:  http.StatusOK,
		List:  make([]*adminv1.NoticeProto, 0, count),
		Total: count,
	}

	for _, bo := range bos {
		reply.List = append(reply.List, noticeBo2dto(bo))
	}
	return reply, nil
}

func (s *NoticeAdmin) Create(ctx context.Context, req *adminv1.CreateRequest) (*adminv1.CreateResponse, error) {
	if err := checkParam(req.Item); err != nil {
		return nil, err
	}

	if err := s.uc.Create(ctx, noticeDto2Bo(req.Item)); err != nil {
		return nil, err
	}
	return &adminv1.CreateResponse{
		Code: http.StatusOK,
	}, nil
}

func (s *NoticeAdmin) Update(ctx context.Context, req *adminv1.UpdateRequest) (*adminv1.UpdateResponse, error) {
	if err := checkParam(req.Item); err != nil {
		return nil, err
	}

	if err := s.uc.Update(ctx, noticeDto2Bo(req.Item)); err != nil {
		return nil, err
	}
	return &adminv1.UpdateResponse{
		Code: http.StatusOK,
	}, nil
}

func (s *NoticeAdmin) Delete(ctx context.Context, req *adminv1.DeleteRequest) (*adminv1.DeleteResponse, error) {
	if err := s.uc.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &adminv1.DeleteResponse{
		Code: http.StatusOK,
	}, nil
}

func buildGetListCond(req *adminv1.GetListRequest) (start, limit int64, err error) {
	if req.PageSize > profile.MaxPageSize {
		err = xerrors.APIPageParamInvalid("page size too large")
		return
	}
	start, limit = profile.PageStartLimit(int64(req.Page), int64(req.PageSize))
	return
}

func noticeBo2dto(bo *biz.Notice) *adminv1.NoticeProto {
	dto := &adminv1.NoticeProto{
		Id:          bo.Id,
		Title:       bo.Title,
		Content:     bo.Content,
		Sort:        bo.Sort,
		StartTime:   bo.StartTime.Unix(),
		EndTime:     bo.EndTime.Unix(),
		CreatedTime: bo.CreatedAt.Unix(),
		UpdatedTime: bo.UpdatedAt.Unix(),
	}
	return dto
}

func noticeDto2Bo(dto *adminv1.NoticeProto) *biz.Notice {
	o := &biz.Notice{
		Id:        dto.Id,
		Title:     dto.Title,
		Content:   dto.Content,
		Sort:      dto.Sort,
		StartTime: xtime.Time(dto.StartTime),
		EndTime:   xtime.Time(dto.EndTime),
		CreatedAt: xtime.Time(dto.CreatedTime),
		UpdatedAt: xtime.Time(dto.UpdatedTime),
	}
	return o
}

func checkParam(dto *adminv1.NoticeProto) error {
	if dto.Title == "" {
		return xerrors.APIParamInvalid("title is empty")
	}
	if dto.Content == "" {
		return xerrors.APIParamInvalid("content is empty")
	}
	if dto.StartTime >= dto.EndTime {
		return xerrors.APIParamInvalid("time invalid")
	}
	return nil
}
