package server

import (
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-pantheon/fabrica-kit/metrics"
	adv1 "github.com/go-pantheon/lares/app/notice/internal/admin/service/v1"
	"github.com/go-pantheon/lares/app/notice/internal/conf"
	ifacev1 "github.com/go-pantheon/lares/app/notice/internal/http/service/v1"
	adminv1 "github.com/go-pantheon/lares/gen/api/server/notice/admin/notice/v1"
	interfacev1 "github.com/go-pantheon/lares/gen/api/server/notice/interface/notice/v1"
)

func NewHTTPServer(label *conf.Label, c *conf.Server, logger log.Logger,
	noticeIfaceV1 *ifacev1.NoticeInterface,
	noticeAdminV1 *adv1.NoticeAdmin,
) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			middleware.Chain(
				recovery.Recovery(),
				metadata.Server(),
				tracing.Server(),
				metrics.Server(),
				logging.Server(logger),
			)),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	s := http.NewServer(opts...)

	interfacev1.RegisterNoticeInterfaceHTTPServer(s, noticeIfaceV1)
	if strings.ToLower(label.Profile) == "dev" {
		adminv1.RegisterNoticeAdminHTTPServer(s, noticeAdminV1)
	}
	return s
}
