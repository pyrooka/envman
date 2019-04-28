package backend

// GitHub Gist backend. https://gist.github.com

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/pyrooka/envman/config"
)

// API URLs for requests.
const (
	apiURL     = "https://api.github.com"
	userAPIURL = apiURL + "/user"
	gistAPIURL = apiURL + "/gists"
	authAPIURL = apiURL + "/authorizations"
)

// GitHub related constants.
const (
	gistDescription = "Envman Data"
	authTokenNote   = "Envman @ "
	reservedName    = "envman"
)

//-------------------------------------------------------------------
// Structs
//-------------------------------------------------------------------

// GitHubGist uses the Gist service of GitHub for backend storage.
type GitHubGist struct {
	Token         string
	EnvManGist *gist
}

// A file in the gist.
type gistFile struct {
	Filename string `json:"filename,omitempty"`
	URL      string `json:"raw_url,omitempty"`
	Content  string `json:"content,omitempty"`
}

// The gist.
type gist struct {
	URL         string               `json:"url,omitempty"`
	Public      bool                 `json:"public,omitempty"`
	Description string               `json:"description,omitempty"`
	Files       map[string]*gistFile `json:"files"`
}

// Auth token response.
type authResponse struct {
	URL    string   `json:"url,omitempty"`
	Scopes []string `json:"scopes"`
	Token  string   `json:"token,omitempty"`
	Note   string   `json:"note"`
}

//-------------------------------------------------------------------
//  HTTP requests
//-------------------------------------------------------------------

// Basic HTTP request.
func makeRequest(url string, method string, content []byte, contentType string, auth ...string) (body []byte, err error) {
	// Decide the type of the auth.
	var token, user, pass string
	if len(auth) == 1 {
		token = auth[0]
	} else if len(auth) == 2 {
		user = auth[0]
		pass = auth[1]
	}

	// Create the client.
	client := &http.Client{}

	// Prepare the request.
	req, err := http.NewRequest(method, url, bytes.NewReader(content))
	if err != nil {
		return
	}
	if token != "" {
		req.Header.Add("Authorization", "token "+token)
	}
	if user != "" && pass != "" {
		req.SetBasicAuth(user, pass)
	}
	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	// Fire the request.
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	// Check the status code.
	successStatusCode := http.StatusOK
	switch method {
	case http.MethodPost:
		successStatusCode = http.StatusCreated
	case http.MethodDelete:
		successStatusCode = http.StatusNoContent

	}
	if resp.StatusCode != successStatusCode {
		return nil, fmt.Errorf("invalid status code: %v (%v). should be: %v", resp.StatusCode, http.StatusText(resp.StatusCode), successStatusCode)
	}

	// Read the body.
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	resp.Body.Close()

	return
}

// HTTP get request.
func makeGet(url string, token string) (body []byte, err error) {
	body, err = makeRequest(url, http.MethodGet, nil, "", token)
	return
}

// HTTP post request.
func makePost(url string, token string, content []byte) (body []byte, err error) {
	body, err = makeRequest(url, http.MethodPost, content, "application/json", token)
	return
}

// HTTP patch request.
func makePatch(url string, token string, content []byte) (body []byte, err error) {
	body, err = makeRequest(url, http.MethodPatch, content, "application/json", token)
	return
}

// HTTP delete request.
func makeDelete(url string, token string) (err error) {
	_, err = makeRequest(url, http.MethodDelete, nil, "", token)
	return
}

//-------------------------------------------------------------------
//  Init functions
//-------------------------------------------------------------------

// Login to GitHub.
func login() (token string, err error) {
	// Username
	fmt.Print("Username: ")
	var username string
	_, err = fmt.Scanln(&username)
	if err != nil {
		return
	}
	// Password
	fmt.Print("Password: ")
	bytePass, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return
	}
	password := string(bytePass)

	// Newline after password entered.
	fmt.Println()

	token, err = createToken(username, password)

	return
}

// Checks the token validity.
func testToken(token string) (err error) {
	// Simple get request to the authenticated user.
	_, err = makeGet(userAPIURL, token)
	return
}

// List the user's authentications.
func listAuths(auth ...string) (auths []*authResponse, err error) {
	// Get the tokens.
	body, err := makeRequest(authAPIURL, http.MethodGet, nil, "", auth...)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &auths)

	return
}

// Get the app's authentication.
func getAuth(note string, creds ...string) (envAuth *authResponse, err error) {
	// Check if we already have an authentication for this machine.
	auths, err := listAuths(creds...)
	if err != nil {
		return
	}

	for _, auth := range auths {
		if auth.Note == note {
			envAuth = auth
			break
		}
	}

	return
}

