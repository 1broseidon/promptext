package format

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTOONEncoder_Primitives(t *testing.T) {
	encoder := NewTOONEncoder()

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "simple string",
			input:    map[string]interface{}{"message": "hello"},
			expected: "message: hello",
		},
		{
			name:     "integer",
			input:    map[string]interface{}{"count": 42},
			expected: "count: 42",
		},
		{
			name:     "float",
			input:    map[string]interface{}{"price": 19.99},
			expected: "price: 19.99",
		},
		{
			name:     "boolean true",
			input:    map[string]interface{}{"active": true},
			expected: "active: true",
		},
		{
			name:     "boolean false",
			input:    map[string]interface{}{"disabled": false},
			expected: "disabled: false",
		},
		{
			name:     "null value",
			input:    map[string]interface{}{"data": nil},
			expected: "data: null",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encoder.Encode(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTOONEncoder_StringQuoting(t *testing.T) {
	encoder := NewTOONEncoder()

	tests := []struct {
		name     string
		input    map[string]interface{}
		expected string
	}{
		{
			name:     "empty string",
			input:    map[string]interface{}{"empty": ""},
			expected: "empty: \"\"",
		},
		{
			name:     "string with comma",
			input:    map[string]interface{}{"csv": "a,b,c"},
			expected: "csv: \"a,b,c\"",
		},
		{
			name:     "string with colon",
			input:    map[string]interface{}{"time": "10:30"},
			expected: "time: \"10:30\"",
		},
		{
			name:     "looks like boolean",
			input:    map[string]interface{}{"text": "true"},
			expected: "text: \"true\"",
		},
		{
			name:     "looks like number",
			input:    map[string]interface{}{"code": "42"},
			expected: "code: \"42\"",
		},
		{
			name:     "leading space",
			input:    map[string]interface{}{"msg": " hello"},
			expected: "msg: \" hello\"",
		},
		{
			name:     "trailing space",
			input:    map[string]interface{}{"msg": "hello "},
			expected: "msg: \"hello \"",
		},
		{
			name:     "no quotes needed",
			input:    map[string]interface{}{"name": "Alice"},
			expected: "name: Alice",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encoder.Encode(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTOONEncoder_MultilineStrings(t *testing.T) {
	encoder := NewTOONEncoder()

	input := map[string]interface{}{
		"code": "package main\n\nfunc main() {\n  println(\"hello\")\n}",
	}

	expected := `code: |
  package main

  func main() {
    println("hello")
  }`

	result, err := encoder.Encode(input)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestTOONEncoder_NestedObjects(t *testing.T) {
	encoder := NewTOONEncoder()

	input := map[string]interface{}{
		"user": map[string]interface{}{
			"id":   123,
			"name": "Alice",
			"settings": map[string]interface{}{
				"theme": "dark",
				"notifications": true,
			},
		},
	}

	expected := `user:
  id: 123
  name: Alice
  settings:
    notifications: true
    theme: dark`

	result, err := encoder.Encode(input)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestTOONEncoder_PrimitiveArrays(t *testing.T) {
	encoder := NewTOONEncoder()

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "string array",
			input:    map[string]interface{}{"tags": []string{"go", "cli", "tool"}},
			expected: "tags[3]: go,cli,tool",
		},
		{
			name:     "int array",
			input:    map[string]interface{}{"numbers": []int{1, 2, 3, 4, 5}},
			expected: "numbers[5]: 1,2,3,4,5",
		},
		{
			name:     "mixed primitive array",
			input:    map[string]interface{}{"data": []interface{}{"text", 42, true}},
			expected: "data[3]: text,42,true",
		},
		{
			name:     "empty array",
			input:    map[string]interface{}{"items": []string{}},
			expected: "items[0]:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encoder.Encode(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTOONEncoder_TabularArrays(t *testing.T) {
	encoder := NewTOONEncoder()

	type Product struct {
		SKU   string
		Name  string
		Price float64
	}

	input := map[string]interface{}{
		"products": []Product{
			{SKU: "A1", Name: "Widget", Price: 9.99},
			{SKU: "B2", Name: "Gadget", Price: 14.50},
			{SKU: "C3", Name: "Doohickey", Price: 7.25},
		},
	}

	// Fields appear in struct definition order: SKU, Name, Price
	expected := `products[3]{SKU,Name,Price}:
  A1,Widget,9.99
  B2,Gadget,14.5
  C3,Doohickey,7.25`

	result, err := encoder.Encode(input)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestTOONEncoder_TabularArraysWithMaps(t *testing.T) {
	encoder := NewTOONEncoder()

	input := map[string]interface{}{
		"items": []map[string]interface{}{
			{"id": 1, "name": "Alice", "active": true},
			{"id": 2, "name": "Bob", "active": false},
		},
	}

	expected := `items[2]{active,id,name}:
  true,1,Alice
  false,2,Bob`

	result, err := encoder.Encode(input)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestTOONEncoder_ListArrays(t *testing.T) {
	encoder := NewTOONEncoder()

	input := map[string]interface{}{
		"mixed": []interface{}{
			"simple string",
			map[string]interface{}{"key": "value"},
			42,
		},
	}

	expected := `mixed[3]:
  - simple string
  -
    key: value
  - 42`

	result, err := encoder.Encode(input)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestTOONEncoder_ComplexStructure(t *testing.T) {
	encoder := NewTOONEncoder()

	input := map[string]interface{}{
		"project": map[string]interface{}{
			"name":   "promptext",
			"version": "0.3.0",
			"tokens": 4250,
		},
		"languages": []string{"Go", "YAML", "Shell"},
		"files": []map[string]interface{}{
			{"path": "cmd/main.go", "lang": "go", "tokens": 450},
			{"path": "internal/config.go", "lang": "go", "tokens": 320},
		},
	}

	expected := `files[2]{lang,path,tokens}:
  go,cmd/main.go,450
  go,internal/config.go,320
languages[3]: Go,YAML,Shell
project:
  name: promptext
  tokens: 4250
  version: 0.3.0`

	result, err := encoder.Encode(input)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestTOONEncoder_EscapingInStrings(t *testing.T) {
	encoder := NewTOONEncoder()

	input := map[string]interface{}{
		"message": "Hello \"world\"",
		"path":    "C:\\Users\\test",
		"newline": "line1\nline2",
	}

	expected := `message: "Hello \"world\""
newline: |
  line1
  line2
path: "C:\\Users\\test"`

	result, err := encoder.Encode(input)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestTOONEncoder_EmptyValues(t *testing.T) {
	encoder := NewTOONEncoder()

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "empty map",
			input:    map[string]interface{}{"data": map[string]interface{}{}},
			expected: "data:",
		},
		{
			name:     "empty array",
			input:    map[string]interface{}{"items": []string{}},
			expected: "items[0]:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encoder.Encode(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTOONEncoder_StructWithTags(t *testing.T) {
	encoder := NewTOONEncoder()

	type User struct {
		ID       int    `toon:"id"`
		Name     string `toon:"name"`
		Internal string `toon:"-"` // Should be skipped
	}

	input := User{
		ID:       123,
		Name:     "Alice",
		Internal: "secret",
	}

	result, err := encoder.Encode(input)
	require.NoError(t, err)

	// Should contain id and name but not Internal
	assert.Contains(t, result, "id: 123")
	assert.Contains(t, result, "name: Alice")
	assert.NotContains(t, result, "Internal")
	assert.NotContains(t, result, "secret")
}

func TestTOONEncoder_ArrayWithQuotedStrings(t *testing.T) {
	encoder := NewTOONEncoder()

	input := map[string]interface{}{
		"patterns": []string{"*.go", "test-*", "foo,bar"},
	}

	expected := `patterns[3]: *.go,test-*,"foo,bar"`

	result, err := encoder.Encode(input)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestTOONEncoder_NonUniformArraysFallbackToList(t *testing.T) {
	encoder := NewTOONEncoder()

	// Arrays with objects that have different keys should use list format
	input := map[string]interface{}{
		"items": []map[string]interface{}{
			{"id": 1, "name": "Alice"},
			{"id": 2, "email": "bob@example.com"}, // Different keys
		},
	}

	result, err := encoder.Encode(input)
	require.NoError(t, err)

	// Should use list format, not tabular
	assert.Contains(t, result, "items[2]:")
	assert.Contains(t, result, "  -\n") // Objects use dash without trailing space
}

func TestTOONEncoder_ArrayWithNestedObjectsFallbackToList(t *testing.T) {
	encoder := NewTOONEncoder()

	// Tabular format requires primitive values only
	input := map[string]interface{}{
		"items": []map[string]interface{}{
			{"id": 1, "metadata": map[string]interface{}{"key": "value"}},
			{"id": 2, "metadata": map[string]interface{}{"key": "value"}},
		},
	}

	result, err := encoder.Encode(input)
	require.NoError(t, err)

	// Should use list format due to nested objects
	assert.Contains(t, result, "items[2]:")
	assert.Contains(t, result, "  -\n") // Objects use dash without trailing space
}
