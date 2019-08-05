package storage

// IStorage is an interface what show what should implement if you want to create a new storage service.
type IStorage interface {
	Init() (err error)                                         // Initialize the storage.
	ListEnvs() (envs []string, err error)                      // Returns a list with the environments in the storage.
	ListVars(envName string) (vars []string, err error)        // Returns a list with the variables in the environment.
	Get(envName string) (vars map[string]string, err error)    // Gets the variables from the environment.
	Update(envName string, vars map[string]string) (err error) // Updates variables in the environment.
	DeleteEnv(envName string, vars []string) (err error)       // Deletes the given variables or the full environment if an empty slice given.
	DeleteVar(envName string, vars []string) (err error)       // Deletes the given variables or the full environment if an empty slice given.
	CleanUp() (err error)                                      // Removes all the created things.
}
