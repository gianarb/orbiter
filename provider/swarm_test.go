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

func TestExcludeUnCountedTasks(t *testing.T) {
	p := SwarmProvider{}
	tasks := []swarm.Task{
		swarm.Task{
			Status: swarm.TaskStatus{
				State: swarm.TaskStatePreparing,
			},
		},
		swarm.Task{
			Status: swarm.TaskStatus{
				State: swarm.TaskStateRunning,
			},
		},
		swarm.Task{
			Status: swarm.TaskStatus{
				State: swarm.TaskStateFailed,
			},
		},
	}
	n := p.calculateActiveTasks(tasks)
	if n != 2 {
		t.Errorf("Tasks need to be two")
	}
}
