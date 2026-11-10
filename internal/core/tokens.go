package core

import "fmt"

type TokenKind int

// All accepted characters are ASCII (less )

const (
	TokenNone             TokenKind = iota
	TokenNumber                     // [1-9][0-9]*
	TokenIdentifier                 // [a-zA-Z][a-zA-Z._]*  // Numbers are not allowed on identifier names to make lexing simpler
	TokenColon                      // ','
	TokenParenLeft                  // '('
	TokenParenRight                 // ')'
	TokenOpMultiplication           // '*'
	TokenOpDivision                 // '//'
	TokenOpModulo                   // '%'
	TokenOpPlus                     // '+'
	TokenOpMinus                    // '-'
	TokenOpBitLeft                  // '<<'
	TokenOpBitRight                 // '>>'
	TokenOpEqual                    // '=='
	TokenOpNotEqual                 // '!='
	TokenOpGreater                  // '>'
	TokenOpLesser                   // '<'
	TokenOpGreaterEqual             // '>='
	TokenOpLesserEqual              // '<='
	TokenOpAnd                      // '&&'
	TokenOpOr                       // '||'
	TokenEnd
)

func (t TokenKind) String() string {
	switch t {
	case TokenNone:
		return "TokenNone"
	case TokenNumber:
		return "NUMBER"
	case TokenIdentifier:
		return "IDENTIFIER"
	case TokenColon:
		return ","
	case TokenParenLeft:
		return "("
	case TokenParenRight:
		return ")"
	case TokenOpMultiplication:
		return "*"
	case TokenOpDivision:
		return "//"
	case TokenOpModulo:
		return "%"
	case TokenOpPlus:
		return "+"
	case TokenOpMinus:
		return "-"
	case TokenOpBitLeft:
		return "<<"
	case TokenOpBitRight:
		return ">>"
	case TokenOpEqual:
		return "=="
	case TokenOpNotEqual:
		return "!="
	case TokenOpGreater:
		return ">"
	case TokenOpLesser:
		return "<"
	case TokenOpGreaterEqual:
		return ">="
	case TokenOpLesserEqual:
		return "<="
	case TokenOpAnd:
		return "&&"
	case TokenOpOr:
		return "||"
	case TokenEnd:
		return "TokenEnd"

	default:
		panic(t)
	}
}

type Token struct {
	Kind       TokenKind
	Pos        int
	Len        int
	Number     int64
	Identifier string
}

func (t Token) String() string {
	switch t.Kind {
	case TokenNumber:
		return fmt.Sprintf("%s(%d)", t.Kind.String(), t.Number)
	case TokenIdentifier:
		return fmt.Sprintf("%s(%s)", t.Kind.String(), t.Identifier)
	default:
		return t.Kind.String()
	}
}
