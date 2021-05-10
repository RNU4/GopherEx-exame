package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

//Token bla bla bla
type Token struct {
	identifiers []string
	tokens      []Token
}

const (
	ERROR_OUT_OF_BOUNDS = iota + 1
	ERROR_WHITESPACE_IN_NAME
	ERROR_UNIDENTIFIED_TYPE
	ERROR_UNIDENTIFIED_OPERATOR
	ERROR_WRONG_ARGUMENT_FORMAT
	ERROR_FUNCTION_ALREADY_EXSITS
	ERROR_FUNCTION_DOESNT_EXIST
	ERROR_WRONG_TYPE_FORMAT
)

//Tokenizer bla bla bla
type Tokenizer struct {
	code      []string
	tokens    []Token
	functions map[string]Token
	index     int
	typeBytes map[string]bool
	html      string
	styleMap  map[string]string
}

func (tokenizer Tokenizer) peak() string {
	if len(tokenizer.code) > tokenizer.index+1 {
		return tokenizer.code[tokenizer.index+1]
	}
	return ""
}

func (tokenizer *Tokenizer) next() string {
	if len(tokenizer.code) > tokenizer.index+1 {
		tokenizer.index++
		return tokenizer.code[tokenizer.index] //Return token
	}
	tokenizer.error(ERROR_OUT_OF_BOUNDS)
	return ""
}

func (tokenizer *Tokenizer) tokenizeCode() Token {
	//Is function
	if len(tokenizer.code)-1 < tokenizer.index {
		return Token{}
	}

	if tokenizer.isFunction() {
		//tokenize function
		functionToken := tokenizer.tokenizeFunction()
		tokenizer.next()
		tokenizer.next()
		//go to next line and Tokenize the body of the function
		for true {
			if tokenizer.code[tokenizer.index] == "]" {
				break
			}
			tok := tokenizer.tokenizeCode()
			//Check if token is empty (quick fix for functions)
			if len(tok.identifiers) > 0 || len(tok.tokens) > 0 {
				functionToken.tokens = append(functionToken.tokens, tok)
			}

		}
		if len(tokenizer.code)-1 > tokenizer.index {
			tokenizer.next()
		}

		tokenizer.functions[functionToken.identifiers[0]] = functionToken

		//return functionToken
	} else if tokenizer.isTypeByte() {
		token := tokenizer.tokenizeType()
		tokenizer.next()
		return token

	} else if tokenizer.isFunctionCall() {
		token := tokenizer.tokenizeFunctionCall()
		tokenizer.next()
		return token
	} else {
		tokenizer.error(ERROR_UNIDENTIFIED_TYPE)
	}

	return Token{}
}

//Type Check
func (tokenizer *Tokenizer) isFunction() bool {
	return tokenizer.peak() == "["
}
func (tokenizer *Tokenizer) isFunctionCall() bool {
	functionName, _ := tokenizer.getFunctionIdentifiers()
	//Check if function exist
	_, ok := tokenizer.functions[functionName]
	if ok {
		return ok //tokenizer.functions[functionName]
	}
	//Error if function does not exist
	tokenizer.error(ERROR_FUNCTION_DOESNT_EXIST)
	return false
}
func (tokenizer *Tokenizer) isTypeByte() bool {

	currentLine := tokenizer.code[tokenizer.index]
	if len(currentLine) < 2 {
		tokenizer.error(ERROR_WRONG_TYPE_FORMAT)
	}
	typeByte := currentLine[0]
	//Checks if the the length of the string is greater than 1 and validates that it's a single byte
	validatorByte := (map[bool]bool{true: (currentLine[1] == ' '), false: false})[len(currentLine) > 1]
	//Checks if it's a valid typeByte
	if tokenizer.typeBytes[string(typeByte)] != false && validatorByte == true {
		return true
	}
	return false
}

