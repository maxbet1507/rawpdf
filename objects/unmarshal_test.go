package objects

import (
	"bytes"
	"testing"
)

func TestTokenizer(t *testing.T) {
	src := bytes.Join([][]byte{
		[]byte("<<"),
		[]byte("%/Name [1234 true false null <00> (test)]"),
		[]byte{0x0a},
		[]byte("/Name [1234 true false null <00> (test\\\\dummy)]"),
		[]byte(">>"),
	}, nil)

	ret := tokenizer(src)

	if tkn := <-ret; tkn.Type != typeDictionaryBegin {
		t.Fatal(tkn)
	}
	if tkn := <-ret; tkn.Type != typeComment || tkn.Value != "%/Name [1234 true false null <00> (test)]" {
		t.Fatal(tkn)
	}
	if tkn := <-ret; tkn.Type != typeName || tkn.Value != "/Name" {
		t.Fatal(tkn)
	}
	if tkn := <-ret; tkn.Type != typeArrayBegin {
		t.Fatal(tkn)
	}
	if tkn := <-ret; tkn.Type != typeUnknown || tkn.Value != "1234" {
		t.Fatal(tkn)
	}
	if tkn := <-ret; tkn.Type != typeTrue {
		t.Fatal(tkn)
	}
	if tkn := <-ret; tkn.Type != typeFalse {
		t.Fatal(tkn)
	}
	if tkn := <-ret; tkn.Type != typeNull {
		t.Fatal(tkn)
	}
	if tkn := <-ret; tkn.Type != typeHexdecimalString || tkn.Value != "<00>" {
		t.Fatal(tkn)
	}
	if tkn := <-ret; tkn.Type != typeLiteralString || tkn.Value != "(test\\\\dummy)" {
		t.Fatal(tkn)
	}
	if tkn := <-ret; tkn.Type != typeArrayEnd {
		t.Fatal(tkn)
	}
	if tkn := <-ret; tkn.Type != typeDictionaryEnd {
		t.Fatal(tkn)
	}
	if tkn, ok := <-ret; ok {
		t.Fatal(tkn)
	}
}

func TestTokenizerErrorQuo(t *testing.T) {
	src := bytes.Join([][]byte{
		[]byte("/("),
	}, nil)

	ret := tokenizer(src)

	if tkn := <-ret; tkn.Type != typeInvalid || tkn.Value != "/(" {
		t.Fatal(tkn)
	}
}

func TestTokenizerErrorLss(t *testing.T) {
	src := bytes.Join([][]byte{
		[]byte("<00"),
	}, nil)

	ret := tokenizer(src)

	if tkn := <-ret; tkn.Type != typeInvalid || tkn.Value != "<00" {
		t.Fatal(tkn)
	}
}

func TestTokenizerErrorGtr(t *testing.T) {
	src := bytes.Join([][]byte{
		[]byte(">00"),
	}, nil)

	ret := tokenizer(src)

	if tkn := <-ret; tkn.Type != typeInvalid || tkn.Value != ">00" {
		t.Fatal(tkn)
	}
}
