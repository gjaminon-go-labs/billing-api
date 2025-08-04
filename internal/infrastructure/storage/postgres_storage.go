// PostgreSQL Storage Implementation
//
// This file implements the Storage interface using PostgreSQL with GORM.
// Provides: Production-ready storage backend, connection management, error handling
// Pattern: Key-value storage abstraction over relational database
// Used by: Production environments, integration tests with real databases
package storage

import (
	"encoding/json"
	"fmt"
	
	"gorm.io/gorm"
)

// PostgreSQLStorage provides a PostgreSQL implementation of the Storage interface
type PostgreSQLStorage struct {
	db *gorm.DB
}

// StorageRecord represents a key-value record in the storage table
type StorageRecord struct {
	Key   string `gorm:"primaryKey;size:255" json:"key"`
	Value string `gorm:"type:text" json:"value"`
}

// TableName specifies the table name for GORM
func (StorageRecord) TableName() string {
	return "storage_records"
}

// NewPostgreSQLStorage creates a new PostgreSQL storage instance
func NewPostgreSQLStorage(db *gorm.DB) *PostgreSQLStorage {
	storage := &PostgreSQLStorage{
		db: db,
	}
	
	// Note: Table creation is handled by the migration system using the migration user
	// The application user only has DML permissions for security
	
	return storage
}


// Store saves a value with the given key
func (s *PostgreSQLStorage) Store(key string, value interface{}) error {
	// Serialize value to JSON
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to serialize value for key %s: %w", key, err)
	}
	
	// Create or update record
	record := StorageRecord{
		Key:   key,
		Value: string(valueBytes),
	}
	
	// Use GORM's Save method which handles both create and update
	if err := s.db.Save(&record).Error; err != nil {
		return fmt.Errorf("failed to store value for key %s: %w", key, err)
	}
	
	return nil
}

// Get retrieves a value by key
func (s *PostgreSQLStorage) Get(key string) (interface{}, error) {
	var record StorageRecord
	
	// Find record by key
	if err := s.db.Where("key = ?", key).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("key not found: %s", key)
		}
		return nil, fmt.Errorf("failed to retrieve value for key %s: %w", key, err)
	}
	
	// Deserialize JSON value
	var value interface{}
	if err := json.Unmarshal([]byte(record.Value), &value); err != nil {
		return nil, fmt.Errorf("failed to deserialize value for key %s: %w", key, err)
	}
	
	return value, nil
}

// Exists checks if a key exists in storage
func (s *PostgreSQLStorage) Exists(key string) bool {
	var count int64
	
	// Count records with the given key
	s.db.Model(&StorageRecord{}).Where("key = ?", key).Count(&count)
	
	return count > 0
}

// ListAll retrieves all stored values
func (s *PostgreSQLStorage) ListAll() ([]interface{}, error) {
	var records []StorageRecord
	
	// Find all records
	if err := s.db.Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve all records: %w", err)
	}
	
	// Deserialize all values
	values := make([]interface{}, 0, len(records))
	for _, record := range records {
		var value interface{}
		if err := json.Unmarshal([]byte(record.Value), &value); err != nil {
			return nil, fmt.Errorf("failed to deserialize value for key %s: %w", record.Key, err)
		}
		values = append(values, value)
	}
	
	return values, nil
}

// Delete removes a value by key
func (s *PostgreSQLStorage) Delete(key string) error {
	// Delete record by key
	result := s.db.Where("key = ?", key).Delete(&StorageRecord{})
	
	if result.Error != nil {
		return fmt.Errorf("failed to delete value for key %s: %w", key, result.Error)
	}
	
	// Check if any record was actually deleted
	if result.RowsAffected == 0 {
		return fmt.Errorf("key not found: %s", key)
	}
	
	return nil
}

// Health checks the health of the PostgreSQL connection
func (s *PostgreSQLStorage) Health() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying SQL DB: %w", err)
	}
	
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	
	return nil
}

// Close closes the database connection
func (s *PostgreSQLStorage) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying SQL DB: %w", err)
	}
	
	return sqlDB.Close()
}

// GetDB returns the underlying GORM DB instance for testing purposes
// This method is intended for use by test helpers and should not be used in production code
func (s *PostgreSQLStorage) GetDB() *gorm.DB {
	return s.db
}

// Stats returns storage statistics
func (s *PostgreSQLStorage) Stats() (map[string]interface{}, error) {
	var count int64
	if err := s.db.Model(&StorageRecord{}).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("failed to get record count: %w", err)
	}
	
	sqlDB, err := s.db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying SQL DB: %w", err)
	}
	
	dbStats := sqlDB.Stats()
	
	return map[string]interface{}{
		"total_records":      count,
		"open_connections":   dbStats.OpenConnections,
		"in_use":            dbStats.InUse,
		"idle":              dbStats.Idle,
		"wait_count":        dbStats.WaitCount,
		"wait_duration":     dbStats.WaitDuration.String(),
		"max_idle_closed":   dbStats.MaxIdleClosed,
		"max_lifetime_closed": dbStats.MaxLifetimeClosed,
	}, nil
}