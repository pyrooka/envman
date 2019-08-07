package githubgist

// List the environments or variables.
func (g *GitHubGist) List(envName string) (result []string, err error) {
	// If no env name given, list the environments.
	if envName == "" {
		result = getEnvironments(g.EnvManGist)
	} else {
		result, err = getVariables(g.EnvManGist, envName, g.Token)
	}

	return
}
