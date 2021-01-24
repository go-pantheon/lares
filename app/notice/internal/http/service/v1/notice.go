package v1

import (
	"context"

	"github.com/luffy050596/rec-account/app/notice/internal/http/biz"
	v1 "github.com/luffy050596/rec-account/gen/api/server/notice/interface/notice/v1"
)

type NoticeInterface struct {
	v1.UnimplementedNoticeInterfaceServer

	log *log.Helper
	uc  *biz.NoticeUseCase
}

func NewNoticeInterface(logger log.Logger, uc *biz.NoticeUseCase) *NoticeInterface {
	return &NoticeInterface{
		log: log.NewHelper(log.With(logger, "module", "notice/interface/v1/notice")),
		uc:  uc,
	}
}

func (s *NoticeInterface) NoticeList(ctx context.Context, req *v1.NoticeListRequest) (*v1.NoticeListResponse, error) {
	reply := &v1.NoticeListResponse{
		Code: v1.NoticeListResponse_CODE_SUCCEEDED,
		List: make([]*v1.Notice, 0),
	}

	list, err := s.uc.NoticeList(ctx)
	if err != nil {
		s.log.WithContext(ctx).Errorf("[NoticeInterface.NoticeList] 失败: %v", err)
		return reply, nil
	}

	if len(list) == 0 {
		return reply, nil
	}

	reply.List = make([]*v1.Notice, 0, len(list))
	for _, v := range list {
		reply.List = append(reply.List, po2dto(v))
	}
	return reply, nil
}

func po2dto(po *biz.Notice) *v1.Notice {
	return &v1.Notice{
		Title:   po.Title,
		Content: po.Content,
	}
}
