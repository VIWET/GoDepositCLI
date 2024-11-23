package helpers

// Generic option
type Option[T any] func(T) error

// Generic option slice
type Options[T any] []Option[T]

// Apply all the options
func (opts Options[T]) Apply(obj T) error {
	for _, opt := range opts {
		if err := opt(obj); err != nil {
			return err
		}
	}

	return nil
}
