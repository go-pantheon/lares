package server

import (
	"github.com/google/wire"
	"github.com/luffy050596/rec-account/app/account/internal/conf"
	"github.com/pkg/errors"
)

var ProviderSet = wire.NewSet(NewHTTPServer, NewRegistrar)

func NewRegistrar(conf *conf.Registry) (registry.Registrar, error) {
	client, err := etcdclient.New(etcdclient.Config{
		Endpoints: conf.Etcd.Endpoints,
		Username:  conf.Etcd.Username,
		Password:  conf.Etcd.Password,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "etcd client create failed")
	}

	return etcd.New(client), nil
}
