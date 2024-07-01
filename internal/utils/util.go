package utils

import (
	"fmt"
	"os"
	"reflect"
)

func GetEnv(key string) string {
	if env, exists := os.LookupEnv(key); exists {
		return env
	}

	return ""
}

func ToString(i interface{}) string {
	val := reflect.ValueOf(i)

	switch val.Kind() {
	case reflect.String:
		return val.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return fmt.Sprintf("%d", val.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", val.Float())
	case reflect.Bool:
		return fmt.Sprintf("%t", val.Bool())
	default:
		return fmt.Sprintf("%v", i)
	}
}
