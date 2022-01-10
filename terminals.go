package main

func termStr(s string) tFunc {
	b := []byte(s)
	return func(it Iterator) bool {
		for _, c := range b {
			if c != it.CC() || it.EOF() {
				return true
			}
			it.GC()
		}
		return it.EOF()
	}
}

func isAlpha(b byte) bool {
	return b >= 'A' && b <= 'Z' || b >= 'a' && b <= 'z'
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func termID() tFunc {
	return func(it Iterator) bool {
		i := it.GP()
		for !it.EOF() && isAlpha(it.CC()) {
			it.GC()
		}
		// true if empty or eof
		if i == it.GP() {
			it.SetError("ID[T]: empty")
			return true
		}
		return it.EOF()
	}
}

func termEmpty() tFunc {
	return func(it Iterator) bool {
		return false
	}
}

func termInteger() tFunc {
	return func(it Iterator) bool {
		i := it.GP()
		for !it.EOF() && isDigit(it.CC()) {
			it.GC()
		}
		if i == it.GP() {
			it.SetError("Integer[T]: empty")
			return true
		}
		return it.EOF()
	}
}

// end symbol
func termAny(end byte, includeEnd bool) tFunc {
	return func(it Iterator) bool {
		i := it.GP()
		for !it.EOF() && it.CC() != end {
			it.GC()
		}
		// next call CC() returns byte end of any
		if includeEnd {
			it.GC()
		}
		if i == it.GP() {
			it.SetError("Any[T]: empty")
			return true
		}
		return it.EOF()
	}
}

func termAnyQuoted() tFunc {
	return func(it Iterator) bool {
		if it.CC() != '"' {
			it.SetError("AnyQuoted[T]: not beginning quote")
			return true
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
		if i == it.GP()-1 {
			it.SetError("AnyQuoted[T]: empty")
			return true
		}
		return it.EOF()
	}
}