//Tokenizers
func (tokenizer *Tokenizer) tokenizeFunction() Token {
	functionName, UnformatedArguments := tokenizer.getFunctionIdentifiers()
	//Check if function allready exists

	if tokenizer.functionNameExists(functionName) {
		tokenizer.error(ERROR_FUNCTION_ALREADY_EXSITS) //Exits with error
	}
	tokenizer.functions[functionName] = Token{}
	//appends function to function list
	//tokenizer.functions[functionName] = true
	argumentNames := strings.Split(strings.Replace(UnformatedArguments, " ", "", -1), ",") //...
	//Clean argument names + Check for whitespaces in argument names and whitespaces in function name
	argumentNames = tokenizer.cleanArgumentNames(argumentNames)
	tokenizer.whitespaceInName(functionName)
	//Creates the function identifiers
	var indentifiers []string
	indentifiers = append(indentifiers, functionName)
	//First element is the function name
	indentifiers = append(indentifiers, argumentNames...)
	//Function arguments [arg1,arg2...]
	token := Token{indentifiers, make([]Token, 0)}
	return token
}

func (tokenizer *Tokenizer) tokenizeType() Token {
	currentLine := tokenizer.code[tokenizer.index]
	typeByte := currentLine[0]
	var indentifiers []string
	switch typeByte {
	case 'T':
		indentifiers = append(indentifiers, string(typeByte)) //Append type
		indentifiers = append(indentifiers, currentLine[2:])  //currentLine[2:] =T |hello
		break

	default:
		indentifiers = append(indentifiers, string(typeByte))                       //Append type
		indentifiers = append(indentifiers, strings.Split(currentLine[2:], " ")...) //Append Arguments
	}
	token := Token{indentifiers, make([]Token, 0)}
	return token
}

func (tokenizer *Tokenizer) tokenizeFunctionCall() Token {
	functionName, UnformatedArguments := tokenizer.getFunctionIdentifiers()
	//Check if function exist and if it dosent error out
	// '\"' ]
	//ARRAY
	//TEMPSTR
	//","
	if !tokenizer.functionNameExists(functionName) {
		tokenizer.error(ERROR_FUNCTION_DOESNT_EXIST)
	}
	var arguments []string
	arguments = append(arguments, functionName)
	argument := ""
	for _, char := range UnformatedArguments {
		if char == ',' {
			arguments = append(arguments, argument)
			argument = ""
		} else {
			argument = argument + string(char)
		}
	}
	if argument != "" {
		arguments = append(arguments, argument)
	}

	token := Token{arguments, make([]Token, 0)}
	return token
}

//Tokenize

// Helper functions
func (tokenizer *Tokenizer) error(errorCode int) {
	switch errorCode {
	case ERROR_OUT_OF_BOUNDS:
		fmt.Println("ERROR_OUT_OF_BOUNDS:", "INDEX", tokenizer.index, "OUT OF", len(tokenizer.code)-1)
		os.Exit(errorCode)
		break

	case ERROR_UNIDENTIFIED_TYPE:
		fmt.Println("ERROR_UNIDENTIFIED_TYPE:", "AT", tokenizer.code[tokenizer.index])
		os.Exit(errorCode)
		break
	case ERROR_WHITESPACE_IN_NAME:
		fmt.Println("ERROR_WHITESPACE_IN_NAME:", "AT", tokenizer.code[tokenizer.index])
		os.Exit(errorCode)
		break
	case ERROR_WRONG_ARGUMENT_FORMAT:
		fmt.Println("ERROR_WRONG_ARGUMENT_FORMAT:", "AT", tokenizer.code[tokenizer.index])
		os.Exit(errorCode)
		break
	case ERROR_FUNCTION_ALREADY_EXSITS:
		fmt.Println("ERROR_FUNCTION_ALREADY_EXSITS:", "AT", tokenizer.code[tokenizer.index])
		os.Exit(errorCode)
		break

	case ERROR_FUNCTION_DOESNT_EXIST:
		fmt.Println("ERROR_FUNCTION_DOESNT_EXIST:", "AT", tokenizer.code[tokenizer.index])
		os.Exit(errorCode)
		break
	case ERROR_WRONG_TYPE_FORMAT:
		fmt.Println("ERROR_WRONG_TYPE_FORMAT:", "AT", tokenizer.code[tokenizer.index])
		os.Exit(errorCode)
		break

	default:
		fmt.Println("ERROR_UNKOWN_ERROR:", "AT", tokenizer.code[tokenizer.index])
		os.Exit(errorCode)
		break
	}
}

