package local

import "fmt"

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
