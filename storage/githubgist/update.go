package githubgist

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Update variables in the environment.
func (g *GitHubGist) Update(envName string, variables map[string]string) (err error) {
	err = updateGist(g.Token, envName, variables, g.EnvManGist)

	return
}

// Updates the gist.
func updateGist(token string, envName string, envVars map[string]string, envmanGist *gist) (err error) {
	// INFO: envman is a reserved name.
	if strings.ToLower(envName) == reservedName {
		err = fmt.Errorf("%v is a reserved name in all variation (lower/uppercase)", reservedName)
		return
	}

	// Get or create the file contains our variables.
	if env, exists := envmanGist.Files[envName]; exists {
		// So the gist file (environment) exists.
		// Get the content of the file.
		content, err := getGistFileContent(env.URL, token)
		if err != nil {
			return err
		}

		// Edit the content.
		for key, value := range envVars {
			content[key] = value
		}

		// Create JSON again.
		contentJSON, err := json.Marshal(content)
		if err != nil {
			return err
		}

		// Set the content to the current env (a gist file).
		env.Content = string(contentJSON)

		// Set this the only file in the gist.
		envmanGist.Files = map[string]*gistFile{
			envName: env,
		}

	} else {
		// This means no gist file (environment) found with the name.
		// Create the JSON content.
		contentJSON, err := json.Marshal(envVars)
		if err != nil {
			return err
		}

		// Create a new gist file.
		env := &gistFile{
			Content: string(contentJSON),
		}

		// Set this file to the only one in the gist.
		envmanGist.Files = map[string]*gistFile{
			envName: env,
		}
	}

	// Create json from the gist struct.
	body, err := json.Marshal(&envmanGist)
	if err != nil {
		return
	}

	// Let's patch the old one.
	_, err = makePatch(envmanGist.URL, token, body)

	return
}

// Deletes a gist.
func deleteGist(token string, envmanGist *gist) (err error) {
	err = makeDelete(envmanGist.URL, token)

	return
}
