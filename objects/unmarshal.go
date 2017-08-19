package objects

import (
	"encoding/hex"
	"strings"
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
)

type tokenType int

const (
	typeInvalid tokenType = iota
	typeUnknown

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
				flush(typeUnknown)
				continue
			}

			if _, ok := delimiters[c]; ok {
				flush(typeUnknown)
				wch <- token{Value: string(c), Type: typeUnknown}
				continue
			}

			switch c {
			case '%':
				flush(typeUnknown)
				buf = append(buf, c)
				tokenREM()
			case '(':
				flush(typeUnknown)
				buf = append(buf, c)
				tokenLPAREN()
			default:
				buf = append(buf, c)
			}
		}
		flush(typeUnknown)

		close(wch)
	}()

	tch := make(chan token)
	go func() {
		tokenGTR := func(buf []string) {
			if c1, ok := <-wch; ok {
				buf = append(buf, c1.Value)
				if c1.Type == typeUnknown && c1.Value == ">" {
					tch <- token{Value: strings.Join(buf, ""), Type: typeDictionaryEnd}
					return
				}
			}

			tch <- token{Value: strings.Join(buf, ""), Type: typeInvalid}
		}

		tokenLSS := func(buf []string) {
			if c1, ok := <-wch; ok {
				buf = append(buf, c1.Value)
				if c1.Type == typeUnknown && c1.Value == "<" {
					tch <- token{Value: strings.Join(buf, ""), Type: typeDictionaryBegin}
					return
				}

				if _, err := hex.DecodeString(c1.Value); c1.Type == typeUnknown && err == nil {
					if c2, ok := <-wch; ok {
						buf = append(buf, c2.Value)
						if c2.Type == typeUnknown && c2.Value == ">" {
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
				if c1.Type == typeUnknown {
					tch <- token{Value: strings.Join(buf, ""), Type: typeName} // name
					return
				}
			}

			tch <- token{Value: strings.Join(buf, ""), Type: typeInvalid}
		}

		for w := range wch {
			if w.Type != typeUnknown {
				tch <- w
				continue
			}

			switch w.Value {
			case "<":
				tokenLSS([]string{w.Value})

			case ">":
				tokenGTR([]string{w.Value})

			case "/":
				tokenQUO([]string{w.Value})

			case "[":
				tch <- token{Value: w.Value, Type: typeArrayBegin}

			case "]":
				tch <- token{Value: w.Value, Type: typeArrayEnd}

			case "true":
				tch <- token{Value: w.Value, Type: typeTrue}

			case "false":
				tch <- token{Value: w.Value, Type: typeFalse}

			case "null":
				tch <- token{Value: w.Value, Type: typeNull}

			default:

				tch <- w
			}
		}
		close(tch)
	}()

	return tch
}

// Unmarshal -
// func Unmarshal(v []byte) (interface{}, error) {

// 	for t := range tokenizer(v) {
// 		fmt.Println(t)
// 	}

// 	return nil, nil
// }
