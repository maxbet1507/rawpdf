package rawpdf

import (
	"bytes"
	"testing"

	"github.com/pkg/errors"
)

func TestMainTrailer(t *testing.T) {
	bin := bytes.Join([][]byte{
		[]byte("dummy"),
		[]byte("trailer"),
		[]byte("<</Size 1234"),
		[]byte("/Test true>>"),
		[]byte("startxref"),
		[]byte("5678"),
		[]byte("%%EOF"),
	}, []byte{0x0d, 0x0a})

	ret, err := newMainTrailer(bin)
	if err != nil {
		t.Fatal(err)
	}

	if ret.Start != 5678 {
		t.Fatal(ret)
	}
	if ret.Dictionary != "<</Size 1234\r\n/Test true>>\r\n" {
		t.Fatal(ret)
	}
}

func TestMainTrailerErrorTooShortLines(t *testing.T) {
	bin := bytes.Join([][]byte{
		// []byte("dummy"),
		// []byte("trailer"),
		// []byte("<</Size 1234"),
		// []byte("/Test true>>"),
		[]byte("startxref"),
		[]byte("5678"),
		[]byte("%%EOF"),
	}, []byte{0x0d, 0x0a})

	if _, err := newMainTrailer(bin); errors.Cause(err) != errInvalidTrailer {
		t.Fatal(err)
	}
}

func TestMainTrailerErrorCheckEOF(t *testing.T) {
	bin := bytes.Join([][]byte{
		[]byte("dummy"),
		[]byte("trailer"),
		[]byte("<</Size 1234"),
		[]byte("/Test true>>"),
		[]byte("startxref"),
		[]byte("5678"),
		// []byte("%%EOF"),
	}, []byte{0x0d, 0x0a})

	if _, err := newMainTrailer(bin); errors.Cause(err) != errInvalidTrailer {
		t.Fatal(err)
	}
}

func TestMainTrailerErrorConvertStartXref(t *testing.T) {
	bin := bytes.Join([][]byte{
		[]byte("dummy"),
		[]byte("trailer"),
		[]byte("<</Size 1234"),
		[]byte("/Test true>>"),
		[]byte("startxref"),
		// []byte("5678"),
		[]byte("%%EOF"),
	}, []byte{0x0d, 0x0a})

	if _, err := newMainTrailer(bin); errors.Cause(err) != errInvalidTrailer {
		t.Fatal(err)
	}
}

func TestMainTrailerErrorCheckStartXref(t *testing.T) {
	bin := bytes.Join([][]byte{
		[]byte("dummy"),
		[]byte("trailer"),
		[]byte("<</Size 1234"),
		[]byte("/Test true>>"),
		// []byte("startxref"),
		[]byte("5678"),
		[]byte("%%EOF"),
	}, []byte{0x0d, 0x0a})

	if _, err := newMainTrailer(bin); errors.Cause(err) != errInvalidTrailer {
		t.Fatal(err)
	}
}

func TestMainTrailerErrorCheckTrailer(t *testing.T) {
	bin := bytes.Join([][]byte{
		[]byte("dummy"),
		// []byte("trailer"),
		[]byte("<</Size 1234"),
		[]byte("/Test true>>"),
		[]byte("startxref"),
		[]byte("5678"),
		[]byte("%%EOF"),
	}, []byte{0x0d, 0x0a})

	if _, err := newMainTrailer(bin); errors.Cause(err) != errInvalidTrailer {
		t.Fatal(err)
	}
}

func TestFirstpageTrailer(t *testing.T) {
	bin := bytes.Join([][]byte{
		[]byte("trailer"),
		[]byte("<</Size 1234"),
		[]byte("/Test true>>"),
		[]byte("startxref"),
		[]byte("5678"),
		[]byte("%%EOF"),
		[]byte("dummy"),
	}, []byte{0x0d, 0x0a})

	ret, err := newFirstpageTrailer(bin)
	if err != nil {
		t.Fatal(err)
	}

	if ret.Start != 5678 {
		t.Fatal(ret)
	}
	if ret.Dictionary != "<</Size 1234\r\n/Test true>>\r\n" {
		t.Fatal(ret)
	}
}

func TestFirstpageTrailerTooShortLines(t *testing.T) {
	bin := bytes.Join([][]byte{
		[]byte("trailer"),
		// []byte("<</Size 1234"),
		// []byte("/Test true>>"),
		[]byte("startxref"),
		[]byte("5678"),
		// []byte("%%EOF"),
		// []byte("dummy"),
	}, []byte{0x0d, 0x0a})

	if _, err := newFirstpageTrailer(bin); errors.Cause(err) != errInvalidTrailer {
		t.Fatal(err)
	}
}

func TestFirstpageTrailerErrorCheckTrailer(t *testing.T) {
	bin := bytes.Join([][]byte{
		// []byte("trailer"),
		[]byte("<</Size 1234"),
		[]byte("/Test true>>"),
		[]byte("startxref"),
		[]byte("5678"),
		[]byte("%%EOF"),
		[]byte("dummy"),
	}, []byte{0x0d, 0x0a})

	if _, err := newFirstpageTrailer(bin); errors.Cause(err) != errInvalidTrailer {
		t.Fatal(err)
	}
}

func TestFirstpageTrailerErrorCheckStartXref(t *testing.T) {
	bin := bytes.Join([][]byte{
		[]byte("trailer"),
		[]byte("<</Size 1234"),
		[]byte("/Test true>>"),
		// []byte("startxref"),
		[]byte("5678"),
		[]byte("%%EOF"),
		[]byte("dummy"),
	}, []byte{0x0d, 0x0a})

	if _, err := newFirstpageTrailer(bin); errors.Cause(err) != errInvalidTrailer {
		t.Fatal(err)
	}
}

func TestFirstpageTrailerErrorConvertStartXref(t *testing.T) {
	bin := bytes.Join([][]byte{
		[]byte("trailer"),
		[]byte("<</Size 1234"),
		[]byte("/Test true>>"),
		[]byte("startxref"),
		// []byte("5678"),
		[]byte("%%EOF"),
		[]byte("dummy"),
	}, []byte{0x0d, 0x0a})

	if _, err := newFirstpageTrailer(bin); errors.Cause(err) != errInvalidTrailer {
		t.Fatal(err)
	}
}

func TestFirstpageTrailerErrorCheckEOF(t *testing.T) {
	bin := bytes.Join([][]byte{
		[]byte("trailer"),
		[]byte("<</Size 1234"),
		[]byte("/Test true>>"),
		[]byte("startxref"),
		[]byte("5678"),
		// []byte("%%EOF"),
		[]byte("dummy"),
	}, []byte{0x0d, 0x0a})

	if _, err := newFirstpageTrailer(bin); errors.Cause(err) != errInvalidTrailer {
		t.Fatal(err)
	}
}
