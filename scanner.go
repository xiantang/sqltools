package sqltools

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode"
)


var eof = rune(0)

// Scanner represents a lexical scanner.
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// read reads the next rune from the bufferred reader.
// Returns the rune(-1) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() { _ = s.r.UnreadRune() }



func (s *Scanner) Scan() (tok Token,lit string) {
	ch := s.read()
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) || ch == '`' {
		s.unread()
		return s.scanIdent()
	}else if isDigit(ch) {
		s.unread()
		return s.scanDigit()
	}

	switch ch {
	case eof:
		return EOF, ""
	case '*':
		return ASTERISK, string(ch)
	case ',':
		return COMMA, string(ch)
	case '(':
		return LeftParentheses,string(ch)
	case ')':
		return RightParentheses,string(ch)
	}
	return ILLEGAL, string(ch)
}

func (s *Scanner)scanWhitespace() (Token, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())
	for {
		if ch := s.read(); ch == eof {
			break
		}else if !isWhitespace(ch) {
			s.unread()
			break
		}else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

func (s *Scanner) scanDigit() (Token, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())
	for {
		if ch := s.read(); ch == eof {
			break
		}else if !isDigit(ch) {
			s.unread()
			break
		}else {
			buf.WriteRune(ch)
		}
	}
	return NUMBER, buf.String()

}

func (s *Scanner) scanIdent() (Token, string) {
	var buf bytes.Buffer
	if ch := s.read(); isLetter(ch) {
		buf.WriteRune(ch)
	}
	for {
		if ch := s.read(); ch == eof {
			break
		}else if ch == '`' {
			break
		} else if !isLetter(ch)&& !isDigit(ch) && ch != '_' {
			s.unread()
			break
		}else {
			buf.WriteRune(ch)
		}
	}

	switch strings.ToUpper(buf.String()) {
	case "SELECT":
		return SELECT, buf.String()
	case "FROM":
		return FROM, buf.String()
	case "ALTER":
		return ALTER, buf.String()
	case "TABLE":
		return TABLE, buf.String()
	case "COLUMN":
		return COLUMN, buf.String()
	case "DROP":
		return DROP, buf.String()
	case "ADD":
		return ADD, buf.String()
	case "VARCHAR":
		return VARCHAR, buf.String()
	case "NULL":
		return NULL, buf.String()
	}
	return IDENT, buf.String()
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}