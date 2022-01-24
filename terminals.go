package parsegen

var anyPrefix = []byte("any")
var emptyTerm = []byte("empty")

func termStr(s string) tFunc {
	b := []byte(s)
	return func(it Iterator) code {
		for _, c := range b {
			if c != it.CC() {
				return missed
			}
			it.GC()
		}
		return zero
	}
}

// any(c), any[c], where c is character
func termBasicAny() tFunc {
	return func(it Iterator) code {
		for _, c := range anyPrefix {
			if c != it.CC() {
				return missed
			}
			it.GC()
		}
		c1 := it.CC()
		it.GC()
		if it.CC() != '0' {
			return missed
		}
		it.GC()
		if it.CC() != 'x' {
			return missed
		}
		it.GC()
		if !isHexDigit(it.CC()) {
			return missed
		}
		it.GC()
		if !isHexDigit(it.CC()) {
			return missed
		}
		it.GC()
		c2 := it.CC()
		if !(c1 == '(' && c2 == ')' || c1 == '[' && c2 == ']') {
			return missed
		}
		it.GC()
		return zero
	}
}

func isHexDigit(b byte) bool {
	return (b >= '0' && b <= '9') || (b >= 'A' && b <= 'F') || (b >= 'a' && b <= 'f')
}

func isHex(h []byte, b1, b2 *byte) bool {
	if len(h) >= 4 && h[0] == '0' && h[1] == 'x' && isHexDigit(h[2]) && isHexDigit(h[3]) {
		*b1 = toHex(h[2], h[3])
		*b2 = *b1
		if len(h) == 7 && h[4] == '-' && isHexDigit(h[5]) && isHexDigit(h[6]) {
			*b2 = toHex(h[5], h[6])
		}
		return true
	}
	return false
}

func toHex(n1, n2 byte) byte {
	n1 = toHexDigit(n1)
	n2 = toHexDigit(n2)
	return n1*16 + n2
}

func toHexDigit(n byte) byte {
	if isDigit(n) {
		return n - '0'
	}
	if n >= 'A' && n <= 'F' {
		return n - 'A' + 10
	}
	return n - 'a' + 10
}

func isAlpha(b byte) bool {
	return b >= 'A' && b <= 'Z' || b >= 'a' && b <= 'z'
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func termID() tFunc {
	return func(it Iterator) code {
		for !it.EOF() && (isAlpha(it.CC()) || isDigit(it.CC()) || it.CC() == '_' || it.CC() == '-') {
			it.GC()
		}

		return zero
	}
}

func basicHex() tFunc {
	return func(it Iterator) code {
		if it.CC() != '0' {
			return missed
		}
		it.GC()
		if it.CC() != 'x' {
			return missed
		}
		it.GC()
		if !isHexDigit(it.CC()) {
			return missed
		}
		it.GC()
		if !isHexDigit(it.CC()) {
			return missed
		}
		it.GC()
		// it may be 0x11-ff
		if it.CC() != '-' {
			return zero
		}
		it.GC()
		if !isHexDigit(it.CC()) {
			return missed
		}
		it.GC()
		if !isHexDigit(it.CC()) {
			return missed
		}
		it.GC()
		return zero
	}
}

func termHex(b byte) tFunc {
	return func(it Iterator) code {
		if b != it.CC() {
			return missed
		}
		it.GC()
		return zero
	}
}

func termHexes(b1, b2 byte) tFunc {
	return func(it Iterator) code {
		c := it.CC()
		if b1 > c || c > b2 {
			return missed
		}
		it.GC()
		return zero
	}
}

func termSpace() tFunc {
	return func(it Iterator) code {
		for !it.EOF() && IsSpace(it.CC()) {
			it.GC()
		}

		if it.EOF() {
			return eof
		}
		return zero
	}
}

// end symbol
// Any(:)
func termAny(end byte, includeEnd bool) tFunc {
	return func(it Iterator) code {
		for !it.EOF() && it.CC() != end {
			it.GC()
		}
		// next call CC() returns byte end of any
		if includeEnd {
			it.GC()
		}

		return zero
	}
}

func termEmpty() tFunc {
	return func(it Iterator) code {
		return exit
	}
}

func termAnyQuoted() tFunc {
	return func(it Iterator) code {
		if it.CC() != '"' {
			// it.SetError("AnyQuoted[T]: not beginning quote")
			return missed
		}
		it.GC()
		for !it.EOF() {
			if it.CC() == '"' {
				break
			}
			it.GC()
		}
		it.GC()
		return zero
	}
}
