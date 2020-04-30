package parser

import (
	"fmt"
	"../ast"
	"../lexer"
	"../token"
	"strconv" //For when we need to obtain the actual int value of numbers inputted in source code
)

const ( //these constants are precedence labels that make the hierarchy of expression parsing order
	_           int = iota // we use iota to give the following constants incrementing numbers (1 - 7) as values. The order of these values decides the order in which expression get parse
	LOWEST                 //							VALUE 1, LOWEST PRECEDENCE
	EQUALS                 // ==						VALUE 2
	LESSGREATER            // > or <					VALUE 3
	SUM                    // +							VALUE 4
	PRODUCT                // *							VALUE 5
	PREFIX                 // -X or !X					VALUE 6
	CALL                   // myFunction(X)				VALUE 7, HIGHEST PRECEDENCE
)

//in what order do we want to parse expressions so the AST is correct (Omit?)
// associates token types with their precedence
//EX: 5 * 5 + 10 The AST should represent this expression like this  ( (5 * 5) + 10 ) so that it is evaluated in proper order
var precedences = map[token.TokenType]int{ //Depending on token type provided, the appropriate precedence *number* is provided (actual numeric value is based upon consts defined above)
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
} //table can tell us that + (token.PLUS) and - (token.MINUS) have the same precedence, but are lower than the precedence of * (token.ASTERISK) and / (token.SLASH), for example

// Whenever a token type is encountered, the parsing functions are called to parse the appropriate expression and return an AST node that represents it
//Each token type can have up to two parsing functions associated with it, depending on whether the token is found in a prefix or an infix position
//prefixParseFns gets called when we encounter the associated token type in prefix position and infixParseFn gets called when we encounter the token type in infix position
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression // argument is “left side” of the infix operator that is being parsed
)

//This is what is constructed; the main template/structure of the parser
type Parser struct {
	l      *lexer.Lexer //pointer to an instance of the lexer, on which we repeatedly call NextToken() to get the next token in the input
	errors []string     //Contains errors we have seen (messages)

	curToken  token.Token //Current token
	peekToken token.Token //Next token

	//In order for our parser to get the correct prefixParseFn or infixParseFn for the current token type, we add two maps to the Parser structure
	//With these maps in place,we can just check if the appropriate map(infix or prefix)has a parsing function associated with curToken.Type
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser { //serves as a parser constructor
	p := &Parser{ //see 'type Parser struct {'
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn) // Prefix Parse functions. Parses based on token type seen in prefix position
	p.registerPrefix(token.IDENT, p.parseIdentifier)           // indentifier
	p.registerPrefix(token.INT, p.parseIntegerLiteral)         // int
	p.registerPrefix(token.BANG, p.parsePrefixExpression)      // not operator
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)     // negative sign
	p.registerPrefix(token.TRUE, p.parseBoolean)               // true bool
	p.registerPrefix(token.FALSE, p.parseBoolean)              // false bool
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)   // open parantheses (
	p.registerPrefix(token.IF, p.parseIfExpression)            // if
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)   // function

	p.infixParseFns = make(map[token.TokenType]infixParseFn) // Infix Parse functions. Parses based on token type seen in infix position
	//Every infix operator gets associated with the same parsing function called parseInfixExpression
	p.registerInfix(token.PLUS, p.parseInfixExpression)     // +
	p.registerInfix(token.MINUS, p.parseInfixExpression)    // -
	p.registerInfix(token.SLASH, p.parseInfixExpression)    // divide "/""
	p.registerInfix(token.ASTERISK, p.parseInfixExpression) // multiply "*"
	p.registerInfix(token.EQ, p.parseInfixExpression)       // == in m, = in ss
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)   // !=
	p.registerInfix(token.LT, p.parseInfixExpression)       // <
	p.registerInfix(token.GT, p.parseInfixExpression)       // >

	p.registerInfix(token.LPAREN, p.parseCallExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken      //current token becomes next
	p.peekToken = p.l.NextToken() //peek token becomes the one after that
}

func (p *Parser) curTokenIs(t token.TokenType) bool { //Is the CURRENT token type what we expect it to be
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool { //Is the NEXT token (peekToken) type what we expect it to be
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	//enforce the correctness of the order of tokens by checking the type of the next token
	//checks the type of the peekToken and only if the type is correct does it advance the tokens by calling nextToken
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else { // automatically adds an error every time one of our expectations about the next token was wrong
		p.peekError(t)
		return false
	}
}

func (p *Parser) Errors() []string { //we can check if the parser encountered any errors. (USED PRIMARILY FOR TESTING)
	return p.errors //returns list of errors
}

func (p *Parser) peekError(t token.TokenType) {
	//used to add an error to errors field of parser struct when the type of peekToken doesn’t match the expectation
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg) //add error message to errors field
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) { //just adds a formatted error message to our parser’s errors field when something is misused as a prefix parse function
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}              //construct the root node of the AST
	program.Statements = []ast.Statement{} //list that we will add parsed statements to

	//Fills the statements list of ast program struct with parsed statements
	for !p.curTokenIs(token.EOF) { //iterates over every token in the input until it encounters an token.EOF token
		stmt := p.parseStatement() //parses the line of code
		if stmt != nil {           //making sure that the statement we are adding is valid within syntax (stmt will be nil if statement is not valid)
			program.Statements = append(program.Statements, stmt) //return value is added to Statements slice of the AST root node |AND/OR| adding parsed statements to the list
		}
		p.nextToken() // advances both p.curToken and p.peekToken to the next statement
	}

	return program //When nothing is left to parse the *ast.Program root node is returned.
}

