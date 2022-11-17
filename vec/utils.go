package vec

func getFromIndex(length int, fromIndex ...int) int {
	if len(fromIndex) > 0 {
		return fixIndex(length, fromIndex[0], true)
	}
	return 0
}

func fixRange(length, start int, end ...int) (fixedStart, fixedEnd int, ok bool) {
	fixedStart = fixIndex(length, start, true)
	if fixedStart == length {
		return
	}
	fixedEnd = length
	if len(end) > 0 {
		fixedEnd = fixIndex(length, end[0], true)
	}
	if fixedEnd-fixedStart <= 0 {
		return
	}
	ok = true
	return
}

func fixIndex(length int, idx int, canLen bool) int {
	if idx < 0 {
		idx = length + idx
		if idx < 0 {
			return 0
		}
		return idx
	}
	if idx >= length {
		if canLen {
			return length
		}
		return length - 1
	}
	return idx
}
