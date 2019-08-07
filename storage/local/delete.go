package local

import "fmt"

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
