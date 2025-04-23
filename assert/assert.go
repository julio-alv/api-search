package assert

// Panics if the input error is not nil
func NotNil(err error) {
	if err != nil {
		panic(err)
	}
}

// Returns T, panics if the error is not nil
func Try[T any](result T, err error) T {
	if err != nil {
		panic(err)
	}
	return result
}
