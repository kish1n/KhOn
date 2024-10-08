package config

import (
	"time"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
)

type jwtRaw struct {
	SecretKey      string        `fig:"secret_key,required"`
	ExpirationTime time.Duration `fig:"expiration_time,required"`
}

type JWT struct {
	SecretKey      string
	ExpirationTime time.Duration
}

func (c *config) JWT() *JWT {
	return c.jwt.Do(func() interface{} {
		cfgRaw := jwtRaw{}
		err := figure.
			Out(&cfgRaw).
			From(kv.MustGetStringMap(c.getter, "jwt")).
			Please()
		if err != nil {
			panic(errors.WithMessage(err, "failed to figure out"))
		}

		return &JWT{
			SecretKey:      cfgRaw.SecretKey,
			ExpirationTime: cfgRaw.ExpirationTime,
		}
	}).(*JWT)
}
