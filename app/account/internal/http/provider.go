package http

import (
	"github.com/go-pantheon/lares/app/account/internal/http/biz"
	"github.com/go-pantheon/lares/app/account/internal/http/data"
	"github.com/go-pantheon/lares/app/account/internal/http/domain"
	"github.com/go-pantheon/lares/app/account/internal/http/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(service.ProviderSet, biz.ProviderSet, domain.ProviderSet, data.ProviderSet)
