package nillabe

import "reflect"

func Nillable[T any](t *T) *T {
	if t == nil {
		return nil
	}
	return t
}

func NonNil[T any](t *T) T {
	if t != nil {
		return *t
	}
	return *new(T)
}

func RequireNonNil[T any](t *T) T {
	if t == nil {
		panic("value " + reflect.TypeFor[T]().Name() + " require non-nil but the value is nil")
	}
	return *t
}
