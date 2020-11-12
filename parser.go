package sqltools

import (
	"io"
)

type SelectStatement struct {
	Fields []string
	TableName string
}

type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	// Literals
	IDENT // fields, table_name

	// Misc characters
	ASTERISK // *
	COMMA    // ,

	// Keywords
	SELECT
	FROM
)

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

type Parser struct {
  s *Scanner
}

func (p Parser) Parse() (interface{}, interface{}) {

}

func NewParser(reader io.Reader)*Parser  {
 return &Parser{
 	s: NewScanner(reader),
 }
}

