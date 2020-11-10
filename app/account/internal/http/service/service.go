package service

import (
	"github.com/google/wire"
	v1 "github.com/luffy050596/rec-account/app/account/internal/http/service/v1"
)

var ProviderSet = wire.NewSet(v1.NewAccountInterface)
