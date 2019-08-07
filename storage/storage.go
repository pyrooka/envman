package storage

import (
	"errors"

	"github.com/pyrooka/envman/storage/local"
)

var currentStorage IStorage

// IStorage is an interface what show what should implement if you want to create a new storage service.
type IStorage interface {
	Init() (err error)                                         // Initialize the storage.
	List(envName string) (envs []string, err error)            // Returns a list of variables in the environment or the environments in the storage if the envName is an empty string.
	Get(envName string) (vars map[string]string, err error)    // Gets the variables from the environment.
	Update(envName string, vars map[string]string) (err error) // Updates variables in the environment.
	Delete(envName string, vars []string) (err error)          // Deletes the given variables or the full environment if an empty slice given.
	CleanUp() (err error)                                      // Removes all the created things.
}

// SetStorage sets the storage to the global variable base on the given name.
func SetStorage(name string) error {
	switch name {
	case "local":
		currentStorage = &local.Local{}
	case "githubgist":
		currentStorage = &local.Local{}
	default:
		return errors.New("storage not found")
	}

	return nil
}

// GetStorage returns the storage have been set for this session.
func GetStorage() (IStorage, error) {
	if currentStorage == nil {
		return nil, errors.New("storage is unset")
	}

	return currentStorage, nil
}
