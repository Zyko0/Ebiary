package core

type Addressable[T any] interface {
	Addr() T
}

func Addr[T any](v T) T {
	if vv, ok := any(v).(Addressable[T]); ok {
		return vv.Addr()
	}
	return v
}
