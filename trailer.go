package rawpdf

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

var (
	errInvalidTrailer = fmt.Errorf("invalid trailer")
)

type trailer struct {
	Start      int
	Dictionary string
}

func newFirstpageTrailer(lines lines) (*trailer, lines, error) {
	pos := 0
	if len(lines) < 4 {
		return nil, nil, errors.Wrap(errInvalidTrailer, "too short lines")
	}

	for lines[pos].String() != "trailer" {
		return nil, nil, errors.Wrap(errInvalidTrailer, "check trailer")
	}
	pos++

	beginpos := pos
	for lines[pos].String() != "startxref" {
		if pos++; pos >= len(lines) {
			return nil, nil, errors.Wrap(errInvalidTrailer, "check startref")
		}
	}
	dictionary := string(lines[beginpos:pos].Join())
	pos++

	start, err := strconv.Atoi(lines[pos].String())
	if err != nil {
		return nil, nil, errors.Wrap(errInvalidTrailer, "convert startxref")
	}
	pos++

	if eof := lines[pos]; eof.String() != "%%EOF" {
		return nil, nil, errors.Wrap(errInvalidTrailer, "check %%EOF")
	}

	return &trailer{
		Start:      start,
		Dictionary: dictionary,
	}, lines[pos+1:], nil
}

func newMainTrailer(lines lines) (*trailer, lines, error) {
	pos := len(lines)
	if pos < 4 {
		return nil, nil, errors.Wrap(errInvalidTrailer, "too short lines")
	}
	pos--

	if eof := lines[pos]; eof.String() != "%%EOF" {
		return nil, nil, errors.Wrap(errInvalidTrailer, "check %%EOF")
	}
	pos--

	start, err := strconv.Atoi(lines[pos].String())
	if err != nil {
		return nil, nil, errors.Wrap(errInvalidTrailer, "convert startxref")
	}
	pos--

	if startxref := lines[pos]; startxref.String() != "startxref" {
		return nil, nil, errors.Wrap(errInvalidTrailer, "check startref")
	}
	pos--

	endpos := pos
	for lines[pos].String() != "trailer" {
		if pos--; pos < 0 {
			return nil, nil, errors.Wrap(errInvalidTrailer, "check trailer")
		}
	}

	return &trailer{
		Start:      start,
		Dictionary: string(lines[pos+1 : endpos+1].Join()),
	}, lines[:pos], nil
}
