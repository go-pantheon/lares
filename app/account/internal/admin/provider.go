package admin

import (
	"github.com/google/wire"
	"github.com/luffy050596/rec-account/app/account/internal/admin/biz"
	"github.com/luffy050596/rec-account/app/account/internal/admin/data"
	"github.com/luffy050596/rec-account/app/account/internal/admin/service"
)

var ProviderSet = wire.NewSet(service.ProviderSet, biz.ProviderSet, data.ProviderSet)
