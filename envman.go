package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"gopkg.in/urfave/cli.v1"

	"github.com/pyrooka/envman/backend"
	"github.com/pyrooka/envman/config"
)

const scriptPrefixTemplate = "loadenv_%s.%s"

// Creates a shell script.
func createScript(name string, content string) (err error) {
	file, err := os.Create(name)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.WriteString(content)

	return
}

// Creates the shell scripts with the environemnt variables based on the OS type.
func createScripts(envName string, envVars map[string]string) (err error) {
	var scriptName, script string

	switch runtime.GOOS {
	case "windows":
		// First create the batch script.
		scriptName = fmt.Sprintf(scriptPrefixTemplate, envName, "bat")
		script = fmt.Sprintln("@echo off")
		for key, value := range envVars {
			script += fmt.Sprintln(fmt.Sprintf("SET %s=%s", key, value))
		}
		err = createScript(scriptName, script)
		if err != nil {
			return
		}

		// Then the powershell script.
		scriptName = fmt.Sprintf(scriptPrefixTemplate, envName, "ps1")
		script = ""
		for key, value := range envVars {
			script += fmt.Sprintln(fmt.Sprintf("$env:%s=\"%s\"", key, value))
		}
		err = createScript(scriptName, script)

	default:
		// If not Windows should be SH compatible, right?
		scriptName = fmt.Sprintf(scriptPrefixTemplate, envName, "sh")
		for key, value := range envVars {
			script += fmt.Sprintln(fmt.Sprintf("export %s=%s", key, value))
		}

		err = createScript(scriptName, script)
	}

	return
}

// Gets the backend based on the input string.
func getBackend(backendStr string) (backend.IBackend, error) {
	switch backendStr {
	case "local":
		l := backend.Local{}
		return &l, nil
	case "githubgist":
		g := backend.GitHubGist{}
		return &g, nil
	default:
		return nil, errors.New("backend not found")
	}
}

func main() {
	// Load the config.
	conf, err := config.Load()
	if err != nil {
		fmt.Println("Error while reading the config: " + err.Error())
		return
	}

	// The backend which we will use.
	var backendObj backend.IBackend

	app := cli.NewApp()
	app.Name = "Envman"
	app.Usage = "Manage your environment variables"
	app.Version = "1.0.0"

	// Get the backend and initialize it before execute the command action.
	app.Before = func(c *cli.Context) error {
		backendStr := c.String("backend")
		if backendStr != "" {
			conf.DefaultBackend = backendStr
		} else {
			backendStr = conf.DefaultBackend
		}

		backendObj, err = getBackend(backendStr)
		if backendObj != nil {
			err = backendObj.Init(conf)
		}

		return err
	}

	// Command line flags.
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "backend, b",
			Usage: "Use and set a different backend as default",
		},
	}

	// Command line commands.
	app.Commands = []cli.Command{
		{
			Name:      "list",
			Aliases:   []string{"ls"},
			Usage:     "List the environments or variables in the environment",
			ArgsUsage: "[environment name]",
			Action: func(c *cli.Context) error {
				// If the first arg is not provided (empty string ""), then list the environments.
				result, err := backendObj.List(c.Args().First())
				if err == nil {
					fmt.Println(result)
				}
				return err
			},
		},
		{
			Name:      "load",
			Aliases:   []string{"l"},
			Usage:     "Load and environment to the current one",
			ArgsUsage: "environment_name",
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					return errors.New("missing environment name")
				}

				envName := c.Args().First()
				vars, err := backendObj.Get(envName)
				if err != nil {
					return err
				}

				err = createScripts(envName, vars)
				return err
			},
		},
		{
			Name:      "save",
			Aliases:   []string{"s"},
			Usage:     "Save environment variables to an environment",
			ArgsUsage: "environment_name environment_variables...",
			Action: func(c *cli.Context) error {
				if c.NArg() < 2 {
					return errors.New("not enough argument")
				}

				envVars := make(map[string]string)
				args := c.Args()
				for _, key := range args[1:] {
					if value, exists := os.LookupEnv(key); exists {
						envVars[key] = value
					} else {
						fmt.Println(fmt.Sprintf("Environment variable %v skipped, because doesn't exist.", key))
					}
				}

				err = backendObj.Update(args[0], envVars)
				return err
			},
		},
		{
			Name:      "remove",
			Aliases:   []string{"rm"},
			Usage:     "Remove a full environment or just a variable",
			ArgsUsage: "environment_name [environment_variables...]",
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					return errors.New("not enough argument")
				}

				args := c.Args()
				if len(args) == 1 {
					err = backendObj.Delete(args[0], []string{})
				} else {
					err = backendObj.Delete(args[0], args[1:])
				}

				return err
			},
		},
		{
			Name:  "cleanup",
			Usage: "Cleanup the backend, delete all the created files",
			Action: func(c *cli.Context) error {
				err = backendObj.CleanUp()
				return err
			},
		},
	}

	// Run the command line application.
	err = app.Run(os.Args)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}

	// Save the config.
	err = conf.Save()
	if err != nil {
		fmt.Println("Error while writing the config: " + err.Error())
		return
	}
}
