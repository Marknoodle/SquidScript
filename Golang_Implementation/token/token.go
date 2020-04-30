//OVERVIEW: Tokenizer gives value to certain words and characters, defines keywords, operators, special characters, etc. 
package token

type TokenType string

const (//these are our token types
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT = "IDENT" // user given input that does not match any pre-defined operators or keywords, (ie variable labels, strings, etc)
	INT   = "INT"   // 1343456

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{//this is a hashmap where inputted text may match a keyword, thus requiring the token thereof to have the appropriate keyword token type
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

//REQUIRES: a string input (the string literal of the identifier we are trying to tokenize)
//MODIFIES: 
//EFFECTS: returns token type of passed in identifier (either a keyword token type or IDENT)
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {//does the input match a keyword in our hashmap (represented by the bool ok). If so, tok holds the appropriate keyword token type we wish to return
		return tok//returns tok, which has the tokenized value of the input, and ok is true
	}
	return IDENT //returns TokenType IDENT, since ok evaluated as false because passed in input does not match a keyword in our hashmap
}