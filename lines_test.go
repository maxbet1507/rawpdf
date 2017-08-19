package rawpdf

import (
	"encoding/hex"
	"testing"
)

func TestLines(t *testing.T) {
	buf, _ := hex.DecodeString("000102030405060708090a0b0c0d0e0f0a3020")
	ls := newLines(buf)

	if len(ls) != 3 {
		t.Fatal(ls)
	}
	if ls[0].Hex() != "000102030405060708090a" {
		t.Fatal(ls)
	}
	if ls[1].Hex() != "0b0c0d0e0f0a" {
		t.Fatal(ls)
	}
	if ls[2].Hex() != "3020" || ls[2].String() != "0" {
		t.Fatal(ls)
	}
	if hex.EncodeToString(ls.Join()) != hex.EncodeToString(buf) {
		t.Fatal(ls)
	}
}
