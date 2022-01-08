package main

func termStr(s string) tFunc {
	b := []byte(s)
	return func(it *Iterator) bool {
		for _, c := range b {
			if c != it.cc() {
				return true
			}
			it.gc()
		}
		return false
	}
}

func termID() tFunc {
	return nil
}

func termEmpty() tFunc {
	return func(it *Iterator) bool {
		return !it.eof
	}
}

func termAnyQuoted() tFunc {
	return nil
}
