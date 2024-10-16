package storage

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for *sql.DB
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Query(query string, args ...any) (*sql.Rows, error) {
	argsCalled := m.Called(query, args)
	return argsCalled.Get(0).(*sql.Rows), argsCalled.Error(1)
}

func TestReadSQLFile(t *testing.T) {
	t.Run("Successful Read", func(t *testing.T) {
		expectedContent := "SELECT * FROM test;"
		filePath := "test.sql"

		// Create a temporary file
		err := os.WriteFile(filePath, []byte(expectedContent), 0644)
		assert.NoError(t, err)
		defer os.Remove(filePath)

		content, err := ReadSQLFile(filePath)
		assert.NoError(t, err)
		assert.Equal(t, expectedContent, content)
	})

	t.Run("File Read Error", func(t *testing.T) {
		_, err := ReadSQLFile("non_existent_file.sql")
		assert.Error(t, err)
	})
}

func TestExecuteSQLQuery_Success(t *testing.T) {
	mockDB := new(MockDB)
	baseQuery := "SELECT * FROM test WHERE id = $1"
	args := []any{1}
	mockRows := new(sql.Rows)
	mockDB.On("Query", baseQuery, args).Return(mockRows, nil)

	rows, err := executeSQLQuery(mockDB, baseQuery, args)

	assert.NoError(t, err)
	assert.Equal(t, mockRows, rows)
}

func TestExecuteSQLQuery_QueryError(t *testing.T) {
	mockDB := new(MockDB)
	baseQuery := "SELECT * FROM test WHERE id = $1"
	args := []any{1}
	mockDB.On("Query", baseQuery, args).Return(&sql.Rows{}, sql.ErrConnDone)

	_, err := executeSQLQuery(mockDB, baseQuery, args)

	assert.Error(t, err, "error expected when query fails")
	assert.EqualError(t, err, sql.ErrConnDone.Error())
}

func TestAddFilterCondition(t *testing.T) {
	query := "SELECT * FROM test WHERE 1=1"
	var args []any

	t.Run("with value", func(t *testing.T) {
		newQuery, newArgs := addFilterCondition(query, args, "field", "value")
		assert.Equal(t, "SELECT * FROM test WHERE 1=1 AND field = $1", newQuery)
		assert.Equal(t, []any{"value"}, newArgs)
	})

	t.Run("without value", func(t *testing.T) {
		newQuery, newArgs := addFilterCondition(query, args, "field", "")
		assert.Equal(t, query, newQuery)
		assert.Equal(t, args, newArgs)
	})
}

func TestAppendFilterIdsInQuery(t *testing.T) {
	baseQuery := "SELECT * FROM test t WHERE 1=1"
	entityAlias := "t"
	var args []any

	t.Run("with ids", func(t *testing.T) {
		ids := []string{"1", "2", "3"}
		newQuery, newArgs := appendFilterIdsInQuery(baseQuery, entityAlias, ids, args)
		assert.Equal(t, "SELECT * FROM test t WHERE 1=1 AND t.id IN ($1, $2, $3)", newQuery)
		assert.Equal(t, []any{"1", "2", "3"}, newArgs)
	})

	t.Run("without ids", func(t *testing.T) {
		var ids []string
		newQuery, newArgs := appendFilterIdsInQuery(baseQuery, entityAlias, ids, args)
		assert.Equal(t, baseQuery, newQuery)
		assert.Equal(t, args, newArgs)
	})
}
