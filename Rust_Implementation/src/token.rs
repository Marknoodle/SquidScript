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
const RETURNS: String = "RETURNS";
const KEY: String = "KEY";
const LOCK: String = "LOCK";

//Now we define the a structure for the tokens themselves
struct Token 
{
    Type : TokenType,
    Literal : String,
}

let mut keywords = HashMap::new();

keywords.insert(String::from("fn"), FUNCTION);
keywords.insert(String::from("let"), LET);
keywords.insert(String::from("true"), TRUE);
keywords.insert(String::from("false"), FALSE);
keywords.insert(String::from("return"), RETURN);
keywords.insert(String::from("returns"), RETURNS);
keywords.insert(String::from("key"), KEY);
keywords.insert(String::from("lock"), LOCK);