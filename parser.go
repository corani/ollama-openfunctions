package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Call struct {
	Tool  string         `json:"tool"`
	Input map[string]any `json:"tool_input"`
}

// TODO(daniel): this is all very quick and dirty. Clean this up with a proper
// tokenizer and better error handling.
func ParseCalls(input string) ([]Call, error) {
	p := &parser{
		lexer: &lexer{
			input:  []byte(input),
			cursor: 0,
		},
	}

	if err := p.ParseCalls(); err != nil {
		return nil, err
	}

	return p.calls, nil
}

type parser struct {
	lexer *lexer
	calls []Call
	call  Call
}

func (p *parser) ParseCalls() error {
	p.lexer.Trim()

	for p.lexer.Header() {
		p.lexer.Trim()

		if err := p.ParseCall(); err != nil {
			return err
		}
	}

	return nil
}

func (p *parser) ParseCall() error {
	p.call.Tool = p.lexer.Ident()
	p.call.Input = make(map[string]any)

	if p.call.Tool == "" {
		return fmt.Errorf("expected identifier at %v", p.lexer.Tail())
	}

	if !p.lexer.Expect('(') {
		return fmt.Errorf("expected lparen at %v", p.lexer.Tail())
	}

	if err := p.ParseArgs(); err != nil {
		return err
	}

	if !p.lexer.Expect(')') {
		return fmt.Errorf("expected rparen at %v", p.lexer.Tail())
	}

	p.calls = append(p.calls, p.call)

	return nil
}

func (p *parser) ParseArgs() error {
	if p.lexer.Peek() == ')' {
		return nil
	}

	for {
		p.lexer.Trim()

		if err := p.ParseArg(); err != nil {
			return err
		}

		if !p.lexer.Expect(',') {
			break
		}
	}

	return nil
}

func (p *parser) ParseArg() error {
	key := p.lexer.Ident()
	if key == "" {
		return fmt.Errorf("expected identifier at %v", p.lexer.Tail())
	}

	if !p.lexer.Expect('=') {
		return fmt.Errorf("expected '=' at %v", p.lexer.Tail())
	}

	var val any

	if v, ok := p.lexer.String(); ok {
		val = strings.ReplaceAll(v, "\\'", "'")
	} else if v, ok := p.lexer.Int(); ok {
		val = v
	}

	p.call.Input[key] = val

	return nil
}

type lexer struct {
	input  []byte
	cursor int
}

func (l *lexer) Tail() string {
	return string(l.input[l.cursor:])
}

func (l *lexer) Header() bool {
	if bytes.HasPrefix(l.input[l.cursor:], []byte("<function>")) {
		l.cursor += len("<function>")

		return true
	} else if bytes.HasPrefix(l.input[l.cursor:], []byte("<<function>>")) {
		l.cursor += len("<<function>>")

		return true
	}

	return false
}

func (l *lexer) Trim() {
	if l.cursor >= len(l.input) {
		return
	}

	for l.input[l.cursor] == ' ' {
		l.cursor++
	}
}

func (l *lexer) Ident() string {
	start := l.cursor

	for {
		ch := l.input[l.cursor]

		if ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch >= '0' && ch <= '9' {
			l.cursor++
		} else {
			break
		}
	}

	if l.cursor > start {
		return string(l.input[start:l.cursor])
	}

	return ""
}

func (l *lexer) String() (string, bool) {
	var delimiter byte

	if l.Expect('\'') {
		delimiter = '\''
	} else if l.Expect('"') {
		delimiter = '"'
	} else {
		return "", false
	}

	start := l.cursor

	for {
		ch := l.input[l.cursor]

		if ch == delimiter {
			break
		} else if ch == '\\' {
			l.cursor++
		}

		l.cursor++
	}

	if !l.Expect(delimiter) {
		return "", false
	}

	if l.cursor > start {
		return string(l.input[start : l.cursor-1]), true
	}

	return "", false
}

func (l *lexer) Int() (int, bool) {
	start := l.cursor

	for {
		ch := l.input[l.cursor]

		if ch >= '0' && ch <= '9' {
			l.cursor++
		} else {
			break
		}
	}

	i, err := strconv.ParseInt(string(l.input[start:l.cursor]), 10, 64)
	if err != nil {
		return 0, false
	}

	return int(i), true
}

func (l *lexer) Expect(ch byte) bool {
	if l.input[l.cursor] == ch {
		l.cursor++

		return true
	}

	return false
}

func (l *lexer) Peek() byte {
	return l.input[l.cursor]
}
