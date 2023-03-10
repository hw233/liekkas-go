package naming

import "bytes"

func isLowerByte(r byte) bool {
	return r >= 'a' && r <= 'z'
}

func isUpperByte(r byte) bool {
	return r >= 'A' && r <= 'Z'
}

func toLowerByte(r byte) byte {
	return r + ('a' - 'A')
}

func toUpperByte(r byte) byte {
	return r - ('a' - 'A')
}

// underline naming to hump naming
// aa_bb_cc to AaBbCc
// aa_Bb_cc to AaBbCc
func HumpNaming(name string) string {
	if len(name) <= 0 {
		return name
	}

	var buf bytes.Buffer

	for i := 0; i < len(name); i++ {
		// this is '_'
		if name[i] == '_' {
			continue
		}

		// the first byte
		if (i == 0 ||
			// last is '_'
			(i > 0 && name[i-1] == '_')) &&
			// this is lower
			isLowerByte(name[i]) {
			buf.WriteByte(toUpperByte(name[i]))
		} else {
			buf.WriteByte(name[i])
		}
	}

	return buf.String()
}

func FirstLower(name string) string {
	if len(name) <= 0 {
		return name
	}
	var buf bytes.Buffer
	if isUpperByte(name[0]) {
		buf.WriteByte(toLowerByte(name[0]))
	}
	for i := 1; i < len(name); i++ {
		buf.WriteByte(name[i])
	}

	return buf.String()
}

// hump naming to underline naming
// AaBbCc to aa_bb_cc
// AaBBCc to aa_bb_cc
func UnderlineNaming(name string) string {
	if len(name) <= 0 {
		return name
	}

	var buf bytes.Buffer

	for i := 0; i < len(name); i++ {
		// this is lower
		if isLowerByte(name[i]) {
			buf.WriteByte(name[i])
			continue
		}

		// this is upper
		// not the first byte
		if i > 1 && (
		// last byte is lower
		isLowerByte(name[i-1]) ||
			// next byte is lower
			(i+1 < len(name) && isLowerByte(name[i+1]))) {
			buf.WriteByte('_')
		}

		buf.WriteByte(toLowerByte(name[i]))
	}

	return buf.String()
}
