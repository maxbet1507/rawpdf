package rawpdf

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

var (
	errInvalidXref = fmt.Errorf("invalid xref")

	reXrefEntry = regexp.MustCompile("^(\\d{10}) (\\d{5}) ([n|f]) *$")
	reXrefRange = regexp.MustCompile("^(\\d+) (\\d+)$")
)

type xrefEntry struct {
	Offset     int
	Generation int
	InUse      bool
}

func newXrefEntry(line line) (*xrefEntry, error) {
	if line.length != 20 {
		return nil, errors.Wrap(errInvalidXref, "check length")
	}

	matches := reXrefEntry.FindStringSubmatch(line.String())
	if matches == nil {
		return nil, errors.Wrap(errInvalidXref, "match entry")
	}

	ret := &xrefEntry{}
	ret.Offset, _ = strconv.Atoi(matches[1])
	ret.Generation, _ = strconv.Atoi(matches[2])
	ret.InUse = matches[3] == "n"

	return ret, nil
}

type xref map[int]*xrefEntry

func newXref(lines lines) (xref, lines, error) {
	if len(lines) < 2 {
		return nil, nil, errors.Wrap(errInvalidXref, "too short lines")
	}
	pos := 0

	if lines[pos].String() != "xref" {
		return nil, nil, errors.Wrap(errInvalidXref, "check xref")
	}
	pos++

	matches := reXrefRange.FindStringSubmatch(lines[pos].String())
	if matches == nil {
		return nil, nil, errors.Wrap(errInvalidXref, "match range")
	}
	pos++

	objnum, _ := strconv.Atoi(matches[1])
	nument, _ := strconv.Atoi(matches[2])
	ret := xref{}

	for i := 0; i < nument && pos < len(lines); i++ {
		ent, err := newXrefEntry(lines[pos])
		if err != nil {
			return nil, nil, err
		}
		pos++
		ret[objnum+i] = ent
	}

	return ret, lines[pos:], nil
}
