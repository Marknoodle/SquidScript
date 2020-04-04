// First we define the datatype for all of our tokens: String
// We chose string instead of some complex type because it is simple
struct TokenType(String);

//Now we define the values for all of our tokens
//Special Characters:
const ILLEGAL: String = "ILLEGAL";
const EOF: String = "EOF";

//Identifiers and Literals
const IDENTIFIER: String = "IDENTIFIER"; //Variable names
const INT: String = "INT"; //integers

//Operators
//math
const ASSIGN: String = ":=";
const PLUS: String = "+";
const MINUS: String = "-";
const BANG: String = "!";
const STAR: String = "*";
const SLASH: String = "/";

//comparators
const LT: String = "<";
const GT: String = ">";
const EQ: String = "=";
const NEQ: String = "!=";

//delimeters
const COMMA: String = ",";
const SEMICOLON: String = ";";
const LPAREN: String = "(";
const RPAREN: String = ")";
const LCURLYBRACKET: String = "{";
const RCURLYBRACKET: String = "}";
const LSQUAREBRACKET: String = "[";
const RSQUAREBRACKET: String = "]";

//Keywords
const FUNCTION: String = "FUNCTION";
const LET: String= "LET";
const TRUE: String "TRUE";
const FALSE: String = "FALSE";
const RETURN: String = "RETURN";

//Now we define the a structure for the tokens themselves
struct Token 
{
    
}