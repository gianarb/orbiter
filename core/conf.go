package core

import (
	"github.com/go-yaml/yaml"
)

type PolicyConf struct {
	// Number of tasks to start during a scale up
	Up int `yaml:"up"`
	// Number of tasks to start during a scale down
	Down int `yaml:"down"`
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
