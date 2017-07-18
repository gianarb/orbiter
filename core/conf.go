package core

import (
	"github.com/go-yaml/yaml"
)

type PolicyConf struct {
	Up       int `yaml:"up"`       // Number of tasks to start during a scale up
	Down     int `yaml:"down"`     // Number of tasks to start during a scale down
	CoolDown int `yaml:"cooldown"` // Number of milliseconds to sleep avoidin too quick scale
}

type AutoscalerConf struct {
	Provider   string                `yaml:"provider"`
	Parameters map[string]string     `yaml:"parameters"`
	Policies   map[string]PolicyConf `yaml:"policies"`
}

type Conf struct {
	//Daemon map[string]Idontknow `yaml:"daemon"`
	AutoscalersConf map[string]AutoscalerConf `yaml:"autoscalers"`
}

func createConfiguration() Conf {
	conf := Conf{
		AutoscalersConf: map[string]AutoscalerConf{},
	}
	return conf
}

func ParseYAMLConfiguration(content []byte) (Conf, error) {
	config := createConfiguration()
	err := yaml.Unmarshal(content, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
