package storage

import (
	"errors"

	"github.com/pyrooka/envman/storage/local"
)

// IStorage is an interface what show what should implement if you want to create a new storage service.
type IStorage interface {
	Init() (err error)                                         // Initialize the storage.
	List(envName string) (envs []string, err error)            // Returns a list of variables in the environment or the environments in the storage if the envName is an empty string.
	Get(envName string) (vars map[string]string, err error)    // Gets the variables from the environment.
	Update(envName string, vars map[string]string) (err error) // Updates variables in the environment.
	Delete(envName string, vars []string) (err error)          // Deletes the given variables or the full environment if an empty slice given.
	CleanUp() (err error)                                      // Removes all the created things.
}

// SelectStorage returns a storage object based on the given string.
func SelectStorage(storageName string) (IStorage, error) {
	switch storageName {
	case "local":
		l := local.Local{}
		return &l, nil
	case "githubgist":
		g := local.Local{}
		return &g, nil
	default:
		return nil, errors.New("storage not found")
	}
}
