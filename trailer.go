package rawpdf

import (
	"fmt"
	"strconv"

	"github.com/maxbet1507/rawpdf/objects"

	"github.com/pkg/errors"
)

var (
	errInvalidTrailer = fmt.Errorf("invalid trailer")
)

type trailer struct {
	Start      int
	Dictionary map[string]interface{}
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
	dictionary, err := objects.Unmarshal(lines[beginpos:pos].Join())
	if _, ok := dictionary.(map[string]interface{}); !ok || err != nil {
		return nil, nil, errors.Wrap(errInvalidTrailer, "trailer dictionary")
	}
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
		Dictionary: dictionary.(map[string]interface{}),
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

	dictionary, err := objects.Unmarshal(lines[pos+1 : endpos+1].Join())
	if _, ok := dictionary.(map[string]interface{}); !ok || err != nil {
		return nil, nil, errors.Wrap(errInvalidTrailer, "trailer dictionary")
	}

	return &trailer{
		Start:      start,
		Dictionary: dictionary.(map[string]interface{}),
	}, lines[:pos], nil
}
