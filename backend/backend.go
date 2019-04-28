package backend

import (
	"github.com/pyrooka/envman/config"
)

// IBackend is an interface what show what should implement if you want to create a new storage service.
type IBackend interface {
	Init(c *config.Config) (err error)                         // Initialize the backend.
	List(envName string) ([]string, error)                     // Returns a list with the variables in the env or the environments in the backend if the name is empty string.
	Get(envName string) (vars map[string]string, err error)    // Gets the variables for the environment.
	Update(envName string, vars map[string]string) (err error) // Updates variables in the environment.
	Delete(envName string, vars []string) (err error)          // Deletes the given variables or the full environment if an empty slice given.
	CleanUp() (err error)                                      // Removes all the created things.
}
