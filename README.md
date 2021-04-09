# GoDI

Go Dependency Injection made easy.

This package allows you to easily inject your function's dependencies.
Here's some code:

```go
package main

import (
	"github.com/ItsaMeTuni/godi"
	"strconv"
)

// Create a simple provider
type MagicProvider struct {
	magic int
}

// Create a function that needs a provider. We're going
// to inject this function with a MagicProvider.
func foo(provider MagicProvider) string {
	return strconv.Itoa(provider.magic)
}

func main() {
	// Create provider
	myProvider := MagicProvider{ magic: 42 }

	// Create a list of providers (currently we only have one,
	// which is myProvider).
	providers := godi.Providers{ myProvider }
	
	// Inject foo with the providers of providers and
	// execute it. 
	// This is equivalent to calling foo like foo(myProvider)	
	retVals, err := godi.Inject(foo, providers)
	if err != nil {
		return
    }
	
    // Print the string foo returned
	println(retVals[0].String())
}
```