package http

import (
	"github.com/google/wire"
	"github.com/luffy050596/rec-account/app/account/internal/http/biz"
	"github.com/luffy050596/rec-account/app/account/internal/http/data"
	"github.com/luffy050596/rec-account/app/account/internal/http/domain"
	"github.com/luffy050596/rec-account/app/account/internal/http/service"
)

var ProviderSet = wire.NewSet(service.ProviderSet, biz.ProviderSet, domain.ProviderSet, data.ProviderSet)
