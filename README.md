# go-parser
This is a simple string parser. It handles an environment full of variables, and a string to parse it down to a simple true or false boolean value.

Supports:
- 'true'
- 'false'
- < less than
- > greater than
- = Equivalence 
- | boolean OR
- & boolean AND
- () Nesting expressions

## Requirements
* golang

## Installation
Simply `import "github.com/brimstone/go-parser"` in your project.

## Usage
```go
package main

import (
	"fmt"
	"github.com/brimstone/go-parser"
)

func main() {
	env := make(parser.Env)
	env["foo"] = true
	env["bar"] = 1
	fmt.Println(parser.Parse(env, "bar=0|foo"))
}
```

## Known issues
* no real string handling
* no integer math
* no negative numbers
* no unequal operator
