package xlinq

type sizeInfo struct {
	size     int
	sizeHint SizeHint
}

type SizeHint func(int) int

func SizeHintExact(s int) SizeHint {
	return func(int) int { return s }
}

func SizeHintRelative(ratio int) SizeHint {
	return func(size int) int { return size * ratio / 100 }
}

func withSize(size int) sizeInfo {
	return sizeInfo{
		size: size,
	}
}

func (si *sizeInfo) SetSizeHint(sizeHint SizeHint) {
	si.sizeHint = sizeHint
}

func (si *sizeInfo) SetSize(size int) {
	si.size = max(0, size)
}

func (si sizeInfo) CopyWithSize(size int) sizeInfo {
	// mutate copy
	si.SetSize(size)
	return si
}

func (si sizeInfo) CopyWithUnknownSize() sizeInfo {
	// mutate copy
	si.SetSize(0)
	return si
}

func (si sizeInfo) CopyWithSizeHint(sizeHint SizeHint) sizeInfo {
	// mutate copy
	si.SetSizeHint(sizeHint)
	return si
}

func (si sizeInfo) GetPreAllocSize() int {
	if si.size > 0 {
		return si.size
	}

	if si.sizeHint != nil {
		return si.sizeHint(si.size)
	}

	return 0
}

func (si sizeInfo) GetSize() int {
	return max(0, si.size)
}