func (tokenizer *Tokenizer) whitespaceInName(name string) {
	isWhiteSpace := func(c rune) bool {
		return c == '\t' || c == '\r' || c == ' '
	}
	if strings.IndexFunc(name, isWhiteSpace) != -1 {
		tokenizer.error(ERROR_WHITESPACE_IN_NAME)
	}

}

func (tokenizer *Tokenizer) functionNameExists(functionName string) bool {
	_, ok := tokenizer.functions[functionName]
	return ok //tokenizer.functions[functionName]

}

func (tokenizer *Tokenizer) cleanArgumentNames(names []string) []string {
	var cleanNames []string
	for _, name := range names {
		cleanNames = append(cleanNames, strings.TrimSpace(name))
		tokenizer.whitespaceInName(cleanNames[len(cleanNames)-1]) //Stops tokinization if true and throws error
	}
	return cleanNames
}

func (tokenizer *Tokenizer) getFunctionIdentifiers() (string, string) {
	currentLine := tokenizer.code[tokenizer.index]
	UnformatedIdentifiers := strings.Split(currentLine, "(") //box(arg1,arg2) =[box, argstring]
	functionName := strings.TrimSpace(UnformatedIdentifiers[0])
	arguments := strings.Split(UnformatedIdentifiers[len(UnformatedIdentifiers)-1], ")")
	return functionName, arguments[0]
}
func (tokenizer *Tokenizer) getFunction(functionName string) Token {
	if tokenizer.functionNameExists(functionName) {
		tokenizer.error(ERROR_FUNCTION_DOESNT_EXIST)
	}

	return tokenizer.functions[functionName]
}
func main() {
	b, err := ioutil.ReadFile("file.txt") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	sourceCode := string(b)
	cleanCode := sanitize(sourceCode)
	parsedCode := parseCode(cleanCode)
	//!MODIFED
	tokenizer := Tokenizer{parsedCode, make([]Token, 0), map[string]Token{}, 0, map[string]bool{"I": true,"T": true, "M": true, "S": true, "V": true}, "", make(map[string]string)} 
	tokens := make([]Token, 0)
	for tokenizer.code[tokenizer.index] != "." {
		tok := tokenizer.tokenizeCode()

		//Check if token is empty (quick fix for functions)
		if len(tok.identifiers) > 0 || len(tok.tokens) > 0 {
			tokens = append(tokens, tok)
		}
	}
	fmt.Println("------------MAIN CODE------------ ")
	printTokens(tokens, 0)
	fmt.Println("------------FUNCTIONS------------ ")
	for k, v := range tokenizer.functions {
		fmt.Println()
		fmt.Println("Function :", k)
		tmpAr := make([]Token, 0)
		tmpAr = append(tmpAr, v)
		printTokens(tmpAr, 0)
	}

	fmt.Println("------------PROGRAM START------------ ")
	tokenizer.generatePrint(tokens, make(map[string]int))

	fmt.Println("------------HTML------------ ")
	fmt.Println(tokenizer.html)
	d1 := []byte(tokenizer.html)
	ioutil.WriteFile("index.html", d1, 0644)
	fmt.Println("------------TEST------------ ")
	/*
		ar := make([]string, 0)
		ar = append(ar, "align")
		ar = append(ar, "left")
		fmt.Println(tokenizer.style(ar))*/

}

func sanitize(str string) []string {
	tempCode := strings.Split(str, "\n")

	for index := range tempCode {
		tempCode[index] = strings.Replace(tempCode[index], string(rune(13)), "", -1) // REMOVE 13/CR

		if tempCode[index] != "" {
			i := 0
			for string(tempCode[index][i]) == " " {
				i++
				if i >= len(tempCode[index]) {
					break
				}
			}

			tempCode[index] = tempCode[index][i:]
		}
	}

	var array []string
	for _, _string := range tempCode {
		if _string != "" {
			array = append(array, _string)
		}

	}
	return array
}