// Basic auth request.
func createToken(user string, pass string) (token string, err error) {
	// Create the note for the token description.
	hostname, err := os.Hostname()
	if err != nil {
		return
	}
	note := authTokenNote + hostname

	// If already have an auth but don't have the token, delete it then create a new one.
	// Why? Because once a token generated cannot get the secret from it again.
	auth, err := getAuth(note, user, pass)
	if err != nil {
		return
	}
	if auth != nil {
		err = deleteAuth(auth.URL, user, pass)
		if err != nil {
			return
		}
	}

	// Create the payload.
	data := authResponse{
		Scopes: []string{"gist"},
		Note:   note,
	}
	payload, _ := json.Marshal(data)

	// Make the request.
	body, err := makeRequest(authAPIURL, http.MethodPost, payload, "application/json", user, pass)
	if err != nil {
		return
	}

	// Parse the response.
	authResp := authResponse{}
	err = json.Unmarshal(body, &authResp)
	if err != nil {
		return
	}

	token = authResp.Token

	return
}

// Deletes the authentication on the given URL.
func deleteAuth(URL string, auth ...string) (err error) {
	_, err = makeRequest(URL, http.MethodDelete, nil, "", auth...)
	return
}

//-------------------------------------------------------------------
//  GitHub Gist functions
//-------------------------------------------------------------------

// Gets the all of the user's gists from GitHub.
func getGists(token string) (userGists *[]gist, err error) {
	// Get all the gists.
	body, err := makeGet(gistAPIURL, token)
	if err != nil {
		return
	}

	// Parse the body json to structs.
	err = json.Unmarshal(body, &userGists)
	if err != nil {
		return
	}

	return
}

// Gets the envman gist.
func getOrCreateGist(token string) (g *gist, err error) {
	// Get all the gists.
	userGists, err := getGists(token)
	if err != nil {
		return
	}

	// Iterate while we didn't find the envman gist.
	for _, userGist := range *userGists {
		if userGist.Description == gistDescription {
			g = &userGist
			return
		}
	}

	// If the gist not found create it now.
	g, err = createGist(token)

	return
}

// Gets the content of a gist file.
func getGistFileContent(url string, token string) (content map[string]string, err error) {
	// Get the body in bytes.
	body, err := makeGet(url, token)
	if err != nil {
		return
	}

	// Decode the JSON.
	json.Unmarshal(body, &content)

	return
}

// Gets all the environment in a gist.
func getEnvironments(g *gist) (envs []string) {
	for env := range g.Files {
		if env == reservedName {
			continue
		}
		envs = append(envs, env)
	}

	return
}

// Gets the variables in an environment.
func getVariables(g *gist, envName string, token string) (vars []string, err error) {
	// First check if the environment exists.
	env, exists := g.Files[envName]
	if !exists {
		err = fmt.Errorf("environment \"%v\" doesn't exist", envName)
		return
	}

	// Get the gist file.
	envVars, err := getGistFileContent(env.URL, token)
	if err != nil {
		return
	}

	for envVar := range envVars {
		vars = append(vars, envVar)
	}

	return
}

// Creates the default gist.
func createGist(token string) (createdGist *gist, err error) {
	// Create the default gist file content.
	content := map[string]string{"created": time.Now().String()}
	contentJSON, err := json.Marshal(content)
	if err != nil {
		return
	}

	// Create the default file.
	baseFile := gistFile{
		Content: string(contentJSON),
	}

	// Create the default gist.
	baseGist := gist{
		Description: gistDescription,
		Public:      false,
		Files:       map[string]*gistFile{reservedName: &baseFile},
	}
	baseJSON, err := json.Marshal(&baseGist)
	if err != nil {
		return
	}

	body, err := makePost(gistAPIURL, token, baseJSON)
	if err != nil {
		return
	}

	// Parse the body json to structs.
	err = json.Unmarshal(body, &createdGist)
	if err != nil {
		return
	}

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

//-------------------------------------------------------------------
//  Interface functions
//-------------------------------------------------------------------

// Init makes the authentication if necessary.
func (g *GitHubGist) Init(c *config.Config) (err error) {
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

	return
}

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

// Update variables in the environment.
func (g *GitHubGist) Update(envName string, variables map[string]string) (err error) {
	err = updateGist(g.Token, envName, variables, g.EnvManGist)

	return
}

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
