package objects

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/pkg/errors"
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
	if tkn := <-ret; tkn.Type != typeIdent || tkn.Value != "1234" {
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

func TestUnmarshal(t *testing.T) {
	src := bytes.Join([][]byte{
		[]byte("<<"),
		[]byte("%/Name [1234 true false null <00> (test)]"),
		[]byte{0x0a},
		[]byte("/Name [1234 true false null %comment"),
		[]byte{0x0a},
		[]byte("<00> (test\\\\dummy)]"),
		[]byte(">>"),
	}, nil)

	ret, err := Unmarshal(src)
	if err != nil {
		t.Fatal(err)
	}

	bin, _ := json.Marshal(ret)
	if string(bin) != `{"/Name":["1234",true,false,null,"\u003c00\u003e","(test\\\\dummy)"]}` {
		t.Fatal(ret)
	}
}

func TestUnmarshalErrorEnd(t *testing.T) {
	src := bytes.Join([][]byte{
		[]byte("<<"),
		[]byte("%/Name [1234 true false null <00> (test)]"),
		[]byte{0x0a},
		[]byte("/Name [1234 true false null %comment"),
		[]byte{0x0a},
		[]byte("<00> (test\\\\dummy)]"),
		// []byte(">>"),
	}, nil)

	if _, err := Unmarshal(src); errors.Cause(err) != errUnmarshalFailure {
		t.Fatal(err)
	}
}

func TestUnmarshalErrorDictionary(t *testing.T) {
	src := bytes.Join([][]byte{
		[]byte("<<"),
		[]byte("%/Name [1234 true false null <00> (test)]"),
		[]byte{0x0a},
		[]byte("/Name [1234 true false null %comment"),
		[]byte{0x0a},
		[]byte("<00> (test\\\\dummy)]"),
		[]byte("1234"),
		[]byte(">>"),
	}, nil)

	if _, err := Unmarshal(src); errors.Cause(err) != errUnmarshalFailure {
		t.Fatal(err)
	}
}

func TestUnmarshalErrorInvalidToken(t *testing.T) {
	src := bytes.Join([][]byte{
		[]byte("<<"),
		[]byte("%/Name [1234 true false null <00> (test)]"),
		[]byte{0x0a},
		[]byte("/Name [1234 true false null %comment"),
		[]byte{0x0a},
		[]byte("/("),
		[]byte("<00> (test\\\\dummy)]"),
		[]byte(">>"),
	}, nil)

	if _, err := Unmarshal(src); errors.Cause(err) != errUnmarshalFailure {
		t.Fatal(err)
	}
}
