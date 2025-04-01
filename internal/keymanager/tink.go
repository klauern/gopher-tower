package keymanager

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/keyset"
	"github.com/google/tink/go/tink"
)

type secretData struct {
	Value     []byte    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Track which key version was used to encrypt this secret
	KeyVersion int `json:"key_version"`
}

type keysetMetadata struct {
	CurrentVersion int       `json:"current_version"`
	CreatedAt      time.Time `json:"created_at"`
	LastRotated    time.Time `json:"last_rotated"`
	NextRotation   time.Time `json:"next_rotation"`
}

type TinkManager struct {
	config     Config
	storageDir string
	primitive  tink.AEAD
	keyHandle  *keyset.Handle
	metadata   keysetMetadata
	mu         sync.RWMutex
}

const (
	metadataFile = "keyset_metadata.json"
)

func init() {
	implementations["tink"] = NewTinkManager
}

func NewTinkManager(config Config) (SecretManager, error) {
	if config.StoragePath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		config.StoragePath = filepath.Join(homeDir, ".gopher-tower", "secrets")
	}

	return &TinkManager{
		config:     config,
		storageDir: config.StoragePath,
	}, nil
}

// checkDirPermissions verifies that a directory has secure permissions
func (t *TinkManager) checkDirPermissions(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("failed to stat directory: %w", err)
	}

	// Check if it's a directory
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", dir)
	}

	// Check permissions - only owner should have access
	if info.Mode().Perm() != 0o700 {
		return fmt.Errorf("unsafe directory permissions: %s has %o, want 0700", dir, info.Mode().Perm())
	}

	return nil
}

// checkFilePermissions verifies that a file has secure permissions
func (t *TinkManager) checkFilePermissions(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return err
		}
		return fmt.Errorf("failed to stat file: %w", err)
	}

	// Check if it's a regular file
	if info.Mode()&os.ModeType != 0 {
		return fmt.Errorf("path is not a regular file: %s", path)
	}

	// Check permissions - only owner should have read/write access
	if info.Mode().Perm() != 0o600 {
		return fmt.Errorf("unsafe file permissions: %s has %o, want 0600", path, info.Mode().Perm())
	}

	return nil
}

func (t *TinkManager) Initialize(ctx context.Context, config Config) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Create storage directory if it doesn't exist
	if err := os.MkdirAll(t.storageDir, 0o700); err != nil {
		return fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Verify directory permissions
	if err := t.checkDirPermissions(t.storageDir); err != nil {
		return fmt.Errorf("storage directory is insecure: %w", err)
	}

	// Initialize or load keyset metadata
	if err := t.loadOrInitMetadata(); err != nil {
		return fmt.Errorf("failed to initialize metadata: %w", err)
	}

	// Initialize Tink AEAD primitive
	kh, err := keyset.NewHandle(aead.AES256GCMKeyTemplate())
	if err != nil {
		return fmt.Errorf("failed to create keyset handle: %w", err)
	}

	primitive, err := aead.New(kh)
	if err != nil {
		return fmt.Errorf("failed to create AEAD primitive: %w", err)
	}

	t.primitive = primitive
	t.keyHandle = kh

	return t.saveMetadata()
}

func (t *TinkManager) loadOrInitMetadata() error {
	path := filepath.Join(t.storageDir, metadataFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Initialize new metadata
			now := time.Now()
			t.metadata = keysetMetadata{
				CurrentVersion: 1,
				CreatedAt:      now,
				LastRotated:    now,
			}
			if t.config.RotationPolicy != nil && t.config.RotationPolicy.Interval > 0 {
				t.metadata.NextRotation = now.Add(t.config.RotationPolicy.Interval)
			}
			return nil
		}
		return fmt.Errorf("failed to read metadata: %w", err)
	}

	if err := json.Unmarshal(data, &t.metadata); err != nil {
		return fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return nil
}

