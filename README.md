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

## What is DI and why should I use it?

Dependency Injection is a way to implement Inversion of Control in your code. It basically
allows you to _declare_ the _dependencies_ of a function (or classes in other languages), instead
of creating the dependencies inside the function or using singletons.

For example, take a function that _depends_ on a (fictional) `DatabaseConn` object to make some queries to the 
database. There are basically two ways to do it, the not so good way (a singleton) and
the better one (dependency injection).

With a singleton
```
var dbConn *DatabaseConn

func main() {
    dbConn := &DatabaseConn{ /* configuration ommitted */ }
    
    foo("0001")
}

func foo(userId string) {
    // Don't write queries like this, it's prone to SQL injection, use a
    // proper query builder sanitizer
    dbConn.query("SELECT * FROM users WHERE id = " + userId)
    
    // do cool stuff with the user
}
```

This is great and all for small and simple projects. But what happens when your project grows and you need to
write tests for `foo`? And what happens if you need a different connection somewhere else? Yes, you could mock
the `dbConn` singleton and you could use multiple singletons for different connections, but those options
are messy and painful to maintain.

Dependency injection solves this problem by having you _declare_ `foo`'s dependency on a `DatabaseConn` in
`foo`'s signature and automatically _injecting_ the dependency when you call `foo`.

```
func main() {
    dbConn := &DatabaseConn{ /* configuration ommitted */ }
    
    providers := godi.Providers { &dbConn }
    
    godi.Inject(foo, providers)
}

func foo(dbConn *DatabaseConn) {
    // Don't write queries like this, it's prone to SQL injection, use a
    // proper query builder sanitizer
    dbConn.query("SELECT * FROM users")
    
    // do cool stuff with the user
}
```

Now, you can specify that you want a `*DatabaseConn` in foo and provide it from a list of _providers_. When
`godi.Inject` is called, it will try to find a provider for each parameter of `foo` by type and then, when it
finds a provider for each parameter it will call `foo` with the providers.

Now when need to test your code you can just mock a database connection and call `foo` manually with it (or
event use DI in the test). No more fiddly and messy singletons.