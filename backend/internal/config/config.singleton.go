package config

import "sync"

var (
	instance *Config
	once     sync.Once
)

// Get returns the singleton config instance.
// Panics on first call if config cannot be loaded.
func Get() *Config {
	once.Do(func() {
		cfg, err := Load()
		if err != nil {
			panic("failed to load config: " + err.Error())
		}
		instance = cfg
	})
	return instance
}
