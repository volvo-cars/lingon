package kubeutil

// P returns a pointer to the given value.
func P[T any](t T) *T {
	return &t
}
