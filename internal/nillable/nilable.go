package nillable

import "reflect"

func NonNil[T any](t *T) T {
	if t != nil {
		return *t
	}
	var zero T
	return zero
}

func RequireNonNil[T any](t *T) T {
	if t == nil {
		panic("value " + reflect.TypeFor[T]().Name() + " require non-nil but the value is nil")
	}
	return *t
}
