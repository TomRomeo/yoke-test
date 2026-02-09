package util

import (
	"reflect"
)

func ValueOrDefault[T any](value T, defaultValue T) T {
	v := reflect.ValueOf(value)
	if v.IsZero() {
		return defaultValue
	}
	return value

}
