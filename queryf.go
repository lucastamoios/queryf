package queryf

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/lib/pq"
)

type ParameterType string

const (
	String       ParameterType = "string"
	Integer      ParameterType = "integer"
	Float        ParameterType = "float"
	Boolean      ParameterType = "boolean"
	Pointer      ParameterType = "pointer"
	Null         ParameterType = "null"
	Time         ParameterType = "time"
	Slice        ParameterType = "slice"
	GenericArray ParameterType = "generic_array"
	PqArray      ParameterType = "pq_array"
	SqlNullType  ParameterType = "sql_null_type"
	Bytes        ParameterType = "bytes"
	Map          ParameterType = "map"
	Struct       ParameterType = "struct"
)

// Format will return the query with the arguments formatted.
// This will replace the $1, $2, etc. with the arguments given, similar to what the
// database/sql package does, but for debugging purposes.
//
//	** This is not meant to be used in production code. **
//	** Passing this resulting string to a database may lead to SQL injections. **
//
// Example:
//
//	query := "SELECT * FROM users WHERE id = $1 AND name = $2"
//	args := []any{1, "John"}
//	fmt.Println(Format(query, args...))
//	// Output: SELECT * FROM users WHERE id = 1 AND name = 'John'
func Format(query string, args ...any) string {
	queryb := []byte(query)
	for i, arg := range args {
		index := i + 1
		re := regexp.MustCompile(fmt.Sprintf(`\$%d\b`, index))
		queryb = re.ReplaceAll(queryb, []byte(NewArgument(arg).format()))
	}
	return string(queryb)
}

func NewArgument(arg any) *Argument {
	return &Argument{arg: arg}
}

type Argument struct {
	arg    any
	rValue *reflect.Value
	rType  *reflect.Type
}

func (a *Argument) getReflectedValue() reflect.Value {
	if a.rValue == nil {
		v := reflect.ValueOf(a.arg)
		a.rValue = &v
	}
	return *a.rValue
}

func (a *Argument) getReflectedType() reflect.Type {
	if a.rType == nil {
		t := reflect.TypeOf(a.arg)
		a.rType = &t
	}
	return *a.rType
}

// GetType will return the type of the argument.
func (a *Argument) GetType() ParameterType {
	if a.isNull() {
		return Null
	} else if a.isPtr() {
		return Pointer
	} else if a.isTime() {
		return Time
	} else if a.isString() {
		return String
	} else if a.isSlice() {
		return Slice
	} else if a.isBytes() {
		return Bytes
	} else if a.isPqArray() {
		return PqArray
	} else if a.isGenericArray() {
		return GenericArray
	} else if a.isSqlNullType() {
		return SqlNullType
	} else if a.isBoolean() {
		return Boolean
	} else if a.isFloat() {
		return Float
	} else if a.isMap() {
		return Map
	} else if a.isStruct() {
		return Struct
	}
	return Integer
}

func (a *Argument) isNull() bool {
	return a.arg == nil || (a.getReflectedValue().Kind() == reflect.Ptr && a.getReflectedValue().IsNil())
}

func (a *Argument) isPtr() bool {
	return a.getReflectedValue().Kind() == reflect.Ptr && !a.getReflectedValue().IsNil()
}

func (a *Argument) isTime() bool {
	_, ok := a.arg.(time.Time)
	return ok
}

func (a *Argument) isString() bool {
	return a.getReflectedType().Kind() == reflect.String
}

func (a *Argument) isBoolean() bool {
	return a.getReflectedType().Kind() == reflect.Bool
}

func (a *Argument) isFloat() bool {
	kind := a.getReflectedType().Kind()
	return kind == reflect.Float32 || kind == reflect.Float64
}

func (a *Argument) isSlice() bool {
	return a.getReflectedType().Kind() == reflect.Slice && !a.isBytes() && !a.isPqArray()
}

func (a *Argument) isBytes() bool {
	_, ok := a.arg.([]byte)
	return ok
}

func (a *Argument) isMap() bool {
	return a.getReflectedType().Kind() == reflect.Map
}

