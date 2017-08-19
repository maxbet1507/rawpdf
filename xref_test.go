package rawpdf

import (
	"bytes"
	"testing"

	"github.com/pkg/errors"
)

func TestXrefEntry(t *testing.T) {
	bin := bytes.Join([][]byte{
		[]byte("0123456789 01234 n"),
		[]byte{0x0d, 0x0a},
	}, nil)

	ret, err := newXrefEntry(newLines(bin)[0])
	if err != nil {
		t.Fatal(err)
	}
	if ret.Offset != 123456789 || ret.Generation != 1234 || !ret.InUse {
		t.Fatal(ret)
	}
}

func TestXrefEntryErrorCheckLength(t *testing.T) {
	bin := bytes.Join([][]byte{
		[]byte("0123456789 01234 n "),
		[]byte{0x0d, 0x0a},
	}, nil)

	if _, err := newXrefEntry(newLines(bin)[0]); errors.Cause(err) != errInvalidXref {
		t.Fatal(err)
	}
}

func TestXrefEntryErrorMatch(t *testing.T) {
	bin := bytes.Join([][]byte{
		[]byte("0123456789 01235 F"),
		[]byte{0x0d, 0x0a},
	}, nil)

	if _, err := newXrefEntry(newLines(bin)[0]); errors.Cause(err) != errInvalidXref {
		t.Fatal(err)
	}
}

func TestXref(t *testing.T) {
	bin := bytes.Join([][]byte{
		[]byte("xref"),
		[]byte("10 2"),
		[]byte("0123456789 01234 n"),
		[]byte("9876543210 43210 f"),
		[]byte("dummy"),
	}, []byte{0x0d, 0x0a})

	ret, left, err := newXref(newLines(bin))
	if err != nil {
		t.Fatal(err)
	}

	if len(ret) != 2 {
		t.Fatal(ret)
	}
	if ent := ret[10]; ent.Offset != 123456789 || ent.Generation != 1234 || !ent.InUse {
		t.Fatal(ent)
	}
	if ent := ret[11]; ent.Offset != 9876543210 || ent.Generation != 43210 || ent.InUse {
		t.Fatal(ent)
	}
	if string(left.Join()) != "dummy" {
		t.Fatal(left)
	}
}

func TestXrefErrorTooShortLines(t *testing.T) {
	bin := bytes.Join([][]byte{
		[]byte("xref"),
		// []byte("10 2"),
		// []byte("0123456789 01234 n"),
		// []byte("9876543210 43210 f"),
		// []byte("dummy"),
	}, []byte{0x0d, 0x0a})

	if _, _, err := newXref(newLines(bin)); errors.Cause(err) != errInvalidXref {
		t.Fatal(err)
	}
}

func TestXrefErrorCheckXref(t *testing.T) {
	bin := bytes.Join([][]byte{
		// []byte("xref"),
		[]byte("10 2"),
		[]byte("0123456789 01234 n"),
		[]byte("9876543210 43210 f"),
		[]byte("dummy"),
	}, []byte{0x0d, 0x0a})

	if _, _, err := newXref(newLines(bin)); errors.Cause(err) != errInvalidXref {
		t.Fatal(err)
	}
}

func TestXrefErrorMatchRange(t *testing.T) {
	bin := bytes.Join([][]byte{
		[]byte("xref"),
		// []byte("10 2"),
		[]byte("0123456789 01234 n"),
		[]byte("9876543210 43210 f"),
		[]byte("dummy"),
	}, []byte{0x0d, 0x0a})

	if _, _, err := newXref(newLines(bin)); errors.Cause(err) != errInvalidXref {
		t.Fatal(err)
	}
}

func TestXrefErrorEntry(t *testing.T) {
	bin := bytes.Join([][]byte{
		[]byte("xref"),
		[]byte("10 2"),
		[]byte("0123456789 01234 n"),
		// []byte("9876543210 43210 f"),
		[]byte("dummy"),
	}, []byte{0x0d, 0x0a})

	if _, _, err := newXref(newLines(bin)); errors.Cause(err) != errInvalidXref {
		t.Fatal(err)
	}
}
