package queryf

import (
	"database/sql"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		args     []any
		expected string
	}{
		{
			name:     "basic types",
			query:    "SELECT * FROM users WHERE id = $1 AND name = $2 AND active = $3",
			args:     []any{1, "John", true},
			expected: "SELECT * FROM users WHERE id = 1 AND name = 'John' AND active = true",
		},
		{
			name:     "null value",
			query:    "SELECT * FROM users WHERE id = $1 AND name = $2",
			args:     []any{1, nil},
			expected: "SELECT * FROM users WHERE id = 1 AND name = NULL",
		},
		{
			name:     "pointer value",
			query:    "SELECT * FROM users WHERE id = $1",
			args:     []any{intPtr(42)},
			expected: "SELECT * FROM users WHERE id = 42",
		},
		{
			name:     "time value",
			query:    "SELECT * FROM users WHERE created_at > $1",
			args:     []any{time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)},
			expected: "SELECT * FROM users WHERE created_at > '2023-01-01T00:00:00Z'",
		},
		{
			name:     "slice value",
			query:    "SELECT * FROM users WHERE id = ANY($1)",
			args:     []any{[]int{1, 2, 3}},
			expected: "SELECT * FROM users WHERE id = ANY('{1,2,3}')",
		},
		{
			name:     "bytes value",
			query:    "SELECT * FROM users WHERE data = $1",
			args:     []any{[]byte{0x01, 0x02, 0x03}},
			expected: "SELECT * FROM users WHERE data = '\\x010203'",
		},
		{
			name:     "float value",
			query:    "SELECT * FROM users WHERE score > $1",
			args:     []any{42.5},
			expected: "SELECT * FROM users WHERE score > 42.5",
		},
		{
			name:     "string with quotes",
			query:    "SELECT * FROM users WHERE name = $1",
			args:     []any{"O'Reilly"},
			expected: "SELECT * FROM users WHERE name = 'O''Reilly'",
		},
		{
			name:     "pq string array",
			query:    "SELECT * FROM users WHERE tags = $1",
			args:     []any{pq.StringArray{"tag1", "tag2"}},
			expected: "SELECT * FROM users WHERE tags = '{\"tag1\",\"tag2\"}'",
		},
		{
			name:     "pq int array",
			query:    "SELECT * FROM users WHERE ids = $1",
			args:     []any{pq.Int64Array{1, 2, 3}},
			expected: "SELECT * FROM users WHERE ids = '{1,2,3}'",
		},
		{
			name:     "pq bool array",
			query:    "SELECT * FROM users WHERE flags = $1",
			args:     []any{pq.BoolArray{true, false, true}},
			expected: "SELECT * FROM users WHERE flags = '{t,f,t}'",
		},
		{
			name:     "sql.NullString valid",
			query:    "SELECT * FROM users WHERE name = $1",
			args:     []any{sql.NullString{String: "John", Valid: true}},
			expected: "SELECT * FROM users WHERE name = 'John'",
		},
		{
			name:     "sql.NullString invalid",
			query:    "SELECT * FROM users WHERE name = $1",
			args:     []any{sql.NullString{Valid: false}},
			expected: "SELECT * FROM users WHERE name = NULL",
		},
		{
			name:     "sql.NullInt64 valid",
			query:    "SELECT * FROM users WHERE id = $1",
			args:     []any{sql.NullInt64{Int64: 42, Valid: true}},
			expected: "SELECT * FROM users WHERE id = 42",
		},
		{
			name:     "sql.NullTime valid",
			query:    "SELECT * FROM users WHERE created_at = $1",
			args:     []any{sql.NullTime{Time: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true}},
			expected: "SELECT * FROM users WHERE created_at = '2023-01-01T00:00:00Z'",
		},
		{
			name:     "pq.NullTime valid",
			query:    "SELECT * FROM users WHERE created_at = $1",
			args:     []any{pq.NullTime{Time: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true}},
			expected: "SELECT * FROM users WHERE created_at = '2023-01-01T00:00:00Z'",
		},
		{
			name:     "map value",
			query:    "SELECT * FROM users WHERE data = $1",
			args:     []any{map[string]any{"name": "John", "age": 30}},
			expected: "SELECT * FROM users WHERE data = '{\"age\":30,\"name\":\"John\"}'",
		},
		{
			name:  "struct value",
			query: "SELECT * FROM users WHERE data = $1",
			args: []any{struct {
				Name string
				Age  int
			}{"John", 30}},
			expected: "SELECT * FROM users WHERE data = '{\"Name\":\"John\",\"Age\":30}'",
		},
		{
			name:  "struct with json tags",
			query: "SELECT * FROM users WHERE data = $1",
			args: []any{struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}{"John", 30}},
			expected: "SELECT * FROM users WHERE data = '{\"name\":\"John\",\"age\":30}'",
		},
		{
			name:     "multiple placeholders",
			query:    "SELECT * FROM users WHERE id IN ($1, $2, $3, $1)",
			args:     []any{1, 2, 3},
			expected: "SELECT * FROM users WHERE id IN (1, 2, 3, 1)",
		},
		{
			name:     "pq.NullTime invalid",
			query:    "SELECT * FROM users WHERE created_at = $1",
			args:     []any{pq.NullTime{Valid: false}},
			expected: "SELECT * FROM users WHERE created_at = NULL",
		},
		{
			name:     "pq float32 array",
			query:    "SELECT * FROM users WHERE scores = $1",
			args:     []any{pq.Float32Array{1.1, 2.2, 3.3}},
			expected: "SELECT * FROM users WHERE scores = '{1.1,2.2,3.3}'",
		},
		{
			name:     "pq float64 array",
			query:    "SELECT * FROM users WHERE scores = $1",
			args:     []any{pq.Float64Array{1.1, 2.2, 3.3}},
			expected: "SELECT * FROM users WHERE scores = '{1.1,2.2,3.3}'",
		},
		{
			name:     "pq int32 array",
			query:    "SELECT * FROM users WHERE ids = $1",
			args:     []any{pq.Int32Array{1, 2, 3}},
			expected: "SELECT * FROM users WHERE ids = '{1,2,3}'",
		},
		{
			name:     "pq bytea array",
			query:    "SELECT * FROM users WHERE data = $1",
			args:     []any{pq.ByteaArray{[]byte{1, 2, 3}, []byte{4, 5, 6}}},
			expected: "SELECT * FROM users WHERE data = '{\"\\\\x010203\",\"\\\\x040506\"}'",
		},
		{
			name:     "sql.NullFloat64 valid",
			query:    "SELECT * FROM users WHERE score = $1",
			args:     []any{sql.NullFloat64{Float64: 42.5, Valid: true}},
			expected: "SELECT * FROM users WHERE score = 42.5",
		},
		{
			name:     "sql.NullFloat64 invalid",
			query:    "SELECT * FROM users WHERE score = $1",
			args:     []any{sql.NullFloat64{Valid: false}},
			expected: "SELECT * FROM users WHERE score = NULL",
		},
		{
			name:     "generic array",
			query:    "SELECT * FROM users WHERE data = $1",
			args:     []any{pq.GenericArray{A: []int{1, 2, 3}}},
			expected: "SELECT * FROM users WHERE data = '{1,2,3}'",
		},
		{
			name:     "generic array nil",
			query:    "SELECT * FROM users WHERE data = $1",
			args:     []any{pq.GenericArray{A: nil}},
			expected: "SELECT * FROM users WHERE data = NULL",
		},
		{
			name:     "float32",
			query:    "SELECT * FROM users WHERE score = $1",
			args:     []any{float32(42.5)},
			expected: "SELECT * FROM users WHERE score = 42.5",
		},
		{
			name:     "boolean false",
			query:    "SELECT * FROM users WHERE active = $1",
			args:     []any{false},
			expected: "SELECT * FROM users WHERE active = false",
		},
		{
			name:     "sql.NullInt32 valid",
			query:    "SELECT * FROM users WHERE id = $1",
			args:     []any{sql.NullInt32{Int32: 42, Valid: true}},
			expected: "SELECT * FROM users WHERE id = 42",
		},
		{
			name:     "sql.NullInt16 valid",
			query:    "SELECT * FROM users WHERE id = $1",
			args:     []any{sql.NullInt16{Int16: 42, Valid: true}},
			expected: "SELECT * FROM users WHERE id = 42",
		},
		{
			name:     "sql.NullByte valid",
			query:    "SELECT * FROM users WHERE id = $1",
			args:     []any{sql.NullByte{Byte: 42, Valid: true}},
			expected: "SELECT * FROM users WHERE id = 42",
		},
		{
			name:     "sql.NullBool valid",
			query:    "SELECT * FROM users WHERE active = $1",
			args:     []any{sql.NullBool{Bool: true, Valid: true}},
			expected: "SELECT * FROM users WHERE active = true",
		},
		{
			name:     "sql.NullBool invalid",
			query:    "SELECT * FROM users WHERE active = $1",
			args:     []any{sql.NullBool{Valid: false}},
			expected: "SELECT * FROM users WHERE active = NULL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Format(tt.query, tt.args...)

			// Special case for map test since map iteration order is not guaranteed
			if tt.name == "map value" {
				// Check that the result contains both expected key-value pairs
				assert.Contains(t, result, "\"name\":\"John\"")
				assert.Contains(t, result, "\"age\":30")
				assert.Contains(t, result, "SELECT * FROM users WHERE data = '{")
				assert.Contains(t, result, "}'")
			} else {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestArgumentGetType(t *testing.T) {
	tests := []struct {
		name     string
		arg      any
		expected ParameterType
	}{
		{"null", nil, Null},
		{"string", "test", String},
		{"integer", 42, Integer},
		{"float", 42.5, Float},
		{"boolean", true, Boolean},
		{"pointer", intPtr(42), Pointer},
		{"time", time.Now(), Time},
		{"slice", []int{1, 2, 3}, Slice},
		{"bytes", []byte{1, 2, 3}, Bytes},
		{"pq string array", pq.StringArray{"a", "b"}, PqArray},
		{"pq int array", pq.Int64Array{1, 2}, PqArray},
		{"sql.NullString", sql.NullString{String: "test", Valid: true}, SqlNullType},
		{"map", map[string]any{"key": "value"}, Map},
		{"struct", struct{ Name string }{"John"}, Struct},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arg := NewArgument(tt.arg)
			assert.Equal(t, tt.expected, arg.GetType())
		})
	}
}

// Helper function to create integer pointers
func intPtr(i int) *int {
	return &i
}
