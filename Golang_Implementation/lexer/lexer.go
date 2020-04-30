//OVERVIEW: Lexer is going to read characters and evaluate what they mean. Lexers give meaning to keywords, operations, and whitespace by tokenizing them as it reads code.

package lexer

import "monkey/token"

//This is what is constructed; the main template/structure of the lexer
type Lexer struct {
	input        string // code provided by user
	position     int    // current position in input (points to current char)
	readPosition int    // current reading position in input (after current char)
	ch           byte   // current char under examination
} //end Lexer struct

//REQUIRES: a string input
//MODIFIES:
//EFFECTS: creates lexer structure
func New(src string) *Lexer { //serves as a lexer constructor
	l := &Lexer{input: src} //the lexer structure l recieves src (source code) as input
	l.readChar()            //reads the first character of the input and adjusts position and read position accordingly
	return l                //returns the lexer
} //end constructor

//REQUIRES: a lexer for input (previous method)
//MODIFIES:
//EFFECTS: tokenizes ch (current char under examination)
func (l *Lexer) NextToken() token.Token {
	var tok token.Token // a token variable of the struct defined above

	l.skipWhitespace() //the function reads characters as long as they are whitespace until it has skipped all the whitespace between two other characters

	switch l.ch { //switch statement that identifies the current character in l and then tokenizes the character based on what it is and what char's surround it
	case '=':
		if l.peekChar() == '=' { //if the next character is a second =, meaning that the two == make a boolean operator
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else { //the next character is not another = and therefore l.ch is an assignment = and not a boolean ==
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' { //if the input is != then it is tokenized as not equals
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else { //the next character is not a = and therefore the ! should be read as simply ! (not)
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case 0: //reached end of input so we need to create a EOF (end of file) token
		tok.Literal = ""
		tok.Type = token.EOF
	default: //the character is not a special character, assignment operator,or boolean operator, so we need to check if it is a letter (check to see if it is a keyword), a digit, or if it is illegal
		if isLetter(l.ch) { //is ch a letter
			tok.Literal = l.readIdentifier()          //ch is a letter so we need to read until the next space and determine if Literal is a keyword
			tok.Type = token.LookupIdent(tok.Literal) //returns the appropraite keyword type if Literal is a keyword, otherwise returns IDENT token type
			return tok                                // returns tok which contains Literal and token type
		} else if isDigit(l.ch) { //is ch a number
			tok.Type = token.INT         //the character is an integer, so its token is assigned token type int
			tok.Literal = l.readNumber() //the literal becomes the entire number (until the next whitespace)
			return tok                   // returns tok which contains Literal and token Type
		} else { //the character is not a digit nor is it a letter, therefore it is some illegal character the language will not support
			tok = newToken(token.ILLEGAL, l.ch)
		}
	} //end cases

	l.readChar()
	return tok
} //end NextToken

//%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%% LEXER HELPER METHODS

//REQUIRES: a lexer structure l
//MODIFIES: position will become the location of the next non-whitespace char in input. readPosition will become the location proceeding position's location
//EFFECTS: the position; calls readChar() and reads the next posititon/ skips to next position as long as the current character is whitespace
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' { //keep reading char's until we have hit nonwhitespace
		l.readChar()
	} //end for
} //end skipWhiteSpace

//REQUIRES: a lexer structure l
//MODIFIES: the current position and readPosition are incremented
//EFFECTS: reads the next character in the input
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) { //the read position is outside the range of the input (we have reached the end of the input) so we need to make l.ch = 0 so that an EOF token can be made
		l.ch = 0 //we return 0 because we've reached the outside of our input
	} else { //the read position is INSIDE the range of our input
		l.ch = l.input[l.readPosition] //we read the next char and assign its value to l.ch
	}
	l.position = l.readPosition //position becomes the location of the current char
	l.readPosition += 1         //we increment readPosition so that the next character is ready to be read
} //end readChar

//REQUIRES: a lexer structure l
//MODIFIES:
//EFFECTS: returns the current char at readPosition of the inputted lexer l, or 0 (for an EOF token)
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) { //the current readPosition is outside the input/ has reached the end of the input
		return 0 //returns 0 because we are outside the range of the input
	} else { //readposition is less than the length of the input/ still inside the range
		return l.input[l.readPosition] //returns the char at the readPosition of the input
	}
} //end peekChar

//REQUIRES:a lexer structure l
//MODIFIES: changes position and readPosition to relect the end of the identifier/label
//EFFECTS: returns a string of chars representing an identifier/label. Ex: returns 'apple' in "let apple = 530;"
func (l *Lexer) readIdentifier() string {
	position := l.position //start of identifier
	for isLetter(l.ch) {   //read all the letters of the identifier
		l.readChar()
	} //end for
	return l.input[position:l.position] // return string of identifier so that it can become the string literal for that identifier's token
} //end readIdentifier

//REQUIRES: a lexer structure l
//MODIFIES: changes position and readPosition to relect the end of the number
//EFFECTS: returns a string of chars representing an number. Ex: returns '530' in "let apple = 530;"
func (l *Lexer) readNumber() string {
	position := l.position //start of number
	for isDigit(l.ch) {    //read all the digits of the number
		l.readChar()
	} //end for
	return l.input[position:l.position] // return string of full number so that it can become the string literal for that number's token
} //end readNumber

//REQUIRES: a char of the input to be examined
//MODIFIES:
//EFFECTS: returns a bool of whether or not the char is a letter
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
} //end isLetter

//REQUIRES: a char of the input to be examined
//MODIFIES:
//EFFECTS: returns a bool of whether or not the char is an number
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
} //end isDigit

//REQUIRES: a tokentype and character input
//MODIFIES:
//EFFECTS: creates token of inputted token type and string literal containing inputted char
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
} //end newToken
