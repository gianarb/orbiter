package core

import "testing"

func TestGetRightOrbiterUp(t *testing.T) {
	r := convertStringLabelToInt("orbiter.up", map[string]string{
		"orbiter.up":   "3",
		"orbiter.down": "3",
	})
	if r != 3 {
		t.Fatalf("%d != 3", r)
	}
}

func TestGetRightOrbiterDown(t *testing.T) {
	r := convertStringLabelToInt("orbiter.down", map[string]string{
		"orbiter.up":   "3",
		"orbiter.down": "2",
	})
	if r != 2 {
		t.Fatalf("%d != 2", r)
	}
}
