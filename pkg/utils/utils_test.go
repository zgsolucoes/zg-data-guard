package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Group string
	Value int
}

func TestGroupByProperty(t *testing.T) {
	items := []TestStruct{
		{Group: "A", Value: 1},
		{Group: "B", Value: 2},
		{Group: "B", Value: 9},
		{Group: "A", Value: 3},
		{Group: "B", Value: 4},
		{Group: "C", Value: 5},
	}

	grouped := GroupByProperty(items, func(item TestStruct) string {
		return item.Group
	})

	assert.Len(t, grouped, 3, "There should be three groups")
	assert.Len(t, grouped["C"], 1, "Group C should have one item")
	assert.Len(t, grouped["A"], 2, "Group A should have two items")
	assert.Len(t, grouped["B"], 3, "Group B should have three items")
}

func TestStringNotEmpty(t *testing.T) {
	err := StringNotEmpty(" ", "name")
	assert.Error(t, err, "Should return error for whitespace string")
	assert.EqualError(t, err, "param 'name' is required", "Error message should be correct")

	err = StringNotEmpty("", "age")
	assert.Error(t, err, "Should return error for empty string")
	assert.EqualError(t, err, "param 'age' is required", "Error message should be correct")

	err = StringNotEmpty("Hello, World!", "test")
	assert.NoError(t, err, "Should not return error for non-empty string")
}

func TestValidEmail(t *testing.T) {
	assert.False(t, ValidEmail("test-invalid-email.com"), "Should return false for invalid email")
	assert.False(t, ValidEmail("test"), "Should return false for invalid email")
	assert.False(t, ValidEmail(""), "Should return false for empty email")

	assert.True(t, ValidEmail("test@email"), "Should return true for valid email")
	assert.True(t, ValidEmail("test@email.com"), "Should return true for valid email")
	assert.True(t, ValidEmail("test@gov.com.br"), "Should return true for valid email")
	assert.True(t, ValidEmail("test@net"), "Should return true for valid email")
}

func TestContains(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}

	assert.True(t, Contains(items, 1), "Should return true for existing item")
	assert.True(t, Contains(items, 3), "Should return true for existing item")
	assert.True(t, Contains(items, 5), "Should return true for existing item")

	assert.False(t, Contains(items, 0), "Should return false for non-existing item")
	assert.False(t, Contains(items, 6), "Should return false for non-existing item")
	assert.False(t, Contains(items, 10), "Should return false for non-existing item")
}
