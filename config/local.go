package config

// LocalConfig structure.
type LocalConfig struct {
	Environments map[string]map[string]string `json:"environments"`
}
