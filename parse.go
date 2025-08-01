package caca

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	lex "github.com/temaxuck/caca/internal/lexer"
)

type parser struct {
	l    *lex.Lexer
	cvs  *Canvas
	buf  []uint8
	path string
}

// TODO: This is dirty. Come up with a better interface.
func ParseCanvasFile(path string) (*Canvas, error) {
	input, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	p := parser{
		l:    &lex.Lexer{Input: string(input)},
		cvs:  NewCanvas(),
		path: path,
	}

	for p.l.NextToken() {
		switch p.l.Token.Kind {
		case lex.TK_Comment:
			p.parseMetadataEntry()
			break
		case lex.TK_Digit:
			err := p.parsePixel()
			if err != nil {
				return nil, err
			}
			break
		case lex.TK_NewLine:
			p.appendAndResetBuf()
		}

	}

	switch p.l.Token.Kind {
	case lex.TK_Invalid:
		line, col := p.l.GetCursorFilePos()
		return nil, fmt.Errorf("unexpected token %q at (%d, %d)", p.l.Token.Str, line, col)
	case lex.TK_EOF:
		p.appendAndResetBuf()
	}

	return p.cvs, nil
}

func (p *parser) parsePixel() error {
	px, err := strconv.ParseUint(p.l.Token.Str, 10, 8)
	if err != nil {
		return err
	}
	if px > 4 {
		return fmt.Errorf("intensity level should be in range [0..4], but got: %d", px)
	}

	p.buf = append(p.buf, uint8(px))
	return nil
}

func (p *parser) appendAndResetBuf() {
	if len(p.buf) == 0 {
		return
	}

	row := make([]uint8, len(p.buf))
	copy(row, p.buf)

	p.cvs.Canvas2D = append(p.cvs.Canvas2D, row)
	p.buf = p.buf[:0]
}

func (p *parser) parseMetadataEntry() error {
	text, _ := strings.CutPrefix(p.l.Token.Str, "#")
	text = strings.TrimSpace(text)
	k, v, ok := strings.Cut(text, ":")
	if !ok {
		// This is just a comment
		return nil
	}
	k, v = strings.TrimSpace(k), strings.TrimSpace(v)
	switch k {
	case "Start date":
		d, err := parseDate(v)
		if err != nil {
			return err
		}
		p.cvs.SetStartDate(d)
	case "Repository":
		r, err := parseRepository(v, p.path)
		if err != nil {
			return err
		}
		p.cvs.SetRepository(r)
	case "Author":
		name, email, err := parseAuthor(v)
		if err != nil {
			return err
		}
		p.cvs.SetAuthor(name, email)
	}

	return nil
}

func parseDate(s string) (time.Time, error) {
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return time.Time{}, err
	}

	return t.Add(time.Hour * 12), err
}

func parseRepository(s, path string) (string, error) {
	if len(s) == 0 {
		return "", fmt.Errorf("repository must be non-empty string: local path or URL")
	}

	if !filepath.IsAbs(s) {
		s = filepath.Join(filepath.Dir(path), s)
	}

	return s, nil
}

func parseAuthor(s string) (string, string, error) {
	if len(s) == 0 {
		return "", "", fmt.Errorf("author must follow format: <username> <email>")
	}

	tokens := strings.Split(s, " ")
	if len(tokens) == 1 {
		return "", "", fmt.Errorf("author must follow format: <username> <email>")
	}

	user, email := strings.Join(tokens[:len(tokens)-1], " "), tokens[len(tokens)-1]
	if !validateEmail(email) {
		return "", "", fmt.Errorf("invalid email format")
	}
	return user, email, nil
}
