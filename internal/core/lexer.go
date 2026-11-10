package core

import (
	"fmt"
	"strconv"
	"unicode"
)

type lexerState struct {
	content       []rune
	cursor        int // Next rune to read
	line          int
	lastLineStart int
}

func newLexer(content []rune) lexerState {
	return lexerState{content: content}
}

func Lex(content string) ([]Token, error) {

	state := newLexer([]rune(content))
	var tokens []Token

	for {
		token, err := state.next()
		if err != nil {
			return nil, fmt.Errorf("Error on input read: %s", err.Error())
		}
		if token.Kind == TokenEnd {
			break
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

func (state *lexerState) next() (Token, error) {
	var token Token

	state.trimLeft()

	if state.cursor >= len(state.content) {
		return Token{Kind: TokenEnd, Pos: state.cursor}, nil
	}

	if isCommentStart(state.content[state.cursor]) {
		state.consumeCommentBody()
		state.trimLeft()
	}

	if state.cursor >= len(state.content) {
		return Token{Kind: TokenEnd, Pos: state.cursor}, nil
	}

	if isIdentifierStart(state.content[state.cursor]) {
		token.Pos = state.cursor
		state.cursor++
		state.consumeIdentifierBody()
		token.Len = state.cursor - token.Pos
		token.Identifier = string(state.content[token.Pos:state.cursor])
		token.Kind = TokenIdentifier
		return token, nil
	}

	// Rejects numbers that start with 0 but also rejects "0".
	if isNumberStart(state.content[state.cursor]) {
		token.Pos = state.cursor
		state.cursor++
		state.consumeNumberBody()
		token.Len = state.cursor - token.Pos
		numberStr := string(state.content[token.Pos:state.cursor])
		number, err := strconv.ParseInt(numberStr, 10, 64)
		if err != nil {
			return token, fmt.Errorf("Error while parsing a number: %w", err)
		}
		token.Number = number
		token.Kind = TokenNumber
		return token, nil
	}

	// Accepts bare "0"
	if state.content[state.cursor] == '0' {
		if state.cursor+1 >= len(state.content) || !isNumberBody(state.content[state.cursor+1]) {
			token.Pos = state.cursor
			state.cursor++
			token.Kind = TokenNumber
			token.Number = 0
			token.Len = 1
			return token, nil
		}
	}

	// Fixed time cost, a hand-rolled implementation could be faster
	for _, literal := range LiteralTokensPrecedence {
		token := state.conditionalConsume(literal.TokenKind, literal.string)
		if token.Kind != TokenNone {
			return token, nil
		}
	}

	return token, state.unmanagedError()
}

// Must be sorted by precedence
// IMPORTANT: Make sure those patterns contained in other ones go last i.e. ">=" and ">"
var LiteralTokensPrecedence []struct {
	TokenKind
	string
} = []struct {
	TokenKind
	string
}{
	{TokenOpDivision, "//"},
	{TokenOpBitLeft, "<<"},
	{TokenOpBitRight, ">>"},
	{TokenOpGreaterEqual, ">="},
	{TokenOpLesserEqual, "<="},
	{TokenOpEqual, "=="},
	{TokenOpNotEqual, "!="},
	{TokenOpAnd, "&&"},
	{TokenOpOr, "||"},
	{TokenParenLeft, "("},
	{TokenParenRight, ")"},
	{TokenOpMultiplication, "*"},
	{TokenOpModulo, "%"},
	{TokenOpPlus, "+"},
	{TokenOpMinus, "-"},
	{TokenOpGreater, ">"},
	{TokenOpLesser, "<"},
	{TokenColon, ","},
}

func (state *lexerState) conditionalConsume(kind TokenKind, literal string) Token {

	literalRunes := []rune(literal)

	if len(literalRunes)+state.cursor > len(state.content) {
		return Token{Kind: TokenNone}
	}

	lineJumps := 0
	lastLineStart := state.lastLineStart
	for i, patternRune := range literalRunes {
		r := state.content[state.cursor+i]
		if r != patternRune {
			return Token{Kind: TokenNone}
		}
		if r == '\n' {
			lineJumps++
			lastLineStart = state.cursor + i
		}
	}

	t := Token{Kind: kind, Pos: state.cursor, Len: len(literal)}

	state.cursor += len(literal)
	state.line += lineJumps
	state.lastLineStart = lastLineStart

	return t
}

func (state *lexerState) consumeCommentBody() {
	for state.cursor < len(state.content) {
		r := state.content[state.cursor]
		state.cursor++
		if r == '\n' {
			state.line++
			state.lastLineStart = state.cursor
			break
		}
	}
}

func (state *lexerState) trimLeft() {
	for state.cursor < len(state.content) && unicode.IsSpace(state.content[state.cursor]) {
		r := state.content[state.cursor]
		state.cursor++
		if r == '\n' {
			state.line++
			state.lastLineStart = state.cursor
		}
	}
}

func (state *lexerState) consumeIdentifierBody() {
	for state.cursor < len(state.content) && isIdentifierBody(state.content[state.cursor]) {
		state.cursor++
	}
}
func (state *lexerState) consumeNumberBody() {
	for state.cursor < len(state.content) && isNumberBody(state.content[state.cursor]) {
		state.cursor++
	}
}

// ====================================

func isIdentifierStart(r rune) bool {
	return unicode.IsLetter(r)
}

func isIdentifierBody(r rune) bool {
	return unicode.IsLetter(r) || r == '.' || r == '_'
}

func isNumberStart(r rune) bool {
	return unicode.IsDigit(r) && r != '0'
}

func isNumberBody(r rune) bool {
	return unicode.IsDigit(r)
}

func isCommentStart(r rune) bool {
	return r == '#'
}

func (state lexerState) positionString() string {
	return fmt.Sprintf("position %d (line=%d, row=%d)", state.cursor, state.line, state.cursor-state.lastLineStart)
}

// Logic for specific error messages
func (state lexerState) unmanagedError() error {

	if state.cursor >= len(state.content) {
		return fmt.Errorf("Unmanaged EOF at %s", state.positionString())
	}

	if state.content[state.cursor] == '0' {
		return fmt.Errorf("Numbers cannot start with a zero (0): %s", state.positionString())
	}

	return fmt.Errorf("Invalid token at %s: '%c'", state.positionString(), state.content[state.cursor])
}

func IsValidIdentifier(name string) bool {
	identifier := []rune(name)
	if len(identifier) <= 0 {
		return false
	}

	if !isIdentifierStart(identifier[0]) {
		return false
	}

	for _, r := range identifier[1:] {
		if !isIdentifierBody(r) {
			return false
		}
	}

	return true
}
