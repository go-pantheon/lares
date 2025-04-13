package http

import (
	"github.com/go-pantheon/lares/app/notice/internal/http/biz"
	"github.com/go-pantheon/lares/app/notice/internal/http/data"
	"github.com/go-pantheon/lares/app/notice/internal/http/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(service.ProviderSet, biz.ProviderSet, data.ProviderSet)
