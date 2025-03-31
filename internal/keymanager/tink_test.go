package keymanager

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDir creates a temporary directory for testing.
func setupTestDir(t *testing.T) string {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "keymanager_test_*")
	require.NoError(t, err, "Failed to create temp dir")
	return tempDir
}

// cleanupTestDir removes the temporary directory.
func cleanupTestDir(t *testing.T, dir string) {
	t.Helper()
	err := os.RemoveAll(dir)
	assert.NoError(t, err, "Failed to remove temp dir")
}

func TestTinkManager(t *testing.T) {
	t.Run("Successful Initialization", func(t *testing.T) {
		tempDir := setupTestDir(t)
		defer cleanupTestDir(t, tempDir)

		config := Config{
			StoragePath: tempDir,
		}

		mgr, err := NewTinkManager(config)
		require.NoError(t, err, "NewTinkManager should succeed")
		require.NotNil(t, mgr, "Manager should not be nil")

		err = mgr.Initialize(context.Background(), config)
		require.NoError(t, err, "Initialize should succeed")

		status, err := mgr.Status(context.Background())
		require.NoError(t, err, "Status check should succeed")
		assert.True(t, status.Healthy, "Manager should be healthy after initialization")

		err = mgr.Close()
		assert.NoError(t, err, "Close should succeed")
	})

	t.Run("Secret Lifecycle", func(t *testing.T) {
		tempDir := setupTestDir(t)
		defer cleanupTestDir(t, tempDir)

		config := Config{
			StoragePath: tempDir,
		}

		mgr, err := NewTinkManager(config)
		require.NoError(t, err, "NewTinkManager should succeed")

		err = mgr.Initialize(context.Background(), config)
		require.NoError(t, err, "Initialize should succeed")

		// Test SetSecret
		err = mgr.SetSecret(context.Background(), "test-key", "test-value")
		require.NoError(t, err, "SetSecret should succeed")

		// Test GetSecret
		value, err := mgr.GetSecret(context.Background(), "test-key")
		require.NoError(t, err, "GetSecret should succeed")
		assert.Equal(t, "test-value", value, "Retrieved value should match set value")

		// Test DeleteSecret
		err = mgr.DeleteSecret(context.Background(), "test-key")
		require.NoError(t, err, "DeleteSecret should succeed")

		// Verify secret is deleted
		_, err = mgr.GetSecret(context.Background(), "test-key")
		assert.Error(t, err, "GetSecret should fail after deletion")

		err = mgr.Close()
		assert.NoError(t, err, "Close should succeed")
	})

	t.Run("Concurrent Operations", func(t *testing.T) {
		tempDir := setupTestDir(t)
		defer cleanupTestDir(t, tempDir)

		config := Config{
			StoragePath: tempDir,
		}

		mgr, err := NewTinkManager(config)
		require.NoError(t, err)
		err = mgr.Initialize(context.Background(), config)
		require.NoError(t, err)

		var wg sync.WaitGroup
		numGoroutines := 10
		numOperations := 50
		errChan := make(chan error, numGoroutines*numOperations)

		// Concurrent writes and reads
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					key := fmt.Sprintf("key-%d-%d", id, j)
					value := fmt.Sprintf("value-%d-%d", id, j)

					// Set secret
					if err := mgr.SetSecret(context.Background(), key, value); err != nil {
						errChan <- fmt.Errorf("set error: %w", err)
						continue
					}

					// Get secret
					got, err := mgr.GetSecret(context.Background(), key)
					if err != nil {
						errChan <- fmt.Errorf("get error: %w", err)
						continue
					}

					if got != value {
						errChan <- fmt.Errorf("value mismatch for key %s: got %s, want %s", key, got, value)
					}

					// Delete secret
					if err := mgr.DeleteSecret(context.Background(), key); err != nil {
						errChan <- fmt.Errorf("delete error: %w", err)
					}
				}
			}(i)
		}

		wg.Wait()
		close(errChan)

		var errors []error
		for err := range errChan {
			errors = append(errors, err)
		}
		assert.Empty(t, errors, "concurrent operations should not produce errors")

		err = mgr.Close()
		assert.NoError(t, err)
	})

	t.Run("Error Conditions", func(t *testing.T) {
		t.Run("Uninitialized Manager", func(t *testing.T) {
			tempDir := setupTestDir(t)
			defer cleanupTestDir(t, tempDir)

			config := Config{
				StoragePath: tempDir,
			}

			mgr, err := NewTinkManager(config)
			require.NoError(t, err)
			// Deliberately skip initialization

			// Operations should fail
			_, err = mgr.GetSecret(context.Background(), "key")
			assert.Error(t, err, "GetSecret should fail when uninitialized")

			err = mgr.SetSecret(context.Background(), "key", "value")
			assert.Error(t, err, "SetSecret should fail when uninitialized")

			status, err := mgr.Status(context.Background())
			require.NoError(t, err)
			assert.False(t, status.Healthy, "Status should report unhealthy when uninitialized")
		})

		t.Run("Invalid Storage Path", func(t *testing.T) {
			config := Config{
				StoragePath: "/nonexistent/path/that/should/not/exist",
			}

			mgr, err := NewTinkManager(config)
			require.NoError(t, err)

			err = mgr.Initialize(context.Background(), config)
			assert.Error(t, err, "Initialize should fail with invalid storage path")
		})

		t.Run("Get Non-existent Secret", func(t *testing.T) {
			tempDir := setupTestDir(t)
			defer cleanupTestDir(t, tempDir)

			config := Config{
				StoragePath: tempDir,
			}

			mgr, err := NewTinkManager(config)
			require.NoError(t, err)
			err = mgr.Initialize(context.Background(), config)
			require.NoError(t, err)

			_, err = mgr.GetSecret(context.Background(), "nonexistent-key")
			assert.Error(t, err, "GetSecret should fail for non-existent key")
		})

		t.Run("Delete Non-existent Secret", func(t *testing.T) {
			tempDir := setupTestDir(t)
			defer cleanupTestDir(t, tempDir)

			config := Config{
				StoragePath: tempDir,
			}

			mgr, err := NewTinkManager(config)
			require.NoError(t, err)
			err = mgr.Initialize(context.Background(), config)
			require.NoError(t, err)

			err = mgr.DeleteSecret(context.Background(), "nonexistent-key")
			assert.NoError(t, err, "DeleteSecret should not error on non-existent key")
		})
	})

	t.Run("File Permissions", func(t *testing.T) {
		t.Run("Storage Directory Permissions", func(t *testing.T) {
			tempDir := setupTestDir(t)
			defer cleanupTestDir(t, tempDir)

			config := Config{
				StoragePath: tempDir,
			}

			mgr, err := NewTinkManager(config)
			require.NoError(t, err)
			err = mgr.Initialize(context.Background(), config)
			require.NoError(t, err)

			// Check storage directory permissions
			info, err := os.Stat(tempDir)
			require.NoError(t, err)
			assert.Equal(t, os.FileMode(0o700), info.Mode().Perm(),
				"Storage directory should have 0700 permissions")

			// Set and verify secret file permissions
			testKey := "permission-test-key"
			err = mgr.SetSecret(context.Background(), testKey, "test-value")
			require.NoError(t, err)

			secretPath := filepath.Join(tempDir, testKey+".secret")
			info, err = os.Stat(secretPath)
			require.NoError(t, err)
			assert.Equal(t, os.FileMode(0o600), info.Mode().Perm(),
				"Secret file should have 0600 permissions")
		})

		t.Run("Permission Changes", func(t *testing.T) {
			tempDir := setupTestDir(t)
			defer cleanupTestDir(t, tempDir)

			config := Config{
				StoragePath: tempDir,
			}

			// Create manager and set initial secret
			mgr, err := NewTinkManager(config)
			require.NoError(t, err)
			err = mgr.Initialize(context.Background(), config)
			require.NoError(t, err)

			testKey := "permission-change-key"
			err = mgr.SetSecret(context.Background(), testKey, "test-value")
			require.NoError(t, err)

			secretPath := filepath.Join(tempDir, testKey+".secret")

			// Attempt to modify permissions
			err = os.Chmod(secretPath, 0o644)
			require.NoError(t, err)

			// Verify that reading the secret fails with wrong permissions
			_, err = mgr.GetSecret(context.Background(), testKey)
			assert.Error(t, err, "GetSecret should fail with incorrect file permissions")

			// Fix permissions and verify it works again
			err = os.Chmod(secretPath, 0o600)
			require.NoError(t, err)

			_, err = mgr.GetSecret(context.Background(), testKey)
			assert.NoError(t, err, "GetSecret should succeed after fixing permissions")
		})

		t.Run("Parent Directory Permissions", func(t *testing.T) {
			tempDir := setupTestDir(t)
			defer cleanupTestDir(t, tempDir)

			// Create a subdirectory with wrong permissions
			unsafeDir := filepath.Join(tempDir, "unsafe")
			err := os.MkdirAll(unsafeDir, 0o777)
			require.NoError(t, err)

			config := Config{
				StoragePath: unsafeDir,
			}

			mgr, err := NewTinkManager(config)
			require.NoError(t, err)

			// Initialize should fail due to unsafe permissions
			err = mgr.Initialize(context.Background(), config)
			assert.Error(t, err, "Initialize should fail with unsafe directory permissions")

			// Fix permissions and retry
			err = os.Chmod(unsafeDir, 0o700)
			require.NoError(t, err)

			err = mgr.Initialize(context.Background(), config)
			assert.NoError(t, err, "Initialize should succeed after fixing permissions")
		})

		t.Run("Symlink Protection", func(t *testing.T) {
			tempDir := setupTestDir(t)
			defer cleanupTestDir(t, tempDir)

			// Create a directory for actual storage
			actualDir := filepath.Join(tempDir, "actual")
			err := os.MkdirAll(actualDir, 0o700)
			require.NoError(t, err)

			// Create a symlink
			symlink := filepath.Join(tempDir, "symlink")
			err = os.Symlink(actualDir, symlink)
			require.NoError(t, err)

			config := Config{
				StoragePath: symlink,
			}

			mgr, err := NewTinkManager(config)
			require.NoError(t, err)

			// Initialize through symlink
			err = mgr.Initialize(context.Background(), config)
			require.NoError(t, err)

			// Verify operations work through symlink
			err = mgr.SetSecret(context.Background(), "symlink-test", "test-value")
			assert.NoError(t, err, "SetSecret should work through symlink")

			// Verify file exists in actual directory
			_, err = os.Stat(filepath.Join(actualDir, "symlink-test.secret"))
			assert.NoError(t, err, "Secret file should exist in actual directory")
		})
	})
}
