package githubgist

import "errors"

// CleanUp removes all the created thing. Gist here.
func (g *GitHubGist) CleanUp() (err error) {
	// Delete the gist first.
	err = deleteGist(g.Token, g.EnvManGist)
	if err != nil {
		return
	}

	// Now the authorization.
	auth, err := getAuth(g.Token)
	if err != nil {
		return
	} else if auth == nil {
		err = errors.New("no authorization found")
		return
	}

	err = deleteAuth(auth.URL, g.Token)

	return
}
