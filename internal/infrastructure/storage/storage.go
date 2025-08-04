package storage

// Storage defines the contract for data storage backends
type Storage interface {
	// Store saves a value with the given key
	Store(key string, value interface{}) error
	
	// Get retrieves a value by key
	Get(key string) (interface{}, error)
	
	// Exists checks if a key exists in storage
	Exists(key string) bool
	
	// ListAll retrieves all stored values
	ListAll() ([]interface{}, error)
	
	// Delete removes a value by key
	Delete(key string) error
}