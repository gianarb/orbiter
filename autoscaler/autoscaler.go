package autoscaler

import (
	"context"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"time"
)

type Provider interface {
	Scale(string, int, bool) error
	Name() string
}

type Autoscalers map[string]Autoscaler

type Autoscaler struct {
	provider   Provider
	serviceId  string
	targetUp   int
	targetDown int
	CoolDown   int
}

func NewAutoscaler(p Provider, serviceId string, targetUp int, targetDown int, coolDown int) Autoscaler {
	a := Autoscaler{
		provider:   p,
		serviceId:  serviceId,
		targetUp:   targetUp,
		targetDown: targetDown,
		CoolDown:   coolDown,
	}
	return a
}

func canScale(serviceId string, coolDown time.Duration) (retval bool, err error) {
	retval = false
	err = nil
	ctx := context.Background()
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		logrus.WithField("error", err).Debug("Problem communication with Docker")
		return
	}

	services, err := dockerClient.ServiceList(ctx, types.ServiceListOptions{})
	if err != nil {
		logrus.WithField("error", err).Debug("Bad comunication with Docker.")
		return
	}

	for _, service := range services {
		if service.Spec.Name != serviceId {
			continue
		}
		// now < updatedAt + coolDown ??
		if time.Now().Before(service.Meta.UpdatedAt.Add(coolDown)) {
			err = errors.New(fmt.Sprintf("Cooldown period for %f seconds", service.Meta.UpdatedAt.Add(coolDown).Sub(time.Now()).Seconds()))
			break
		} else {
			retval = true
			break
		}
	}
	return
}

func (a *Autoscaler) ScaleUp() error {
	logrus.WithFields(logrus.Fields{
		"service":   a.serviceId,
		"direction": true,
	}).Infof("Received a new request to scale up %s with %d task.", a.serviceId, a.targetUp)

	if ok, err := canScale(a.serviceId, time.Duration(a.CoolDown)*time.Second); ok == false {
		logrus.Warn("Cannot scale up during coolDown period!")
		return err
	}

	err := a.provider.Scale(a.serviceId, a.targetUp, true)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"service":   a.serviceId,
			"direction": true,
			"error":     err.Error(),
		}).Warnf("We had some problems to scale up %s.", a.serviceId)
	} else {
		logrus.WithFields(logrus.Fields{
			"service":   a.serviceId,
			"direction": true,
		}).Infof("Service %s scaled up.", a.serviceId)
	}

	return err
}

func (a *Autoscaler) ScaleDown() error {
	logrus.WithFields(logrus.Fields{
		"service":   a.serviceId,
		"direction": false,
	}).Infof("Received a new request to scale down %s with %d task.", a.serviceId, a.targetDown)

	if ok, err := canScale(a.serviceId, time.Duration(a.CoolDown)*time.Second); ok == false {
		logrus.Warn("Cannot scale down during coolDown period!")
		return err
	}
	err := a.provider.Scale(a.serviceId, a.targetDown, false)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"service":   a.serviceId,
			"direction": false,
			"error":     err.Error(),
		}).Warnf("We had some problems to scale down %s.", a.serviceId)
	} else {
		logrus.WithFields(logrus.Fields{
			"service":   a.serviceId,
			"direction": false,
		}).Infof("Service %s scaled down.", a.serviceId)
	}
	return err
}
