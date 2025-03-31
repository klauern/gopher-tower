package keymanager

import (
	"context"
	"time"
)

// Status represents the current state of the secret management system
type Status struct {
	Healthy    bool
	Message    string
	LastRotate time.Time
}

// Config holds the configuration for secret management
type Config struct {
	// StoragePath for the encrypted keystore
	StoragePath string `json:"storage_path,omitempty" yaml:"storage_path,omitempty"`

	// MasterPassword used to secure the keystore
	MasterPassword string `json:"master_password,omitempty" yaml:"master_password,omitempty"`
}

// SecretManager defines the interface for secret management operations
type SecretManager interface {
	// Initialize sets up the secret management system
	Initialize(ctx context.Context, config Config) error

	// GetSecret retrieves a secret by key
	GetSecret(ctx context.Context, key string) (string, error)

	// SetSecret stores a secret
	SetSecret(ctx context.Context, key, value string) error

	// DeleteSecret removes a secret
	DeleteSecret(ctx context.Context, key string) error

	// Status returns the current status
	Status(ctx context.Context) (Status, error)

	// Close cleans up any resources
	Close() error
}

// Factory function type for creating secret managers
type SecretManagerFactory func(Config) (SecretManager, error)

// Registry of available implementations
var implementations = make(map[string]SecretManagerFactory)
