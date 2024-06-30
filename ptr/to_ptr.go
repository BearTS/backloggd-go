package ptr

// String returns a pointer value for the string value passed in.
func String(v string) *string {
	return &v
}
