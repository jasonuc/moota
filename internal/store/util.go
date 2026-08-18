package store

func nullIfEmpty[T comparable](v T) *T {
	var zero T

	if v == zero {
		return nil
	}

	return &v
}