func (t *TinkManager) saveMetadata() error {
	data, err := json.Marshal(t.metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	path := filepath.Join(t.storageDir, metadataFile)
	return os.WriteFile(path, data, 0o600)
}

func (t *TinkManager) RotateKeys(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.primitive == nil || t.keyHandle == nil {
		return fmt.Errorf("secret manager not initialized")
	}

	// Create new key
	newHandle, err := keyset.NewHandle(aead.AES256GCMKeyTemplate())
	if err != nil {
		return fmt.Errorf("failed to create new key: %w", err)
	}

	newPrimitive, err := aead.New(newHandle)
	if err != nil {
		return fmt.Errorf("failed to create new primitive: %w", err)
	}

	// Re-encrypt all secrets with new key
	err = t.reencryptSecrets(ctx, newPrimitive)
	if err != nil {
		return fmt.Errorf("failed to re-encrypt secrets: %w", err)
	}

	// Update metadata
	now := time.Now()
	t.metadata.CurrentVersion++
	t.metadata.LastRotated = now
	if t.config.RotationPolicy != nil && t.config.RotationPolicy.Interval > 0 {
		t.metadata.NextRotation = now.Add(t.config.RotationPolicy.Interval)
	}

	// Save metadata before switching keys
	if err := t.saveMetadata(); err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	// Switch to new key
	t.primitive = newPrimitive
	t.keyHandle = newHandle

	return nil
}

func (t *TinkManager) reencryptSecrets(ctx context.Context, newPrimitive tink.AEAD) error {
	entries, err := os.ReadDir(t.storageDir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == metadataFile {
			continue
		}

		// Load secret
		key := entry.Name()
		if !strings.HasSuffix(key, ".secret") {
			continue
		}
		key = key[:len(key)-7] // remove .secret suffix

		data, err := t.loadSecretData(key)
		if err != nil {
			return fmt.Errorf("failed to load secret %s: %w", key, err)
		}

		// Decrypt with old key
		decrypted, err := t.primitive.Decrypt(data.Value, nil)
		if err != nil {
			return fmt.Errorf("failed to decrypt secret %s: %w", key, err)
		}

		// Encrypt with new key
		encrypted, err := newPrimitive.Encrypt(decrypted, nil)
		if err != nil {
			return fmt.Errorf("failed to encrypt secret %s: %w", key, err)
		}

		// Save with new encryption
		data.Value = encrypted
		data.UpdatedAt = time.Now()
		data.KeyVersion = t.metadata.CurrentVersion + 1

		if err := t.saveSecretData(key, data); err != nil {
			return fmt.Errorf("failed to save re-encrypted secret %s: %w", key, err)
		}
	}

	return nil
}

func (t *TinkManager) GetKeyMetadata(ctx context.Context) (KeyMetadata, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.primitive == nil {
		return KeyMetadata{}, fmt.Errorf("secret manager not initialized")
	}

	return KeyMetadata{
		CurrentVersion: t.metadata.CurrentVersion,
		CreatedAt:      t.metadata.CreatedAt,
		LastRotated:    t.metadata.LastRotated,
		NextRotation:   t.metadata.NextRotation,
	}, nil
}

func (t *TinkManager) GetSecret(ctx context.Context, key string) (string, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.primitive == nil {
		return "", fmt.Errorf("secret manager not initialized")
	}

	// Verify directory permissions before proceeding
	if err := t.checkDirPermissions(t.storageDir); err != nil {
		return "", fmt.Errorf("storage directory is insecure: %w", err)
	}

	path := filepath.Join(t.storageDir, key+".secret")
	if err := t.checkFilePermissions(path); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("secret not found: %s", key)
		}
		return "", fmt.Errorf("secret file is insecure: %w", err)
	}

	data, err := t.loadSecretData(key)
	if err != nil {
		return "", err
	}

	// Decrypt the value using Tink
	decrypted, err := t.primitive.Decrypt(data.Value, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt secret: %w", err)
	}

	return string(decrypted), nil
}

func (t *TinkManager) SetSecret(ctx context.Context, key, value string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.primitive == nil {
		return fmt.Errorf("secret manager not initialized")
	}

	// Verify directory permissions before proceeding
	if err := t.checkDirPermissions(t.storageDir); err != nil {
		return fmt.Errorf("storage directory is insecure: %w", err)
	}

	// Encrypt the value using Tink
	encrypted, err := t.primitive.Encrypt([]byte(value), nil)
	if err != nil {
		return fmt.Errorf("failed to encrypt secret: %w", err)
	}

	data := secretData{
		Value:      encrypted,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		KeyVersion: t.metadata.CurrentVersion,
	}

	return t.saveSecretData(key, data)
}

func (t *TinkManager) DeleteSecret(ctx context.Context, key string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	path := filepath.Join(t.storageDir, key+".secret")
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete secret: %w", err)
	}
	return nil
}

func (t *TinkManager) Status(ctx context.Context) (Status, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.primitive == nil {
		return Status{
			Healthy: false,
			Message: "secret manager not initialized",
		}, nil
	}

	return Status{
		Healthy: true,
		Message: "secret manager operational",
	}, nil
}

func (t *TinkManager) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.primitive = nil
	return nil
}

func (t *TinkManager) loadSecretData(key string) (secretData, error) {
	var data secretData

	path := filepath.Join(t.storageDir, key+".secret")
	file, err := os.ReadFile(path)
	if err != nil {
		return data, fmt.Errorf("failed to read secret file: %w", err)
	}

	if err := json.Unmarshal(file, &data); err != nil {
		return data, fmt.Errorf("failed to unmarshal secret data: %w", err)
	}

	return data, nil
}

func (t *TinkManager) saveSecretData(key string, data secretData) error {
	file, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal secret data: %w", err)
	}

	path := filepath.Join(t.storageDir, key+".secret")
	return os.WriteFile(path, file, 0o600)
}
