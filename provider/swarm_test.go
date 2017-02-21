package provider

import (
	"testing"

	"github.com/docker/docker/api/types/swarm"
)

func TestTaskCannotScaleDownBecauseItHasNotEnoughTaskRunning(t *testing.T) {
	p := SwarmProvider{}
	tasks := []swarm.Task{
		swarm.Task{},
	}
	err := p.isAcceptable(tasks, 3, false)
	if err == nil {
		t.Errorf("It can not scale")
	}
}
