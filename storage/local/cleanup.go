package local

// CleanUp removes all the environments from the config.
func (l *Local) CleanUp() (err error) {
	for key := range l.Environments {
		delete(l.Environments, key)
	}
	return
}
