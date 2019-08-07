package local

import "fmt"

// Get returns the environment variables with its values.
func (l *Local) Get(envName string) (vars map[string]string, err error) {
	if env, exists := l.Environments[envName]; exists {
		vars = env
	} else {
		err = fmt.Errorf("environment \"%v\" doesn't exist", envName)
	}

	return
}
