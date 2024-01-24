package utils

func BytesEqual(a []byte, b []byte) bool {
	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}
	aLen, bLen := len(a), len(b)
	if aLen != bLen {
		return false
	}
	for i := 0; i < aLen; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
