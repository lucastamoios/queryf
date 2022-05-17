package queryf

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
)

func isNull(arg any) bool {
	return arg == nil || (reflect.ValueOf(arg).Kind() == reflect.Ptr && reflect.ValueOf(arg).IsNil())
}

func isPtr(arg any) bool {
	return reflect.ValueOf(arg).Kind() == reflect.Ptr
}

func isTime(arg any) bool {
	_, ok := arg.(time.Time)
	return ok
}

func isString(arg any) bool {
	return reflect.TypeOf(arg).Kind() == reflect.String
}

func isBoolean(arg any) bool {
	return reflect.TypeOf(arg).Kind() == reflect.Bool
}

func isSlice(arg any) bool {
	return reflect.TypeOf(arg).Kind() == reflect.Slice
}

func sliceToString(input []any) string {
	var result []string
	for _, arg := range input {
		result = append(result, format(arg))
	}
	return fmt.Sprintf("{%s}", strings.Join(result, ","))
}

func format(arg any) string {
	if isNull(arg) {
		return "NULL"
	} else if isPtr(arg) {
		v := reflect.ValueOf(arg).Elem()
		return format(v)
	} else if isTime(arg) {
		t, _ := arg.(time.Time)
		return format(t.Format(time.RFC3339))
	} else if isString(arg) {
		s, _ := arg.(string)
		return fmt.Sprintf("'%s'", s)
	} else if isSlice(arg) {
		return sliceToString(arg.([]any))
	} else if isBoolean(arg) {
		b, _ := arg.(bool)
		if b {
			return "true"
		}
		return "false"
	}
	return fmt.Sprintf("%v", arg)
}

func Print(query string, args ...any) string {
	queryb := []byte(query)
	for i, arg := range args {
		index := i + 1
		re := regexp.MustCompile(fmt.Sprintf(`\$%d\b`, index))
		queryb = re.ReplaceAll(queryb, []byte(format(arg)))
	}
	return string(queryb)
}
