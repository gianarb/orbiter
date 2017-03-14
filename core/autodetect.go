package core

import "github.com/gianarb/orbiter/autoscaler"

// This function use diferent strategies to get information from
// the system itself to configure the autoloader.
// They can be environment variables for example or other systems.
func Autodetect() (Core error) {
	scalers := autoscaler.Autoscalers{}
	var core Core
	autoDetectSwarmMode(&scalers)
	return core, nil
}

func autoDetectSwarmMode(a *autoscaler.Autoscalers) {
	// Create Docker Client by EnvVar and check if it's working.

	// Get List of Services

	// Check which services has labels and register them to orbiter.
}
