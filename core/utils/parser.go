package utils

import (
	"fmt"
	"strconv"
	"strings"
)

//this function doesn't care about unicode and should be rewritten
func ParseExpr(input string) (*map[string]string, error) {
	//{} is supposed to presented
	if len(input) < 2 {
		return nil, fmt.Errorf("malformed query string")
	}
	//trim first and last character
	expr := strings.Split(input[1:len(input)-1], ",")
	parsedExpr := make(map[string]string)
	for _, token := range expr {
		labelExprTokens := strings.SplitN(token, "=", 2)
		if len(labelExprTokens) != 2 {
			return nil, fmt.Errorf("malformed query string")
		}
		label := labelExprTokens[0]
		value, err := strconv.Unquote(labelExprTokens[1])
		if err != nil {
			return nil, fmt.Errorf("malformed query string: %s", err)
		}
		parsedExpr[label] = value
	}
	return &parsedExpr, nil
}
