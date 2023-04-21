package modules

import (
	"fmt"
)

type Range struct {
	Min int64
	Max int64
}

func (r *Range) BytesRange() string {
	return fmt.Sprintf("bytes=%d-%d", r.Min, r.Max)
}

func MakeRange(i, processes int, size, contentLength int64) Range {
	low := size * int64(i)

	if i == processes-1 {
		return Range{Min: low, Max: contentLength}
	}

	return Range{Min: low, Max: low + size - 1}
}
