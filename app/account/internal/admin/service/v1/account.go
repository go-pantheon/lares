package v1

import (
	"context"
	"net/http"

	"github.com/luffy050596/rec-account/app/account/internal/admin/biz"
	"github.com/luffy050596/rec-account/app/account/internal/pkg/i64"
	adminv1 "github.com/luffy050596/rec-account/gen/api/server/account/admin/account/v1"
	"github.com/luffy050596/rec-kit/profile"
	"github.com/luffy050596/rec-kit/xerrors"
)

type AccountAdmin struct {
	adminv1.UnimplementedAccountAdminServer

	log *log.Helper
	ac  *biz.AccountUseCase
}

func NewAccountAdmin(logger log.Logger, ac *biz.AccountUseCase) *AccountAdmin {
	return &AccountAdmin{
		log: log.NewHelper(log.With(logger, "module", "account/admin/account/v1")),
		ac:  ac,
	}
}

func (s *AccountAdmin) GetList(ctx context.Context, req *adminv1.ListRequest) (reply *adminv1.ListResponse, err error) {
	cond, page, pageSize, err := buildListCond(req)
	if err != nil {
		return nil, err
	}

	accounts, count, err := s.ac.GetList(ctx, i64.Max(int64(page-1), 0)*pageSize, pageSize, cond)
	if err != nil {
		return nil, err
	}

	reply = &adminv1.ListResponse{
		Code:     http.StatusOK,
		Accounts: make([]*adminv1.AccountProto, 0, count),
		Total:    count,
	}

	for _, bo := range accounts {
		reply.Accounts = append(reply.Accounts, accountBo2dto(bo))
	}
	return reply, nil
}

func buildListCond(req *adminv1.ListRequest) (cond *biz.Account, start, limit int64, err error) {
	if req.PageSize > profile.MaxPageSize {
		err = xerrors.APIPageParamInvalid("page size <= %d", profile.MaxPageSize)
		return
	}

	start, limit = profile.PageStartLimit(req.Page, req.PageSize)

	cond = &biz.Account{}
	if req.Condition == nil {
		return
	}

	if len(req.Condition.Name) > 0 {
		cond.Name = req.Condition.Name
	}
	if len(req.Condition.AppleId) > 0 {
		cond.AppleId = req.Condition.AppleId
	}
	if len(req.Condition.GoogleId) > 0 {
		cond.GoogleId = req.Condition.GoogleId
	}
	return
}

func (s *AccountAdmin) GetById(ctx context.Context, req *adminv1.GetByIdRequest) (*adminv1.GetByIdResponse, error) {
	o, err := s.ac.GetById(ctx, req.Id)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	return &adminv1.GetByIdResponse{
		Account: accountBo2dto(o),
	}, nil
}

func accountBo2dto(bo *biz.Account) *adminv1.AccountProto {
	dto := &adminv1.AccountProto{
		Id:          bo.Id,
		IdStr:       bo.IdStr,
		AppleId:     bo.AppleId,
		GoogleId:    bo.GoogleId,
		Name:        bo.Name,
		Channel:     bo.Channel,
		RegisterIp:  bo.RegisterIp,
		LastLoginIp: bo.LastLoginIp,
		CreatedAt:   bo.CreatedAt.Unix(),
	}
	return dto
}