func (p *Parser) parseStatement() ast.Statement { // Deciding how to parse a statment based upon the token type that lets us know what kind of statement we are looking at
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement() //let statement
	case token.RETURN:
		return p.parseReturnStatement() //return statement
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement { //constructs an *ast.LetStatement node with the token it’s currently sitting on (a token.LET token) and then advances the tokens while making assertions about the next token with calls to expectPeek
	//let <identifier> := <expression>;
	stmt := &ast.LetStatement{Token: p.curToken} //let statement struct in AST obtains the let token

	if !p.expectPeek(token.IDENT) { //we expect to see a identifier/label to have some value assigned to it
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal} //constructs an *ast.Identifier node (ie variable label)

	if !p.expectPeek(token.ASSIGN) { //we expect to see assignment operator
		return nil
	}

	p.nextToken() //advancing tokens to expression (this is what comes after the assignment operator)

	stmt.Value = p.parseExpression(LOWEST) //value that we put in <identifier>

	if p.peekTokenIs(token.SEMICOLON) { //checking to see that the let statement has ended and advancing the cur and peek tokens
		p.nextToken()
	}

	return stmt // send back the statement so that it can be added to program's statement list (if it is valid)
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement { // constructs a ast.ReturnStatement
	//return <expression>;
	stmt := &ast.ReturnStatement{Token: p.curToken} //return statement struct in AST obtains the return token

	p.nextToken() //advancing tokens to expression (what comes after 'return')

	stmt.ReturnValue = p.parseExpression(LOWEST) //value that we are returning

	if p.peekTokenIs(token.SEMICOLON) { //checking to see that the return statement has ended and advancing the cur and peek tokens
		p.nextToken()
	}

	return stmt // send back the statement so that it can be added to program's statement list (if it is valid)
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement { // constructs a ast.ExpressionStatement
	stmt := &ast.ExpressionStatement{Token: p.curToken} //expression statement struct in AST obtains the current token

	stmt.Expression = p.parseExpression(LOWEST) //we pass the lowest possible precedence to parseExpression, since we didn’t parse anything yet and we can’t compare precedences

	if p.peekTokenIs(token.SEMICOLON) { //checking to see that the expression statement has ended and advancing the cur and peek tokens
		p.nextToken()
	}

	return stmt
}

//%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%LAST LEFT OFF

//Determines which parsing function (if any) should parse the given expression based off of token type seen
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type] //checkes whether we have a parsing function associated with p.curToken.Type in the prefix position
	if prefix == nil {                          //if we don't
		p.noPrefixParseFnError(p.curToken.Type) //add a error message to our parser’s errors field
		return nil
	}
	leftExp := prefix() //If we do, it calls that parsing function and stores that parsed expression in leftExp

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

//returns the precedence associated with the token type of p.peekToken. If it doesn’t find a precedence for p.peekToken it defaults to LOWEST, the lowest possible precedence any operator can have
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

//does the same thing as peekPrecedence, but for p.curToken.
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

//%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%% Below are the parsing functions are registered in the parser constructor

func (p *Parser) parseIdentifier() ast.Expression { //returns a *ast.Identifier node with the current token in the Token field and the literal value of the token in Value field
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression { //returns a *ast.IntegerLiteral node with the current token in the Token field and the literal value of the token in Value field
	lit := &ast.IntegerLiteral{Token: p.curToken} // IntegerLiteral node obtains int token

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64) //turning the token literal(a string) into a int variable called value.
	if err != nil {                                           //if the inputted token literal could not be parsed into an int, err != nil (meaning an error had occured)
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg) //adding error to parser struct's errors list
		return nil
	}

	lit.Value = value //placing the int-parsed value of curToken.Literal into ast.IntegerLiteral node's field

	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression { // token seen is prefix operator "-" or "!"
	//<prefix operator><expression>;
	expression := &ast.PrefixExpression{ //creates prefix expression node with the current token and its literal
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken() //looking at what we are applying the prefix operator to (ie 'boolVar' in '!boolVar' or 'randInt' in '-randInt')

	//parseExpression will now return a newly constructed node and parsePrefixExpression uses it to fill the Right field of *ast.PrefixExpression
	expression.Right = p.parseExpression(PREFIX) //parses expression following the prefix operator

	return expression //returns the whole parsed expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{ //fills infixExpression node in ast (except Right)
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence() //assigns the precedence of the current token (which is the operator of the infixexpression)to precedence local var
	p.nextToken()                   //advances token so we can fill the Right field of the node
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

//below are two helper methods for the parser that add entries to the prefixParseFns and infixParseFns maps

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
