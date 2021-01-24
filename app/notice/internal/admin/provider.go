package admin

import (
	"github.com/google/wire"
	"github.com/luffy050596/rec-account/app/notice/internal/admin/biz"
	"github.com/luffy050596/rec-account/app/notice/internal/admin/data"
	"github.com/luffy050596/rec-account/app/notice/internal/admin/service"
)

var ProviderSet = wire.NewSet(service.ProviderSet, biz.ProviderSet, data.ProviderSet)
