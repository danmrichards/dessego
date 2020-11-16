package math

// MakeSignedInt returns a forcibly signed version of n, wrapping around at the
// maximum 32-bit integer value.
//
// For some reason, Demon's Souls isn't using signed integers for some of it's
// request params.
func MakeSignedInt(n int) int {
	if n >= (1 << 31) {
		return n - (1 << 32)
	} else {
		return n
	}
}
