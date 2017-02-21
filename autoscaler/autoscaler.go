package autoscaler

import "github.com/Sirupsen/logrus"

type Provider interface {
	Scale(string, int, bool) error
}

type Autoscalers map[string]Autoscaler

type Autoscaler struct {
	provider   Provider
	serviceId  string
	targetUp   int
	targetDown int
}

func NewAutoscaler(p Provider, serviceId string, targetUp int, targetDown int) Autoscaler {
	a := Autoscaler{
		provider:   p,
		serviceId:  serviceId,
		targetUp:   targetUp,
		targetDown: targetDown,
	}
	return a
}

func (a *Autoscaler) ScaleUp() error {
	logrus.WithFields(logrus.Fields{
		"service":   a.serviceId,
		"direction": true,
	}).Infof("Received a new request to scale up %s with %d task.", a.serviceId, a.targetUp)

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
