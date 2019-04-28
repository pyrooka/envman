package backend

import (
	"fmt"

	"github.com/pyrooka/envman/config"
)

// Local uses the computer for backend storage.
type Local struct {
	Environments map[string]map[string]string `json:"environments"`
}

// Init loads the environments from the config.
func (l *Local) Init(c *config.Config) (err error) {
	// If the local config is null,
	if c.Local.Environments == nil {
		// init a new map object.
		c.Local.Environments = map[string]map[string]string{}
	}

	l.Environments = c.Local.Environments

	return
}

// List returns the name of environments or variables.
func (l *Local) List(envName string) (result []string, err error) {
	if envName == "" {
		// Get the name of the environments.
		for key := range l.Environments {
			result = append(result, key)
		}
	} else {
		// Get the name of the variables in the environment if it exists.
		if env, exists := (l.Environments)[envName]; exists {
			for key := range env {
				result = append(result, key)
			}
		} else {
			err = fmt.Errorf("environment \"%v\" doesn't exist", envName)
		}
	}

	return
}

// Get returns the environment variables with its values.
func (l *Local) Get(envName string) (vars map[string]string, err error) {
	if env, exists := l.Environments[envName]; exists {
		vars = env
	} else {
		err = fmt.Errorf("environment \"%v\" doesn't exist", envName)
	}

	return
}

// Update saves the given variable to the environments. Overwrites if exists.
func (l *Local) Update(envName string, variables map[string]string) (err error) {
	// Create the environment if not already exists.
	if _, exists := l.Environments[envName]; !exists {
		l.Environments[envName] = map[string]string{}
	}

	// Add the variables to the env.
	for key, value := range variables {
		l.Environments[envName][key] = value
	}

	return
}

// Delete removes an environment.
func (l *Local) Delete(envName string, envVars []string) (err error) {
	// We need the environment in both cases.
	if env, exists := l.Environments[envName]; exists {
		if len(envVars) == 0 {
			// Delete the environment.
			delete(l.Environments, envName)
		} else {
			// Delete variables.
			for _, envVar := range envVars {
				// NOTE: should we check is the variables exists?
				delete(env, envVar)
			}

		}
	} else {
		err = fmt.Errorf("environment \"%v\" doesn't exist", envName)
	}

	return
}

// CleanUp removes all the environments from the config.
func (l *Local) CleanUp() (err error) {
	for key := range l.Environments {
		delete(l.Environments, key)
	}
	return
}
