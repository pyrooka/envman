package githubgist

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Delete an environment.
func (g *GitHubGist) Delete(envName string, envVars []string) (err error) {
	// If no variable delete the whole gist file.
	if len(envVars) == 0 {
		err = deleteGistFile(g.Token, envName, g.EnvManGist)
	} else {
		err = deleteEnvVars(g.Token, envName, envVars, g.EnvManGist)
	}

	return
}

// Deletes a variable from the environment (gist file).
func deleteEnvVars(token string, envName string, envVars []string, envmanGist *gist) (err error) {
	// Check the environment.
	if env, exists := envmanGist.Files[envName]; exists {
		// Get the content of the file.
		content, err := getGistFileContent(env.URL, token)
		if err != nil {
			return err
		}

		// Delete the keys.
		for _, envVar := range envVars {
			delete(content, envVar)
		}

		// Create JSON from the map.
		contentJSON, err := json.Marshal(content)
		if err != nil {
			return err
		}
		// Set the content to this env.
		env.Content = string(contentJSON)

		// Edit the files.
		envmanGist.Files = map[string]*gistFile{
			envName: env,
		}
	} else {
		err = fmt.Errorf("environment %v not found", envName)
		return
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

// Deletes a while gist file (environment).
func deleteGistFile(token string, envName string, envmanGist *gist) (err error) {
	// INFO: envman is a reserved name.
	if strings.ToLower(envName) == reservedName {
		err = fmt.Errorf("%v is a reserved name in all variation (lower/uppercase)", reservedName)
		return
	}

	// Get the environment.
	if _, exists := envmanGist.Files[envName]; exists {
		// Set the name to an empty string.
		envmanGist.Files = map[string]*gistFile{
			envName: nil,
		}
	} else {
		err = fmt.Errorf("environment %v doesn't exists", envName)
		return
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
