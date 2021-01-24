package http

import (
	"github.com/google/wire"
	"github.com/luffy050596/rec-account/app/notice/internal/http/biz"
	"github.com/luffy050596/rec-account/app/notice/internal/http/data"
	"github.com/luffy050596/rec-account/app/notice/internal/http/service"
)

var ProviderSet = wire.NewSet(service.ProviderSet, biz.ProviderSet, data.ProviderSet)
