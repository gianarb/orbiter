package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	docker "github.com/docker/docker/client"
	"github.com/gianarb/orbiter/autoscaler"
)

type SwarmProvider struct {
	dockerClient *docker.Client
}

func NewSwarmProvider(c map[string]string) (autoscaler.Provider, error) {
	var p autoscaler.Provider
	client, err := docker.NewEnvClient()
	if err != nil {
		logrus.WithField("error", err).Warn("problem to communicate with docker")
		return p, err
	} else {
		logrus.Info("Successfully connected to a Docker daemon")
	}
	p = SwarmProvider{
		dockerClient: client,
	}
	return p, nil

}

func (p SwarmProvider) Scale(serviceId string, target int, direction bool) error {
	ctx := context.Background()
	service, _, err := p.dockerClient.ServiceInspectWithRaw(ctx, serviceId)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err.Error(),
			"provider": "swarm",
		}).Debugf("Service %s didn't scale. We didn't get it from docker.", serviceId)
		return err
	}

	filters := filters.NewArgs()
	filters.Add("service", serviceId)
	tasks, err := p.dockerClient.TaskList(ctx, types.TaskListOptions{
		Filters: filters,
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err.Error(),
			"provider": "swarm",
		}).Debugf("Service %s didn't scale. Impossibile to get current number of running tasks.", serviceId)
		return err
	}

	err = p.isAcceptable(tasks, target, direction)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err.Error(),
			"provider": "swarm",
		}).Info("Service %s is not scaling.", serviceId)
		return err
	}

	spec := service.Spec
	var ptrFromSystem uint64
	if direction == true {
		ptrFromSystem = uint64(len(tasks) + target)
	} else {
		ptrFromSystem = uint64(len(tasks) - target)
	}
	spec.Mode.Replicated.Replicas = &ptrFromSystem
	_, err = p.dockerClient.ServiceUpdate(ctx, serviceId, service.Version, spec, types.ServiceUpdateOptions{})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":    err.Error(),
			"provider": "swarm",
		}).Debugf("We had some trouble to updated %s on docker", serviceId)
		return err
	}

	logrus.WithFields(logrus.Fields{
		"provider": "swarm",
	}).Debugf("Service %s scaled.", serviceId)
	return nil
}

func (p *SwarmProvider) isAcceptable(tasks []swarm.Task, target int, direction bool) error {
	if len(tasks) < target && direction == false {
		return errors.New(fmt.Sprintf("I can not scale down because it has only %d running.", target))
	}
	return nil
}
