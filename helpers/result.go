package helpers

type Result[T any] struct {
	ok  T
	err error
}

func Ok[T any](ok T) Result[T] {
	return Result[T]{ok: ok}
}

func Error[T any](err error) Result[T] {
	return Result[T]{err: err}
}

func (r *Result[T]) Unwrap() (T, error) {
	return r.ok, r.err
}

func (r *Result[T]) UnwrapOr(value T) T {
	if r.err != nil {
		return value
	}
	return r.ok
}
