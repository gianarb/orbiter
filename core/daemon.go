package core

import (
	"fmt"

	"github.com/gianarb/orbiter/autoscaler"
	"github.com/gianarb/orbiter/provider"
)

type Core struct {
	Autoscalers autoscaler.Autoscalers
}

func NewCoreByConfig(c map[string]AutoscalerConf, core *Core) error {
	scalers := autoscaler.Autoscalers{}
	for scalerName, scaler := range c {
		p, err := provider.NewProvider(scaler.Provider, scaler.Parameters)
		if err != nil {
			return err
		}
		for serviceName, policy := range scaler.Policies {
			scalers[fmt.Sprintf("%s/%s", scalerName, serviceName)] = autoscaler.NewAutoscaler(p, serviceName, policy.Up, policy.Down)
		}
	}
	core.Autoscalers = scalers
	return nil
}
