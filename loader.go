package dbuffer

// Loader wraps the interface
type Loader[T any] interface {
	// Load data to dest pointer. ok is true if data updated.
	Load(dest *T) (ok bool, err error)
}
