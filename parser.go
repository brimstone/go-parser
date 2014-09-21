package parser

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
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
	return token //, fmt.Errorf("Don't know how to handle this")
}

func Parse(env Env, input string) (bool, error) {
	tokenExp := regexp.MustCompile("[0-9]+|[-+]|[a-z]+|[=<>]|[|&]")
	tokens := tokenExp.FindAllString(input, -1)
	fmt.Println(tokens)
	//spew.Dump(tokens)

	// easy out
	if len(tokens) == 0 {
		return true, nil
	}
	// Attempt to loop through our tokens until we can reduce our tokens to 1

	dirty := false
TokenCheck:
	for dirty == false && len(tokens) > 1 {
		for i, _ := range tokens {
			fmt.Println("Evaluating", i)
			// check dualnary operators
			if i > 0 && len(tokens) > i {
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
					dirty = true
					tokens = append(tokens[:i], tokens[i+2:]...)
					break TokenCheck
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
					dirty = true
					tokens = append(tokens[:i], tokens[i+2:]...)
					break TokenCheck
				case ">":
					return false, fmt.Errorf("No greater than yet")
				case "|":
					if left == "true" {
						dirty = true
						tokens[i-1] = "true"
						tokens = append(tokens[:i], tokens[i+2:]...)
						//spew.Dump(tokens)
						break TokenCheck
					} else if left != "false" {
						return false, fmt.Errorf("Expected 'true' or 'false', found", left)
					}
					if right == "true" {
						dirty = true
						tokens[i-1] = "true"
						tokens = append(tokens[:i], tokens[i+2:]...)
						//spew.Dump(tokens)
						break TokenCheck
					} else if right != "false" {
						return false, fmt.Errorf("Expected 'true' or 'false', found", right)
					}
				case "&":
					if left == "true" && right == "true" {
						dirty = true
						tokens[i-1] = "true"
						tokens = append(tokens[:i], tokens[i+2:]...)
						//spew.Dump(tokens)
						break TokenCheck
					}
					if left != "false" && left != "true" {
						return false, fmt.Errorf("Expected 'true' or 'false', found '%s'", left)
					}
					if right != "false" && right != "true" {
						return false, fmt.Errorf("Expected 'true' or 'false', found '%s'", right)
					}
					dirty = true
					tokens[i-1] = "false"
					tokens = append(tokens[:i], tokens[i+2:]...)
					spew.Dump(tokens)
					break TokenCheck
				default:
					return false, fmt.Errorf("Don't know how to parse", tokens[i])
				}
			}
			//tokens = tokens[:len(tokens)-1]
			// since we changed the length of our array, we need to break this loop
			//break
		}
	}

	if item, ok := env[tokens[0]]; ok {
		switch t := item.(type) {
		case int:
			if t == 0 {
				tokens[0] = "false"
			} else {
				tokens[0] = "true"
			}
		case bool:
			fmt.Println("it's a bool")
			// maybe we should keep it in the determined form?
			if t {
				tokens[0] = "true"
			} else {
				tokens[0] = "false"
			}
		}
		//tokens[0] = env[tokens[0]].(string)
	}

	// Now that we've reduced our tokens to only one element,
	// figure out what it is and return properly
	switch tokens[0] {
	case "":
		return true, nil
	case "true":
		return true, nil
	case "false":
		return false, nil
	}
	return false, fmt.Errorf("Error parsing '%s'", tokens[0])
}
