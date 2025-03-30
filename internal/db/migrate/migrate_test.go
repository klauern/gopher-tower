package migrate

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrateDB(t *testing.T) {
	// Create a temporary test database
	tmpDir, err := os.MkdirTemp("", "gopher-tower-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")

	// Create a mock filesystem with test migrations
	mockFS := fstest.MapFS{
		"000001_init_schema.up.sql": &fstest.MapFile{
			Data: []byte(`CREATE TABLE test_table (id INTEGER PRIMARY KEY);`),
		},
		"000001_init_schema.down.sql": &fstest.MapFile{
			Data: []byte(`DROP TABLE IF EXISTS test_table;`),
		},
	}

	tests := []struct {
		name    string
		dbPath  string
		fs      FS
		wantErr bool
	}{
		{
			name:    "successful migration",
			dbPath:  dbPath,
			fs:      mockFS,
			wantErr: false,
		},
		{
			name:    "invalid db path",
			dbPath:  "/nonexistent/path/db.sqlite",
			fs:      mockFS,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use a unique database file for each test case
			testDBPath := tt.dbPath
			if !tt.wantErr {
				testDBPath = filepath.Join(tmpDir, tt.name+".db")
			}

			err := MigrateDB(testDBPath, tt.fs)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Verify migration version
			version, dirty, err := GetMigrationVersion(testDBPath, tt.fs)
			require.NoError(t, err)
			assert.Equal(t, uint(1), version)
			assert.False(t, dirty)
		})
	}
}

func TestRollbackDB(t *testing.T) {
	// Create a temporary test database
	tmpDir, err := os.MkdirTemp("", "gopher-tower-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")

	// Create a mock filesystem with test migrations
	mockFS := fstest.MapFS{
		"000001_init_schema.up.sql": &fstest.MapFile{
			Data: []byte(`CREATE TABLE test_table (id INTEGER PRIMARY KEY);`),
		},
		"000001_init_schema.down.sql": &fstest.MapFile{
			Data: []byte(`DROP TABLE IF EXISTS test_table;`),
		},
	}

	tests := []struct {
		name    string
		dbPath  string
		fs      FS
		wantErr bool
	}{
		{
			name:    "successful rollback",
			dbPath:  dbPath,
			fs:      mockFS,
			wantErr: false,
		},
		{
			name:    "invalid db path",
			dbPath:  "/nonexistent/path/db.sqlite",
			fs:      mockFS,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use a unique database file for each test case
			testDBPath := tt.dbPath
			if !tt.wantErr {
				testDBPath = filepath.Join(tmpDir, tt.name+".db")
				// First apply the migration
				err = MigrateDB(testDBPath, tt.fs)
				require.NoError(t, err)
			}

			err := RollbackDB(testDBPath, tt.fs)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Verify migration version is back to 0
			version, dirty, err := GetMigrationVersion(testDBPath, tt.fs)
			if !tt.wantErr {
				require.NoError(t, err)
				assert.Equal(t, uint(0), version)
				assert.False(t, dirty)
			}
		})
	}
}

func TestGetMigrationVersion(t *testing.T) {
	// Create a temporary test database
	tmpDir, err := os.MkdirTemp("", "gopher-tower-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a mock filesystem with test migrations
	mockFS := fstest.MapFS{
		"000001_init_schema.up.sql": &fstest.MapFile{
			Data: []byte(`CREATE TABLE test_table (id INTEGER PRIMARY KEY);`),
		},
		"000001_init_schema.down.sql": &fstest.MapFile{
			Data: []byte(`DROP TABLE IF EXISTS test_table;`),
		},
	}

	tests := []struct {
		name          string
		dbPath        string
		fs            FS
		shouldMigrate bool
		wantVersion   uint
		wantDirty     bool
		wantErr       bool
	}{
		{
			name:          "get version after migration",
			dbPath:        filepath.Join(tmpDir, "migrated.db"),
			fs:            mockFS,
			shouldMigrate: true,
			wantVersion:   1,
			wantDirty:     false,
			wantErr:       false,
		},
		{
			name:          "get version without migration",
			dbPath:        filepath.Join(tmpDir, "fresh.db"),
			fs:            mockFS,
			shouldMigrate: false,
			wantVersion:   0,
			wantDirty:     false,
			wantErr:       false,
		},
		{
			name:          "invalid db path",
			dbPath:        "/nonexistent/path/db.sqlite",
			fs:            mockFS,
			shouldMigrate: false,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldMigrate {
				err := MigrateDB(tt.dbPath, tt.fs)
				require.NoError(t, err)
			}

			version, dirty, err := GetMigrationVersion(tt.dbPath, tt.fs)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantVersion, version)
			assert.Equal(t, tt.wantDirty, dirty)
		})
	}
}
