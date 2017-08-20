package objects

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

var (
	whitespaces = map[byte]struct{}{
		0x00: struct{}{}, // null
		0x09: struct{}{}, // tab
		0x0a: struct{}{}, // LF
		0x0c: struct{}{}, // FF
		0x0d: struct{}{}, // CR
		0x20: struct{}{}, // space
	}
	delimiters = map[byte]struct{}{
		// '(': struct{}{},
		// ')': struct{}{},
		'<': struct{}{},
		'>': struct{}{},
		'[': struct{}{},
		']': struct{}{},
		// '{': struct{}{},
		// '}': struct{}{},
		'/': struct{}{},
		// '%': struct{}{},
	}

	errUnmarshalFailure = fmt.Errorf("unmarshal failure")
)

type tokenType int

const (
	typeInvalid tokenType = iota
	typeIdent

	typeComment
	typeLiteralString
	typeHexdecimalString

	typeDictionaryBegin
	typeDictionaryEnd

	typeArrayBegin
	typeArrayEnd

	typeName
	typeTrue
	typeFalse
	typeNull
)

type token struct {
	Value string
	Type  tokenType
}

func tokenizer(v []byte) <-chan token {
	cch := make(chan byte)
	go func() {
		for _, c := range v {
			cch <- c
		}
		close(cch)
	}()

	wch := make(chan token)
	go func() {
		buf := []byte{}

		flush := func(t tokenType) {
			if len(buf) > 0 {
				wch <- token{Value: strings.TrimRight(string(buf), "\r"), Type: t}
				buf = []byte{}
			}
		}

		tokenREM := func() {
			for c := range cch {
				if c == 0x0a {
					break
				}
				buf = append(buf, c)
			}
			flush(typeComment)
		}

		tokenLPAREN := func() {
			for c := range cch {
				buf = append(buf, c)
				if c == ')' {
					flush(typeLiteralString)
					return
				}
				if c == '\\' {
					if c, ok := <-cch; ok {
						buf = append(buf, c)
					}
				}
			}
			flush(typeInvalid)
		}

		for c := range cch {
			if _, ok := whitespaces[c]; ok {
				flush(typeIdent)
				continue
			}

			if _, ok := delimiters[c]; ok {
				flush(typeIdent)
				wch <- token{Value: string(c), Type: typeIdent}
				continue
			}

			switch c {
			case '%':
				flush(typeIdent)
				buf = append(buf, c)
				tokenREM()
			case '(':
				flush(typeIdent)
				buf = append(buf, c)
				tokenLPAREN()
			default:
				buf = append(buf, c)
			}
		}
		flush(typeIdent)

		close(wch)
	}()

	tch := make(chan token)
	go func() {
		tokenGTR := func(buf []string) {
			if c1, ok := <-wch; ok {
				buf = append(buf, c1.Value)
				if c1.Type == typeIdent && c1.Value == ">" {
					tch <- token{Value: strings.Join(buf, ""), Type: typeDictionaryEnd}
					return
				}
			}

			tch <- token{Value: strings.Join(buf, ""), Type: typeInvalid}
		}

		tokenLSS := func(buf []string) {
			if c1, ok := <-wch; ok {
				buf = append(buf, c1.Value)
				if c1.Type == typeIdent && c1.Value == "<" {
					tch <- token{Value: strings.Join(buf, ""), Type: typeDictionaryBegin}
					return
				}

				if _, err := hex.DecodeString(c1.Value); c1.Type == typeIdent && err == nil {
					if c2, ok := <-wch; ok {
						buf = append(buf, c2.Value)
						if c2.Type == typeIdent && c2.Value == ">" {
							tch <- token{Value: strings.Join(buf, ""), Type: typeHexdecimalString}
							return
						}
					}
				}
			}

			tch <- token{Value: strings.Join(buf, ""), Type: typeInvalid}
		}

		tokenQUO := func(buf []string) {
			if c1, ok := <-wch; ok {
				buf = append(buf, c1.Value)
				if c1.Type == typeIdent {
					tch <- token{Value: strings.Join(buf, ""), Type: typeName} // name
					return
				}
			}

			tch <- token{Value: strings.Join(buf, ""), Type: typeInvalid}
		}

		buf := []string{}

		flushbuf := func() {
			if val := strings.Join(buf, " "); len(val) > 0 {
				tch <- token{Value: val, Type: typeIdent}
			}
			buf = []string{}
		}

		for w := range wch {
			if w.Type != typeIdent {
				flushbuf()
				tch <- w
				continue
			}

			switch w.Value {
			case "<":
				flushbuf()
				tokenLSS([]string{w.Value})

			case ">":
				flushbuf()
				tokenGTR([]string{w.Value})

			case "/":
				flushbuf()
				tokenQUO([]string{w.Value})

			case "[":
				flushbuf()
				tch <- token{Value: w.Value, Type: typeArrayBegin}

			case "]":
				flushbuf()
				tch <- token{Value: w.Value, Type: typeArrayEnd}

			case "true":
				flushbuf()
				tch <- token{Value: w.Value, Type: typeTrue}

			case "false":
				flushbuf()
				tch <- token{Value: w.Value, Type: typeFalse}

			case "null":
				flushbuf()
				tch <- token{Value: w.Value, Type: typeNull}

			default:
				buf = append(buf, w.Value)
			}
		}
		flushbuf()
		close(tch)
	}()

	return tch
}

func toObject(tn <-chan token) (interface{}, error) {
	if t, ok := <-tn; ok {
		switch t.Type {
		case typeComment:
			return toObject(tn)

		case typeTrue:
			return true, nil

		case typeFalse:
			return false, nil

		case typeNull:
			return nil, nil

		case typeName, typeLiteralString, typeHexdecimalString, typeIdent:
			return t.Value, nil

		case typeDictionaryBegin:
			ret := map[string]interface{}{}

			for t := range tn {
				switch t.Type {
				case typeComment:
				case typeName:
					val, err := toObject(tn)
					if err != nil {
						return val, err
					}
					ret[t.Value] = val

				case typeDictionaryEnd:
					return ret, nil

				default:
					return nil, errors.Wrap(errUnmarshalFailure, "dictionary")
				}
			}

		case typeArrayBegin:
			ret := []interface{}{}

			for {
				val, err := toObject(tn)
				if err != nil || val == typeArrayEnd {
					return ret, err
				}
				ret = append(ret, val)
			}

		case typeArrayEnd:
			return typeArrayEnd, nil

		default:
			return nil, errors.Wrap(errUnmarshalFailure, "invalid token")
		}
	}

	return nil, errors.Wrap(errUnmarshalFailure, "end")
}

// Unmarshal -
func Unmarshal(v []byte) (interface{}, error) {
	return toObject(tokenizer(v))
}
