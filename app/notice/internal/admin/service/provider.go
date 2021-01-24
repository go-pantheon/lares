package service

import (
	"github.com/google/wire"
	v1 "github.com/luffy050596/rec-account/app/notice/internal/admin/service/v1"
)

var ProviderSet = wire.NewSet(v1.NewNoticeAdmin)