func parseCode(code []string) []string {
	breaks := make(map[string]int)
	breaks["["] = 1 //"()"
	breaks["]"] = 1
	tokenString := ""
	var tempCode []string
	for _, _string := range code {
		for _, char := range _string { //"a,"
			if breaks[string(char)] == 1 {
				if tokenString != "" {
					tempCode = append(tempCode, tokenString)
				}
				tempCode = append(tempCode, string(char))
				tokenString = ""
			} else {
				tokenString += string(char)
			}
		}
		if tokenString != "" {
			tempCode = append(tempCode, tokenString)
			tokenString = ""
		}
	}
	if tokenString != "" {
		tempCode = append(tempCode, tokenString)
	}
	tempCode = append(tempCode, ".") //EOF
	return tempCode
}

func tokenizeCode(code []string) []Token {

	tokenizer := Tokenizer{code, make([]Token, 0), map[string]Token{}, -1, map[string]bool{}, "", make(map[string]string)}
	tokens := make([]Token, 0)
	tempToken := Token{make([]string, 0), make([]Token, 0)}

	for tokenizer.index < len(code)-1 {
		element := tokenizer.next()

		if tokenizer.peak() == "[" {
			brackLeft := 0
			brackRight := -1
			i := 2
			for brackLeft != brackRight {
				if code[i+tokenizer.index] == "[" {
					brackLeft++
				}
				if code[i+tokenizer.index] == "]" {
					brackRight++
				}
				i++
			}
			tempToken.identifiers = strings.Split(element, " ")
			tempToken.tokens = tokenizeCode(code[tokenizer.index+2 : i+tokenizer.index]) //CALL FUNCTION

			tokenizer.index += i              //SKIP BEYOND FUNCTION
			if tokenizer.index >= len(code) { //NO MORE TOKENS RETURN LAST
				tempToken.identifiers = nil
				tempToken.identifiers = append(tempToken.identifiers, strings.Split(element, " ")...)
				tokens = append(tokens, tempToken)
				return tokens
			}

			if string(code[tokenizer.index][0]) == "(" { //I

				tempToken.identifiers = append(tempToken.identifiers, append(append(make([]string, 0), "()"), strings.Split(code[tokenizer.index][1:len(code[tokenizer.index])-1], ",")...)...) //REMOVE ( and ), AND SPLIT STRING WITH ","
			}
			tokens = append(tokens, tempToken)
			tempToken.identifiers = nil
			tokenizer.index-- //GO ONE BACK FOR NEXT LOOP
		} else {
			//MAKE SURE THAT IT IS IN A FUNCTION
			if element != "]" && element != "[" && string(element[0]) != "(" {
				if len(tempToken.identifiers) > 0 {
					if tempToken.identifiers[0] == "T" || tempToken.identifiers[0] == "I"{ // T FOR TEXT //!modifed
						tokens = append(tokens, Token{append(make([]string, 0), element), make([]Token, 0)}) //TEXT DONT SPLIT INTO ARRAY
					} else {
						tokens = append(tokens, Token{strings.Split(element, " "), make([]Token, 0)}) //NOT TEXT SPLIT INTO ARRAY
					}
				} else {
					if strings.Split(element, " ")[0] == "T" || strings.Split(element, " ")[0] == "I"{ //!modifed
						tokens = append(tokens, Token{append(make([]string, 0), element), make([]Token, 0)}) //TEXT DONT SPLIT INTO ARRAY
					} else {
						tokens = append(tokens, Token{strings.Split(element, " "), make([]Token, 0)}) //IF NOT IN FUNCTION SPLIT INTO ARRAY}
					}
				}

			}
		}
	}

	return tokens
}

func printTokens(tokens []Token, liftOff int) {
	for _, element := range tokens {
		fmt.Println("├" + strings.Repeat("─", liftOff*4) + " " + strings.Join(element.identifiers, " "))
		if len(element.tokens) > 0 { //Check if token has sub tokens
			printTokens(element.tokens, liftOff+1)
		}
	}
}

