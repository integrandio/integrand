package persistence

import "sync"

// Data structure to store API keys
var apiKeys struct {
	sync.RWMutex
	keys map[string]bool
}

func init() {
	apiKeys.keys = make(map[string]bool)
}

// AddAPIKey adds a new API key to the store
func AddAPIKey(key string) {
	apiKeys.Lock()
	defer apiKeys.Unlock()
	apiKeys.keys[key] = true
}

// DeleteAPIKey removes an API key from the store
func DeleteAPIKey(key string) {
	apiKeys.Lock()
	defer apiKeys.Unlock()
	delete(apiKeys.keys, key)
}

// IsAPIKeyValid checks if an API key is valid
func IsAPIKeyValid(key string) bool {
	apiKeys.RLock()
	defer apiKeys.RUnlock()
	return apiKeys.keys[key]
}
