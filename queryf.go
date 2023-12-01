package queryf

import (
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
	Boolean      ParameterType = "boolean"
	Pointer      ParameterType = "pointer"
	Null         ParameterType = "null"
	Time         ParameterType = "time"
	Slice        ParameterType = "slice"
	GenericArray ParameterType = "generic_array"
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
	} else if a.isGenericArray() {
		return GenericArray
	} else if a.isBoolean() {
		return Boolean
	}
	return Integer
}

func (a *Argument) isNull() bool {
	return a.arg == nil || (a.getReflectedValue().Kind() == reflect.Ptr && a.getReflectedValue().IsNil())
}

func (a *Argument) isPtr() bool {
	return a.getReflectedValue().Kind() == reflect.Ptr
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

func (a *Argument) isSlice() bool {
	return a.getReflectedType().Kind() == reflect.Slice
}

func (a *Argument) isGenericArray() bool {
	_, ok := a.arg.(pq.GenericArray)
	return ok
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
	} else if a.isSlice() {
		return a.formatSlice()
	} else if a.isGenericArray() {
		return a.formatGenericArray(a.arg)
	} else if a.isBoolean() {
		return a.formatBoolean(a.arg)
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
	return NewArgument(rv.Elem()).format()
}

func (a *Argument) formatTime(arg any) string {
	t, _ := arg.(time.Time)
	return NewArgument(t.Format(time.RFC3339)).format()
}

func (a *Argument) formatString(arg any) string {
	s, _ := arg.(string)
	return fmt.Sprintf("'%s'", s)
}

func (a *Argument) formatBoolean(arg any) string {
	b, _ := arg.(bool)
	if b {
		return "true"
	}
	return "false"
}

func (a *Argument) formatGenericArray(arg any) string {
	m, ok := arg.(string)
	if ok {
		return m
	}
	n, _ := arg.(pq.GenericArray).Value()
	if n == nil {
		return "NULL"
	}
	return n.(string)
}
