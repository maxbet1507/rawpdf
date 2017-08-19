package rawpdf

import (
	"bytes"
	"encoding/hex"
	"strings"
)

type line struct {
	data   []byte
	offset int
	length int
}

func (s line) Bytes() []byte {
	return s.data[s.offset : s.offset+s.length]
}

func (s line) Hex() string {
	return hex.EncodeToString(s.Bytes())
}

func (s line) String() string {
	return strings.TrimRight(string(s.Bytes()), " \r\n")
}

// type line []byte
type lines []line

func (s lines) Join() []byte {
	buf := new(bytes.Buffer)

	for _, val := range s {
		// bytes.Bufferはエラーを返さない
		buf.Write(val.Bytes())
	}

	return buf.Bytes()
}

func newLines(v []byte) lines {
	ret := lines{}

	offset := 0
	length := 0
	for _, val := range v {
		length++
		if val == 0x0a {
			ret = append(ret, line{
				data:   v,
				offset: offset,
				length: length,
			})
			offset, length = offset+length, 0
		}
	}
	if length > 0 {
		ret = append(ret, line{
			data:   v,
			offset: offset,
			length: length,
		})
	}

	return ret
}
