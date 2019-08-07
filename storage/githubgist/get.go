package githubgist

import "errors"

// Get the environment from the gist and returns it as a map.
func (g *GitHubGist) Get(envName string) (vars map[string]string, err error) {
	// Get the environment from the map.
	env, exists := g.EnvManGist.Files[envName]
	if !exists {
		return nil, errors.New("environment not found")
	}

	// Get the content of the gist file.
	vars, err = getGistFileContent(env.URL, g.Token)

	return
}