func printCode(code []string) {
	for _, _string := range code {
		fmt.Println(_string)
	}
}
func joinMap(smap map[string]string) string {

	str := ""
	if len(smap) > 0 {
		str = str + " style=\""
	} else {
		return str
	}
	for key, value := range smap {
		str = str + key + ":" + value + ";"
	}
	str = str + "\""
	return str
}
func (tokenizer *Tokenizer) style(styles []string) string {
	if styles[0] == "clear" {
		tokenizer.styleMap = make(map[string]string)
		return ""
	}
	args := 2
	for i := 0; i < len(styles)-1; i += args {
		switch styles[i] {
		case "align":
			if styles[i+1] == "text" {
				if styles[i+2] == "left" {
					tokenizer.styleMap["text-align"] = "left"
				}

				if styles[i+2] == "right" {
					tokenizer.styleMap["text-align"] = "right"
				}

				if styles[i+2] == "center" {
					tokenizer.styleMap["text-align"] = "center"
				}
			}

			if styles[i+1] == "box" {
				if styles[i+2] == "center" {
					tokenizer.styleMap["margin"] = "0 auto"
				}
			}

			args = 3
			break

		case "color":
			tokenizer.styleMap["color"] = styles[i+1]
			args = 2
			break

		case "size":
			tokenizer.styleMap["size"] = styles[i+1]
			args = 2
			break

		case "padding":
			if styles[i+1] == "left" {
				tokenizer.styleMap["padding-left"] = styles[i+2]
			}

			if styles[i+1] == "right" {
				tokenizer.styleMap["padding-right"] = styles[i+2]
			}

			if styles[i+1] == "top" {
				tokenizer.styleMap["padding-top"] = styles[i+2]
			}
			if styles[i+1] == "bottom" {
				tokenizer.styleMap["padding-bottom"] = styles[i+2]
			}
			args = 3
			break

		case "margin":
			tokenizer.styleMap["margin"] = styles[i+1]
			if styles[i+1] == "left" {
				tokenizer.styleMap["margin-left"] = styles[i+2]
			}

			if styles[i+1] == "right" {
				tokenizer.styleMap["margin-right"] = styles[i+2]
			}

			if styles[i+1] == "top" {
				tokenizer.styleMap["margin-top"] = styles[i+2]
			}
			if styles[i+1] == "bottom" {
				tokenizer.styleMap["margin-bottom"] = styles[i+2]
			}
			args = 3
			break

		case "box":
			if styles[i+1] == "start" {
				tokenizer.html = tokenizer.html + "<div" + joinMap(tokenizer.styleMap) + ">"
			}
			if styles[i+1] == "end" {
				tokenizer.html = tokenizer.html + "</div>"
			}
			args = 2
			break
		case "border":
			if styles[i+1] == "radius" {
				tokenizer.styleMap["border-radius"] = styles[i+2]
			}

			if styles[i+1] == "style" {
				tokenizer.styleMap["border-style"] = styles[i+2]
			}

			if styles[i+1] == "width" {
				tokenizer.styleMap["border-width"] = styles[i+2]
			}
			args = 3
			break

		case "width":
			tokenizer.styleMap["width"] = styles[i+1]
			args = 2
			break
		case "height":
			tokenizer.styleMap["height"] = styles[i+1]
			args = 2
			break
		case "background":
			if styles[i+1] == "color" {
				tokenizer.styleMap["background-color"] = styles[i+2]
			}
			args = 3
			break
		case "font":
			if styles[i+1] == "weight" {
				tokenizer.styleMap["font-weight"] = styles[i+2]
			}
			if styles[i+1] == "size" {
				tokenizer.styleMap["font-size"] = styles[i+2]
			}
			args=3
		break
		case "image":
			tokenizer.styleMap["src"] = styles[i+1]
			args = 2
		break
		case "float":
			tokenizer.styleMap["float"] = styles[i+1]
			args = 2
		break
		}

	}
	return joinMap(tokenizer.styleMap)
}
func (tokenizer *Tokenizer) generatePrint(tokens []Token, varibels map[string]int) {

	tokenizer.html = tokenizer.html + "<link href=\"style.css\" rel=\"stylesheet\" type=\"text/css\">"
	for _, element := range tokens {
		switch element.identifiers[0] {
		case "T":
			fmt.Println(element.identifiers[1])
			tokenizer.html = tokenizer.html + "<p" + joinMap(tokenizer.styleMap) + ">" + element.identifiers[1] + "</p>"
			break

		case "I":
			//fmt.Println(element.identifiers[1])
			tokenizer.html = tokenizer.html + "<img src=\""+element.identifiers[1]+"\""+joinMap(tokenizer.styleMap) + "></img>"
			break

		case "V":
			//Printsvarible
			for _, variabelName := range element.identifiers[1:] {
				fmt.Println(varibels[variabelName])
				tokenizer.html = tokenizer.html + "<p" + joinMap(tokenizer.styleMap) + ">" + strconv.Itoa(varibels[variabelName]) + "</p>"
			}
			break
		case "M":
			result, command := evaluateMath(element.identifiers[1:], varibels)

			if command == "SET" {
				varibels[element.identifiers[1]] = result
			}
			if command == "EXIT" {
				return
			}
			break
		case "S":
			tokenizer.style(element.identifiers[1:])
			break
		}
		token, isFunction := tokenizer.functions[element.identifiers[0]]
		//call function and handle return
		if isFunction && len(element.tokens) == 0 {
			parsedVaribels := make(map[string]int)
			if len(element.identifiers) > 1 {
				for index, variabelName := range element.identifiers[1:] {
					number := 0
					if varibels[variabelName] != 0 {
						number = varibels[variabelName]
					} else {
						parsedNumber, err := strconv.Atoi(variabelName)
						if err != nil {
							number = 0
						} else {
							number = parsedNumber
						}
					}
					parsedVaribels[token.identifiers[index+1]] = number
				}
			}

			tmpTok := make([]Token, 0)
			tmpTok = append(tmpTok, token)
			tokenizer.generatePrint(tmpTok, parsedVaribels)
		}

		//If function
		if len(element.tokens) > 0 { //Check if token has sub tokens
			//Function start (ignore)
			//Function content called
			//Return argument
			tokenizer.generatePrint(element.tokens, varibels)
		}
	}

}
func evaluateMath(mathExpression []string, varibels map[string]int) (int, string) {
	command := ""
	lastValue := 0
	currentValue := 0
	currentOperator := ""
	result := 0
	for _, value := range mathExpression {
		parsedInteger, err := strconv.Atoi(value)
		if err != nil {
			currentValue = varibels[value]
		} else {
			currentValue = parsedInteger
		}
		if !isMathOperator(value) && !isCheckOperator(value) {
			fmt.Println("currentOperator:", currentOperator, "MathOperator: ", isMathOperator(currentOperator))
			if currentOperator == "" {
				lastValue = currentValue
			} else if isMathOperator(currentOperator) {
				result = evaluateMathOperation(lastValue, currentValue, currentOperator)
				currentOperator = ""
				lastValue = result
				command = "SET"
			} else {
				if evaluateCheckOperation(lastValue, currentValue, currentOperator) {
					command = "EXIT" //exits function
					break
				} else {
					return 0, "" //Nothing happens
				}
			}
		} else {
			fmt.Println(value)
			currentOperator = value
		}
	}
	fmt.Println("result:", result)
	return result, command
}
func evaluateMathOperation(x int, y int, operator string) int {
	switch operator {
	case "+":
		return x + y
	case "-":
		return x - y
	case "*":
		return x * y
	case "/":
		return x / y
	case "=":
		return y
	default:
		os.Exit(ERROR_UNIDENTIFIED_OPERATOR)
	}
	return 0
}
func evaluateCheckOperation(x int, y int, operator string) bool {
	switch operator {
	case "==":
		return x == y
	case "<":
		return x < y
	case ">":
		return x > y
	default:
		os.Exit(ERROR_UNIDENTIFIED_OPERATOR)
	}
	return false
}
func isMathOperator(operator string) bool {
	if operator == "+" ||
		operator == "-" ||
		operator == "*" ||
		operator == "/" ||
		operator == "=" {
		return true
	}
	return false
}

func isCheckOperator(operator string) bool {
	if operator == "==" ||
		operator == ">" ||
		operator == "<" {
		return true
	}
	return false
}
