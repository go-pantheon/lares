package data

import (
	"github.com/go-pantheon/fabrica-util/data/cache"
	"github.com/go-pantheon/lares/app/notice/internal/conf"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Data struct {
	Mdb *gorm.DB
	Rdb cache.Cacheable
}

func NewData(c *conf.Data) (d *Data, cleanup func(), err error) {
	var (
		pdb *gorm.DB
		rdb cache.Cacheable
	)

	if c.Redis.Cluster {
		rdb, cleanup, err = cache.NewRedisCluster(&redis.ClusterOptions{
			Addrs:        []string{c.Redis.Addr},
			Password:     c.Redis.Password,
			DialTimeout:  c.Redis.DialTimeout.AsDuration(),
			WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
			ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		})
		if err != nil {
			return
		}
	} else {
		rdb, cleanup, err = cache.NewRedis(&redis.Options{
			Addr:         c.Redis.Addr,
			Password:     c.Redis.Password,
			DialTimeout:  c.Redis.DialTimeout.AsDuration(),
			WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
			ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		})
		if err != nil {
			return
		}
	}

	pdb, err = gorm.Open(postgres.New(postgres.Config{
		DSN: c.Postgres.Source,
	}), &gorm.Config{})
	if err != nil {
		err = errors.Wrapf(err, "new postgres db failed")
		return
	}

	sdb, err := pdb.DB()
	if err != nil {
		err = errors.Wrapf(err, "get raw db failed")
		return
	}

	sdb.SetMaxIdleConns(int(c.Postgres.MaxIdleConns))
	sdb.SetMaxOpenConns(int(c.Postgres.MaxOpenConns))

	d = &Data{
		Mdb: pdb,
		Rdb: rdb,
	}
	return
}
