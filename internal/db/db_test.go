package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockDB implements DBTX interface for testing
type mockDB struct {
	execContext     func(context.Context, string, ...interface{}) (sql.Result, error)
	prepareContext  func(context.Context, string) (*sql.Stmt, error)
	queryContext    func(context.Context, string, ...interface{}) (*sql.Rows, error)
	queryRowContext func(context.Context, string, ...interface{}) *sql.Row
}

func (m *mockDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if m.execContext != nil {
		return m.execContext(ctx, query, args...)
	}
	return nil, nil
}

func (m *mockDB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	if m.prepareContext != nil {
		return m.prepareContext(ctx, query)
	}
	return nil, nil
}

func (m *mockDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if m.queryContext != nil {
		return m.queryContext(ctx, query, args...)
	}
	return nil, nil
}

func (m *mockDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if m.queryRowContext != nil {
		return m.queryRowContext(ctx, query, args...)
	}
	return &sql.Row{}
}

func TestNew(t *testing.T) {
	mock := &mockDB{}
	queries := New(mock)
	assert.NotNil(t, queries)
	assert.Equal(t, mock, queries.db)
}

func TestWithTx(t *testing.T) {
	mock := &mockDB{}
	queries := New(mock)

	tx := &sql.Tx{}
	queriesWithTx := queries.WithTx(tx)

	assert.NotNil(t, queriesWithTx)
	assert.Equal(t, tx, queriesWithTx.db)
	assert.NotEqual(t, queries.db, queriesWithTx.db)
}

// mockResult implements sql.Result for testing
type mockResult struct {
	lastInsertId int64
	rowsAffected int64
	err          error
}

func (m *mockResult) LastInsertId() (int64, error) {
	return m.lastInsertId, m.err
}

func (m *mockResult) RowsAffected() (int64, error) {
	return m.rowsAffected, m.err
}

func TestQueriesWithMockDB(t *testing.T) {
	ctx := context.Background()
	expectedResult := &mockResult{lastInsertId: 1, rowsAffected: 1}

	mock := &mockDB{
		execContext: func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
			return expectedResult, nil
		},
	}

	queries := New(mock)
	assert.NotNil(t, queries)

	// Test ExecContext through the mock
	result, err := queries.db.ExecContext(ctx, "INSERT INTO test (col) VALUES (?)", "value")
	require.NoError(t, err)

	lastID, err := result.LastInsertId()
	require.NoError(t, err)
	assert.Equal(t, expectedResult.lastInsertId, lastID)

	rowsAff, err := result.RowsAffected()
	require.NoError(t, err)
	assert.Equal(t, expectedResult.rowsAffected, rowsAff)
}
