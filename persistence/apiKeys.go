package persistence

import "errors"

// adds a new API key to the store
func AddAPIKey(key string) error {
	API_KEYS.Lock()
	defer API_KEYS.Unlock()
	for _, k := range API_KEYS.keys {
		if k == key {
			return errors.New("API key already exists")
		}
	}
	API_KEYS.keys = append(API_KEYS.keys, key)
	return nil
}

// removes an API key from the store
func DeleteAPIKey(key string) error {
	API_KEYS.Lock()
	defer API_KEYS.Unlock()
	for i, k := range API_KEYS.keys {
		if k == key {
			API_KEYS.keys = append(API_KEYS.keys[:i], API_KEYS.keys[i+1:]...)
			return nil
		}
	}
	return errors.New("API key not found")
}

// checks if an API key is valid
func IsAPIKeyValid(key string) bool {
	API_KEYS.RLock()
	defer API_KEYS.RUnlock()
	for _, k := range API_KEYS.keys {
		if k == key {
			return true
		}
	}
	return false
}
