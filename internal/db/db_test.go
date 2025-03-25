package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDB implements DBTX interface for testing
type MockDB struct {
	mock.Mock
}

func (m *MockDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	args = append([]interface{}{ctx, query}, args...)
	ret := m.Called(args...)
	result, _ := ret.Get(0).(sql.Result)
	return result, ret.Error(1)
}

func (m *MockDB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	ret := m.Called(ctx, query)
	stmt, _ := ret.Get(0).(*sql.Stmt)
	return stmt, ret.Error(1)
}

func (m *MockDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	args = append([]interface{}{ctx, query}, args...)
	ret := m.Called(args...)
	rows, _ := ret.Get(0).(*sql.Rows)
	return rows, ret.Error(1)
}

func (m *MockDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	args = append([]interface{}{ctx, query}, args...)
	ret := m.Called(args...)
	row, _ := ret.Get(0).(*sql.Row)
	return row
}

// MockResult implements sql.Result for testing
type MockResult struct {
	mock.Mock
}

func (m *MockResult) LastInsertId() (int64, error) {
	ret := m.Called()
	return ret.Get(0).(int64), ret.Error(1)
}

func (m *MockResult) RowsAffected() (int64, error) {
	ret := m.Called()
	return ret.Get(0).(int64), ret.Error(1)
}

func TestNew(t *testing.T) {
	mockDB := new(MockDB)
	queries := New(mockDB)

	assert.NotNil(t, queries)
	assert.Equal(t, mockDB, queries.db)
}

func TestQueriesWithTx(t *testing.T) {
	mockDB := new(MockDB)
	queries := New(mockDB)

	tx := &sql.Tx{}
	queriesWithTx := queries.WithTx(tx)

	assert.NotNil(t, queriesWithTx)
	assert.Equal(t, tx, queriesWithTx.db)
	assert.NotEqual(t, queries.db, queriesWithTx.db)
}
