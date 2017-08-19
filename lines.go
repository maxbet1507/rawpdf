package rawpdf

import (
	"bytes"
	"strings"
)

type line []byte
type lines []line

func (s line) String() string {
	return strings.TrimRight(string(s), "\r\n")
}

func (s lines) Join() []byte {
	lst := [][]byte{}
	for _, val := range s {
		lst = append(lst, val)
	}

	return bytes.Join(lst, nil)
}

func newLines(v []byte) lines {
	ret := lines{}

	vals := bytes.Split(v, []byte{0x0a})
	for idx, val := range vals {
		if idx != len(vals)-1 {
			ret = append(ret, line(append(val, 0x0a)))
		} else {
			ret = append(ret, line(val))
		}
	}

	return ret
}
