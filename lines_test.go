package rawpdf

import "testing"
import "encoding/hex"

func TestLines(t *testing.T) {
	buf, _ := hex.DecodeString("000102030405060708090a0b0c0d0e0f0a20")
	ls := newLines(buf)

	if len(ls) != 3 {
		t.Fatal(ls)
	}
	if hex.EncodeToString(ls[0]) != "000102030405060708090a" {
		t.Fatal(ls)
	}
	if hex.EncodeToString(ls[1]) != "0b0c0d0e0f0a" {
		t.Fatal(ls)
	}
	if hex.EncodeToString(ls[2]) != "20" || ls[2].String() != " " {
		t.Fatal(ls)
	}
	if hex.EncodeToString(ls.Join()) != hex.EncodeToString(buf) {
		t.Fatal(ls)
	}
}
