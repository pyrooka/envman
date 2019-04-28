# envman
Envman (environment variable manager) is a simple tool which stores your environment variables in the defined backend (e.g. GitHub Gist) so you can easily load them when you want to continue your work.

## Usage
```
envman [global options] command [command options] [arguments...]

COMMANDS:
     list, ls    List the environments or variables in the environment
     load, l     Load and environment to the current one
     save, s     Save environment variables to an environment
     remove, rm  Remove a full environment or just a variable
     cleanup     Cleanup the backend, delete all the created files
     help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --backend value, -b value  Use and set a different backend as default
   --help, -h                 show help
   --version, -v              print the version
```

## Sample commands
`envman ls`  
`envman ls ENV_NAME`  
`envman save ENV_NAME VAR_1 VAR_2`  
`envman load ENV_NAME`  
`envman rm ENV_NAME`  
`envman rm ENV_NAME VAR_1`

## Backend development
- Implement the Backend interface.
- If want to use config for your backend add it to the Config struct.
- In the main package declare it in backend parsing.

## TODO
- Autocomplete
- Security: E.g. AES encrypt any text which is uploaded.
- Display better error messages?
- Option to rename env?
### Backends
#### GitHub Gist
- Refactor (GraphQL?)
