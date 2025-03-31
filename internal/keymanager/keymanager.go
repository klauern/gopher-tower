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

// RotationPolicy defines when and how keys should be rotated
type RotationPolicy struct {
	// Interval between automatic key rotations
	Interval time.Duration `json:"interval,omitempty" yaml:"interval,omitempty"`
	// MaxAge is the maximum age of a key before it must be rotated
	MaxAge time.Duration `json:"max_age,omitempty" yaml:"max_age,omitempty"`
	// MaxVersions is the maximum number of key versions to keep
	MaxVersions int `json:"max_versions,omitempty" yaml:"max_versions,omitempty"`
}

// KeyMetadata contains information about the current key state
type KeyMetadata struct {
	// CurrentVersion is the active key version
	CurrentVersion int `json:"current_version"`
	// CreatedAt is when this key version was created
	CreatedAt time.Time `json:"created_at"`
	// LastRotated is when the key was last rotated
	LastRotated time.Time `json:"last_rotated"`
	// NextRotation is when the key is scheduled to be rotated next
	NextRotation time.Time `json:"next_rotation"`
}

// Config holds the configuration for secret management
type Config struct {
	// StoragePath for the encrypted keystore
	StoragePath string `json:"storage_path,omitempty" yaml:"storage_path,omitempty"`

	// MasterPassword used to secure the keystore
	MasterPassword string `json:"master_password,omitempty" yaml:"master_password,omitempty"`

	// RotationPolicy defines the key rotation settings
	RotationPolicy *RotationPolicy `json:"rotation_policy,omitempty" yaml:"rotation_policy,omitempty"`
}

// SecretManager defines the core interface for secret management operations
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

// KeyRotator defines the interface for key rotation operations
type KeyRotator interface {
	// RotateKeys performs a manual key rotation
	RotateKeys(ctx context.Context) error

	// GetKeyMetadata returns metadata about the current key state
	GetKeyMetadata(ctx context.Context) (KeyMetadata, error)
}

// RotatableSecretManager combines both secret management and key rotation capabilities
type RotatableSecretManager interface {
	SecretManager
	KeyRotator
}

// Factory function type for creating secret managers
type SecretManagerFactory func(Config) (SecretManager, error)

// Registry of available implementations
var implementations = make(map[string]SecretManagerFactory)
