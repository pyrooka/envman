package local

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
