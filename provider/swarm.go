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

func (p SwarmProvider) Name() string {
	return "swarm"
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
		}).Infof("Service %s is not scaling.", serviceId)
		return err
	}

	spec := service.Spec
	var ptrFromSystem uint64
	base := p.calculateActiveTasks(tasks)
	if direction == true {
		ptrFromSystem = uint64(base + target)
	} else {
		ptrFromSystem = uint64(base - target)
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
	}).Debugf("Service %s scaled from %d to %d", serviceId, base, ptrFromSystem)
	return nil
}

// This function validate if a request is acceptable or not.
func (p *SwarmProvider) isAcceptable(tasks []swarm.Task, target int, direction bool) error {
	if direction == false && (p.calculateActiveTasks(tasks) < target || p.calculateActiveTasks(tasks) < 2) {
		return errors.New(fmt.Sprintf("I can not scale down because it has only %d running.", target))
	}
	return nil
}

// Calculate the number of tasks to use as started poit to scale up or down.
// This function is necesarry because we need to exclude shutted down or
// rejected tasks.
func (p *SwarmProvider) calculateActiveTasks(tasks []swarm.Task) int {
	c := 0
	for _, task := range tasks {
		if task.Status.State == swarm.TaskStateNew ||
			task.Status.State == swarm.TaskStateAccepted ||
			task.Status.State == swarm.TaskStatePending ||
			task.Status.State == swarm.TaskStateAssigned ||
			task.Status.State == swarm.TaskStateStarting ||
			task.Status.State == swarm.TaskStatePreparing ||
			task.Status.State == swarm.TaskStateReady ||
			task.Status.State == swarm.TaskStateRunning {
			c = c + 1
		}
	}
	return c
}
