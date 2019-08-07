package githubgist

// Init makes the authentication if necessary.
func (g *GitHubGist) Init() (err error) {
	/*
		var token string
		// Check if we have auth token.
		if token = c.GitHubGist.Token; len(token) > 0 {
			// If have, test is.
			err = testToken(token)
			// If no error occured it means we got HTTP 200.
			if err != nil {
				// Clear the error.
				err = nil
				fmt.Println("Invalid token. Please create a new one.")
			}
		} else {
			// Otherwise we need to authenticate the user and create the token.
			fmt.Println("No token found for authentication. Please login.")
		}

		if token == "" {
			token, err = login()
			if err != nil {
				return
			}
		}

		// Set the token to the struct and the config.
		g.Token = token
		c.GitHubGist.Token = token

		// Load our gist to the struct.
		envGist, err := getOrCreateGist(token)
		if err != nil {
			return
		}

		g.EnvManGist = envGist
	*/
	return
}