func (a *Argument) isStruct() bool {
	return a.getReflectedType().Kind() == reflect.Struct && !a.isTime() && !a.isSqlNullType()
}

func (a *Argument) isGenericArray() bool {
	_, ok := a.arg.(pq.GenericArray)
	return ok
}

func (a *Argument) isPqArray() bool {
	switch a.arg.(type) {
	case pq.BoolArray, pq.ByteaArray, pq.Float32Array, pq.Float64Array,
		pq.Int32Array, pq.Int64Array, pq.StringArray:
		return true
	default:
		return false
	}
}

func (a *Argument) isSqlNullType() bool {
	switch a.arg.(type) {
	case sql.NullBool, sql.NullByte, sql.NullFloat64, sql.NullInt16,
		sql.NullInt32, sql.NullInt64, sql.NullString, sql.NullTime,
		pq.NullTime:
		return true
	default:
		return false
	}
}

func (a *Argument) format() string {
	if a.isNull() {
		return a.formatNull()
	} else if a.isPtr() {
		return a.formatPtr(a.getReflectedValue())
	} else if a.isTime() {
		return a.formatTime(a.arg)
	} else if a.isString() {
		return a.formatString(a.arg)
	} else if a.isBytes() {
		return a.formatBytes(a.arg)
	} else if a.isSlice() {
		return a.formatSlice()
	} else if a.isPqArray() {
		return a.formatPqArray(a.arg)
	} else if a.isGenericArray() {
		return a.formatGenericArray(a.arg)
	} else if a.isSqlNullType() {
		return a.formatSqlNullType(a.arg)
	} else if a.isBoolean() {
		return a.formatBoolean(a.arg)
	} else if a.isMap() {
		return a.formatMap()
	} else if a.isStruct() {
		return a.formatStruct()
	} else if a.isFloat() {
		return a.formatFloat(a.arg)
	}
	return fmt.Sprintf("%v", a.arg)
}

func (a *Argument) formatSlice() string {
	var result []string
	for i := 0; i < a.getReflectedValue().Len(); i++ {
		newArg := NewArgument(a.getReflectedValue().Index(i).Interface())
		result = append(result, newArg.format())
	}
	return fmt.Sprintf("'{%s}'", strings.Join(result, ","))
}

func (a *Argument) formatNull() string {
	return "NULL"
}

func (a *Argument) formatPtr(rv reflect.Value) string {
	return NewArgument(rv.Elem().Interface()).format()
}

func (a *Argument) formatTime(arg any) string {
	t, _ := arg.(time.Time)
	return fmt.Sprintf("'%s'", t.Format(time.RFC3339))
}

func (a *Argument) formatString(arg any) string {
	s, _ := arg.(string)
	// Escape single quotes for SQL
	s = strings.ReplaceAll(s, "'", "''")
	return fmt.Sprintf("'%s'", s)
}

func (a *Argument) formatBytes(arg any) string {
	b, _ := arg.([]byte)
	// Format as bytea literal
	return fmt.Sprintf("'\\x%x'", b)
}

func (a *Argument) formatBoolean(arg any) string {
	b, _ := arg.(bool)
	if b {
		return "true"
	}
	return "false"
}

func (a *Argument) formatFloat(arg any) string {
	switch v := arg.(type) {
	case float32:
		return fmt.Sprintf("%g", v)
	case float64:
		return fmt.Sprintf("%g", v)
	default:
		return fmt.Sprintf("%v", arg)
	}
}

func (a *Argument) formatGenericArray(arg any) string {
	m, ok := arg.(string)
	if ok {
		return m
	}
	n, err := arg.(pq.GenericArray).Value()
	if err != nil || n == nil {
		return "NULL"
	}
	strVal := n.(string)
	if !strings.HasPrefix(strVal, "'") {
		strVal = "'" + strVal + "'"
	}
	return strVal
}

