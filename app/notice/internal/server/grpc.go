package server

import (
	"math"

	v1 "github.com/luffy050596/rec-account/app/notice/internal/admin/service/v1"
	"github.com/luffy050596/rec-account/app/notice/internal/conf"
	adminv1 "github.com/luffy050596/rec-account/gen/api/server/notice/admin/notice/v1"
	"github.com/luffy050596/rec-kit/metrics"
	"google.golang.org/grpc"
)

func NewGRPCServer(c *conf.Server, logger log.Logger,
	noticeAdminV1 *v1.NoticeAdmin,
) *kgrpc.Server {
	var opts = []kgrpc.ServerOption{
		kgrpc.Middleware(
			middleware.Chain(
				recovery.Recovery(),
				metadata.Server(),
				tracing.Server(),
				metrics.Server(),
				logging.Server(logger),
			),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, kgrpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, kgrpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, kgrpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	opts = append(opts, kgrpc.Options(
		grpc.InitialConnWindowSize(1<<30),
		grpc.InitialWindowSize(1<<30),
		grpc.MaxConcurrentStreams(math.MaxInt32),
	))
	srv := kgrpc.NewServer(opts...)
	adminv1.RegisterNoticeAdminServer(srv, noticeAdminV1)
	return srv
}
