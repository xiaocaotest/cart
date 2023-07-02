package common

import (
	"github.com/go-micro/plugins/v4/config/source/consul"
	"go-micro.dev/v4/config"
	"strconv"
)

func GetConsulConfig(host string, port int64, prefix string) (config.Config, error) {
	consulSource := consul.NewSource(
		consul.WithAddress(host+":"+strconv.FormatInt(port, 10)),
		consul.WithPrefix(prefix),
		consul.StripPrefix(true),
	)
	config, err := config.NewConfig()
	if err != nil {
		return config, err
	}

	err = config.Load(consulSource)
	return config, err
}
