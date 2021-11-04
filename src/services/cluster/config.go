package cluster

import (
	"github.com/linshenqi/sptty"
	"github.com/linshenqi/taskmate/src/base"
)

type Config struct {
	sptty.BaseConfig

	Enable   bool     `yaml:"enable"`
	EtcdURLs []string `yaml:"etcdURLs"`
}

func (s *Config) ConfigName() string {
	return base.ServiceCluster
}
