package local

// Local uses the computer for backend storage.
type Local struct {
	Environments map[string]map[string]string `json:"environments"`
}
