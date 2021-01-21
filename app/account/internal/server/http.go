package server

import (
	"strings"

	adv1 "github.com/luffy050596/rec-account/app/account/internal/admin/service/v1"
	"github.com/luffy050596/rec-account/app/account/internal/conf"
	ifacev1 "github.com/luffy050596/rec-account/app/account/internal/http/service/v1"
	adminv1 "github.com/luffy050596/rec-account/gen/api/server/account/admin/account/v1"
	interfacev1 "github.com/luffy050596/rec-account/gen/api/server/account/interface/account/v1"
	"github.com/luffy050596/rec-kit/metrics"
)

func NewHTTPServer(label *conf.Label, c *conf.Server, logger log.Logger,
	v1Iface *ifacev1.AccountInterface,
	v1Admin *adv1.AccountAdmin,
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

	interfacev1.RegisterAccountInterfaceHTTPServer(s, v1Iface)
	if strings.ToLower(label.Profile) == "dev" {
		adminv1.RegisterAccountAdminHTTPServer(s, v1Admin)
	}
	return s
}
