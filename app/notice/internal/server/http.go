package server

import (
	"strings"

	adv1 "github.com/luffy050596/rec-account/app/notice/internal/admin/service/v1"
	"github.com/luffy050596/rec-account/app/notice/internal/conf"
	ifacev1 "github.com/luffy050596/rec-account/app/notice/internal/http/service/v1"
	adminv1 "github.com/luffy050596/rec-account/gen/api/server/notice/admin/notice/v1"
	interfacev1 "github.com/luffy050596/rec-account/gen/api/server/notice/interface/notice/v1"
	"github.com/luffy050596/rec-kit/metrics"
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