func (a *Argument) formatPqArray(arg any) string {
	var valuer driver.Valuer
	var ok bool

	// Handle different PostgreSQL array types
	switch arg.(type) {
	case pq.BoolArray, pq.ByteaArray, pq.Float32Array, pq.Float64Array,
		pq.Int32Array, pq.Int64Array, pq.StringArray:
		valuer, ok = arg.(driver.Valuer)
	default:
		return fmt.Sprintf("%v", arg)
	}

	if !ok {
		return fmt.Sprintf("%v", arg)
	}

	val, err := valuer.Value()
	if err != nil || val == nil {
		return "NULL"
	}

	// Ensure the array value is properly quoted
	strVal := val.(string)
	if !strings.HasPrefix(strVal, "'") {
		strVal = "'" + strVal + "'"
	}
	return strVal
}

func (a *Argument) formatSqlNullType(arg any) string {
	// Handle SQL null types
	switch v := arg.(type) {
	case sql.NullBool:
		if !v.Valid {
			return "NULL"
		}
		return a.formatBoolean(v.Bool)
	case sql.NullByte:
		if !v.Valid {
			return "NULL"
		}
		return fmt.Sprintf("%d", v.Byte)
	case sql.NullFloat64:
		if !v.Valid {
			return "NULL"
		}
		return a.formatFloat(v.Float64)
	case sql.NullInt16:
		if !v.Valid {
			return "NULL"
		}
		return fmt.Sprintf("%d", v.Int16)
	case sql.NullInt32:
		if !v.Valid {
			return "NULL"
		}
		return fmt.Sprintf("%d", v.Int32)
	case sql.NullInt64:
		if !v.Valid {
			return "NULL"
		}
		return fmt.Sprintf("%d", v.Int64)
	case sql.NullString:
		if !v.Valid {
			return "NULL"
		}
		return a.formatString(v.String)
	case sql.NullTime:
		if !v.Valid {
			return "NULL"
		}
		return a.formatTime(v.Time)
	case pq.NullTime:
		if !v.Valid {
			return "NULL"
		}
		return a.formatTime(v.Time)
	default:
		return fmt.Sprintf("%v", arg)
	}
}

func (a *Argument) formatMap() string {
	// For JSON-like data
	rv := a.getReflectedValue()
	keys := rv.MapKeys()
	pairs := make([]string, 0, len(keys))

	for _, key := range keys {
		k := fmt.Sprintf("%v", key.Interface())
		v := NewArgument(rv.MapIndex(key).Interface()).format()
		// If the value is already quoted (starts with '), we need to handle it specially
		if strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'") {
			// Extract the value without the quotes
			innerValue := v[1 : len(v)-1]
			// Escape any double quotes
			innerValue = strings.ReplaceAll(innerValue, "\"", "\\\"")
			pairs = append(pairs, fmt.Sprintf("\"%s\":\"%s\"", k, innerValue))
		} else {
			pairs = append(pairs, fmt.Sprintf("\"%s\":%s", k, v))
		}
	}

	return fmt.Sprintf("'{%s}'", strings.Join(pairs, ","))
}

func (a *Argument) formatStruct() string {
	// For struct types, format as JSON-like object
	rv := a.getReflectedValue()
	rt := a.getReflectedType()
	pairs := make([]string, 0, rt.NumField())

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		fieldName := field.Name
		// Check for json tag
		if tag, ok := field.Tag.Lookup("json"); ok {
			parts := strings.Split(tag, ",")
			if parts[0] != "" && parts[0] != "-" {
				fieldName = parts[0]
			}
		}

		fieldValue := rv.Field(i)
		if fieldValue.CanInterface() {
			v := NewArgument(fieldValue.Interface()).format()
			// If the value is already quoted (starts with '), we need to handle it specially
			if strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'") {
				// Extract the value without the quotes
				innerValue := v[1 : len(v)-1]
				// Escape any double quotes
				innerValue = strings.ReplaceAll(innerValue, "\"", "\\\"")
				pairs = append(pairs, fmt.Sprintf("\"%s\":\"%s\"", fieldName, innerValue))
			} else {
				pairs = append(pairs, fmt.Sprintf("\"%s\":%s", fieldName, v))
			}
		}
	}

	return fmt.Sprintf("'{%s}'", strings.Join(pairs, ","))
}
