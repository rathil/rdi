package rdi

// Get a single dependency
func Get[T any](di DI) (T, error) {
	var result T
	err := di.Invoke(func(out T) { result = out })
	return result, err
}

// MustGet is like Get but panics if an error has occurred
func MustGet[T any](di DI) T {
	var result T
	di.MustInvoke(func(out T) { result = out })
	return result
}
