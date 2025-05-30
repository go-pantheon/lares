package service

import (
	v1 "github.com/go-pantheon/lares/app/account/internal/admin/service/v1"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(v1.NewAccountAdmin)
