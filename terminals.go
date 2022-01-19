package main

func termStr(s string) tFunc {
	b := []byte(s)
	return func(it Iterator) code {
		for _, c := range b {
			if c != it.CC() {
				// log.Printf("Str[T]: not matched %c != %c", c, it.CC())
				return missed
			}
			it.GC()
		}
		if it.EOF() {
			return eof
		}
		return zero
	}
}

func isAlpha(b byte) bool {
	return b >= 'A' && b <= 'Z' || b >= 'a' && b <= 'z'
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func termID() tFunc {
	return func(it Iterator) code {
		i := it.GP()
		for !it.EOF() && isAlpha(it.CC()) {
			it.GC()
		}
		// true if empty or eof
		if i == it.GP() {
			// it.SetError("ID[T]: empty")
			return missed
		}
		if it.EOF() {
			return eof
		}
		return zero
	}
}

func termEmpty() tFunc {
	return func(it Iterator) code {
		if it.EOF() {
			return eof
		}
		return zero
	}
}

func termInteger() tFunc {
	return func(it Iterator) code {
		i := it.GP()
		for !it.EOF() && isDigit(it.CC()) {
			it.GC()
		}
		if i == it.GP() {
			// it.SetError("Integer[T]: empty")
			return missed
		}
		if it.EOF() {
			return eof
		}
		return zero
	}
}

// end symbol
func termAny(end byte, includeEnd bool) tFunc {
	return func(it Iterator) code {
		i := it.GP()
		for !it.EOF() && it.CC() != end {
			it.GC()
		}
		// next call CC() returns byte end of any
		if includeEnd {
			it.GC()
		}
		if i == it.GP() {
			// it.SetError("Any[T]: empty")
			return missed
		}
		if it.EOF() {
			return eof
		}
		return zero
	}
}

func termAnyQuoted() tFunc {
	return func(it Iterator) code {
		if it.CC() != '"' {
			// it.SetError("AnyQuoted[T]: not beginning quote")
			return missed
		}
		it.GC()
		i := it.GP()
		for !it.EOF() {
			if it.CC() == '"' {
				break
			}
			it.GC()
		}
		it.GC()
		if i >= it.GP()-1 {
			// it.SetError("AnyQuoted[T]: empty")
			return missed
		}
		if it.EOF() {
			return eof
		}
		return zero
	}
}
