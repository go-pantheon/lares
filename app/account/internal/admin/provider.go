package admin

import (
	"github.com/go-pantheon/lares/app/account/internal/admin/biz"
	"github.com/go-pantheon/lares/app/account/internal/admin/data"
	"github.com/go-pantheon/lares/app/account/internal/admin/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(service.ProviderSet, biz.ProviderSet, data.ProviderSet)
