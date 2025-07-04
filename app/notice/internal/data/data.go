package data

import (
	"context"

	kdb "github.com/go-pantheon/fabrica-kit/trace/postgresql"
	udb "github.com/go-pantheon/fabrica-util/data/db/postgresql"
	"github.com/go-pantheon/fabrica-util/data/redis"
	"github.com/go-pantheon/lares/app/notice/internal/conf"
	"github.com/pkg/errors"
	goredis "github.com/redis/go-redis/v9"
)

type Data struct {
	Pdb *udb.DB
	Rdb goredis.UniversalClient
}

func NewData(c *conf.Data) (d *Data, cleanup func(), err error) {
	var (
		pdb *udb.DB
		rdb goredis.UniversalClient
	)

	if c.Redis.Cluster {
		rdb, cleanup, err = redis.NewCluster(&goredis.ClusterOptions{
			Addrs:        []string{c.Redis.Addr},
			Password:     c.Redis.Password,
			DialTimeout:  c.Redis.DialTimeout.AsDuration(),
			WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
			ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		})
	} else {
		rdb, cleanup, err = redis.NewStandalone(&goredis.Options{
			Addr:         c.Redis.Addr,
			Password:     c.Redis.Password,
			DialTimeout:  c.Redis.DialTimeout.AsDuration(),
			WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
			ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		})
	}

	if err != nil {
		return nil, nil, err
	}

	pdb, cleanupPg, err := kdb.NewTracingDB(context.Background(), kdb.DefaultPostgreSQLConfig(udb.Config{
		DSN:      c.Postgresql.Source,
		DBName:   c.Postgresql.Database,
		MinConns: 10,
		MaxConns: 100,
	}))
	if err != nil {
		return nil, nil, errors.Wrapf(err, "new postgres db failed")
	}

	d = &Data{
		Pdb: pdb,
		Rdb: rdb,
	}

	combinedCleanup := func() {
		if cleanup != nil {
			cleanup()
		}

		if cleanupPg != nil {
			cleanupPg()
		}
	}

	return d, combinedCleanup, nil
}
