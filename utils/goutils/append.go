package goutils

// InPlaceAppend Assigning the result of `append` to the same slice variable is immune from the issue shown above.
func InPlaceAppend[T ~[]I, I any](old *T, elems ...I) {
	*old = append(*old, elems...)
}

// ImmutableAppend If the result of `append` needs to be reassigned, it needs to be done in an immutable way.
func ImmutableAppend[T ~[]I, I any](old T, elems ...I) T {
	res := make([]I, len(old)+len(elems))
	copy(res, old)
	copy(res[len(old):], elems)
	return res
}
