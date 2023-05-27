package configparser

// Provider is config provider
type Provider interface {
	// Get returns the Parser if succeed or error otherwise.
	Get() (*Parser, error)
}

// Default is a default config provider
func Default() Provider {
	return NewFile()
}
