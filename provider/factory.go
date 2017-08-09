package provider

import (
	"errors"
	"fmt"

	"github.com/gianarb/orbiter/autoscaler"
)

func NewProvider(t string, c map[string]string) (autoscaler.Provider, error) {
	var p autoscaler.Provider
	var err error
	switch t {
	case "swarm":
		p, err = NewSwarmProvider(c)
	case "aws_ec2":
		p, err = NewAwsEc2Provider(c)
	case "digitalocean":
		p, err = NewDigitalOceanProvider(c)
	case "fake":
		p = FakeProvider{}
	default:
		err = errors.New(fmt.Sprintf("%s not supported.", t))
	}
	return p, err
}
