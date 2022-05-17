package queryf

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
)

func isNull(arg any) bool {
	return arg == nil
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
	replaceTo := ""
	if isNull(arg) {
		replaceTo = "NULL"
	} else if isTime(arg) {
		t, _ := arg.(time.Time)
		replaceTo = t.Format(time.RFC3339)
	} else if isString(arg) {
		s, _ := arg.(string)
		replaceTo = fmt.Sprintf("'%s'", s)
	} else if isSlice(arg) {
		replaceTo = sliceToString(arg.([]any))
	} else if isBoolean(arg) {
		b, _ := arg.(bool)
		replaceTo = "false"
		if b {
			replaceTo = "true"
		}
	}
	return replaceTo
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
