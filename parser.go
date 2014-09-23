package parser

import (
	"fmt"
	"regexp"
	"strconv"
)

type Env map[string]interface{}

func resolveEnv(env Env, token string) string {
	if item, ok := env[token]; ok {
		switch t := item.(type) {
		case string:
			return item.(string)
		case int:
			return strconv.Itoa(item.(int))
		case bool:
			if t {
				return "true"
			} else {
				return "false"
			}
		}
	}
	return token
}

func parseTokens(env Env, tokens []string) (bool, error) {
	// easy out
	if len(tokens) == 0 {
		return true, nil
	}
	// Attempt to loop through our tokens until we can reduce our tokens to 1

	level := 2
TokenCheck:
	for level > 0 && len(tokens) > 1 {
		for i, _ := range tokens {
			switch level {
			case 2:
				// check dualnary operators
				if i > 0 && len(tokens) > i+1 {
					left := resolveEnv(env, tokens[i-1])
					right := resolveEnv(env, tokens[i+1])
					switch tokens[i] {
					case "<":
						// attempt to convert left to an int
						leftInt, err := strconv.Atoi(left)
						if err != nil {
							return false, fmt.Errorf("Expect type int, found", left)
						}
						// attempt to convert right to an int
						rightInt, err := strconv.Atoi(right)
						if err != nil {
							return false, fmt.Errorf("Expect type int, found", right)
						}
						// do the op
						if leftInt < rightInt {
							tokens[i-1] = "true"
						} else {
							tokens[i-1] = "false"
						}
						tokens = append(tokens[:i], tokens[i+2:]...)
						break TokenCheck
					case ">":
						// attempt to convert left to an int
						leftInt, err := strconv.Atoi(left)
						if err != nil {
							return false, fmt.Errorf("Expect type int, found", left)
						}
						// attempt to convert right to an int
						rightInt, err := strconv.Atoi(right)
						if err != nil {
							return false, fmt.Errorf("Expect type int, found", right)
						}
						// do the op
						if leftInt > rightInt {
							tokens[i-1] = "true"
						} else {
							tokens[i-1] = "false"
						}
						tokens = append(tokens[:i], tokens[i+2:]...)
						break TokenCheck
					}
				}
			case 1:
				// check dualnary operators
				if i > 0 && len(tokens) > i+1 {
					left := resolveEnv(env, tokens[i-1])
					right := resolveEnv(env, tokens[i+1])
					switch tokens[i] {
					case "=":
						// do the op
						// this also covers bools and ints since they're both strings
						// FALSE != false however
						if left == right {
							tokens[i-1] = "true"
						} else {
							tokens[i-1] = "false"
						}
						tokens = append(tokens[:i], tokens[i+2:]...)
						break TokenCheck
					case "|":
						if left == "true" {
							tokens[i-1] = "true"
							tokens = append(tokens[:i], tokens[i+2:]...)
							break TokenCheck
						} else if left != "false" {
							return false, fmt.Errorf("Expected 'true' or 'false', found", left)
						}
						if right == "true" {
							tokens[i-1] = "true"
							tokens = append(tokens[:i], tokens[i+2:]...)
							break TokenCheck
						} else if right != "false" {
							return false, fmt.Errorf("Expected 'true' or 'false', found", right)
						}
					case "&":
						if left == "true" && right == "true" {
							tokens[i-1] = "true"
							tokens = append(tokens[:i], tokens[i+2:]...)
							break TokenCheck
						}
						if left != "false" && left != "true" {
							return false, fmt.Errorf("Expected 'true' or 'false', found '%s'", left)
						}
						if right != "false" && right != "true" {
							return false, fmt.Errorf("Expected 'true' or 'false', found '%s'", right)
						}
						tokens[i-1] = "false"
						tokens = append(tokens[:i], tokens[i+2:]...)
						break TokenCheck
					}
				}
			}
		}
		level--
	}

	// Now that we've reduced our tokens to only one element,
	// figure out what it is and return properly
	if item, ok := env[tokens[0]]; ok {
		switch t := item.(type) {
		case int:
			if item == 0 {
				tokens[0] = "false"
			} else {
				tokens[0] = "true"
			}
		case bool:
			// maybe we should keep it in the determined form?
			if t {
				tokens[0] = "true"
			} else {
				tokens[0] = "false"
			}
		}
	}

	// Handle if the last element is an integer
	itemInt, err := strconv.Atoi(tokens[0])
	if err == nil {
		if itemInt == 0 {
			tokens[0] = "false"
		} else {
			tokens[0] = "true"
		}
	}

	// We should now be down to simply true or false
	switch tokens[0] {
	case "true":
		return true, nil
	case "false":
		return false, nil
	}
	// Only invalid tokens should be this far.
	return false, fmt.Errorf("Error parsing '%s'", tokens[0])
}

func Parse(env Env, input string) (bool, error) {
	// Handle:
	// - numbers
	// - letters (assuming they're variables in the environment
	// - equal signs, equivalence
	// - less than
	// - greater than
	// - boolean OR
	// - boolean AND
	tokenExp := regexp.MustCompile("[0-9]+|[a-z]+|[=<>]|[|&]")
	tokens := tokenExp.FindAllString(input, -1)

	return parseTokens(env, tokens)

}
