package godi

import (
	"errors"
	"reflect"
)

type Providers = []interface{}
type Fn = interface{}

// Assert whether fn is a func and whether its return value types match
// the types of the values in returnValues (order matters).
//
// To skip assertion of the return values, just set returnValues to nil.
//
// Example:
// To check if foo is a function and has (int, error) as return value,
// do the following:
// func foo() (int, error) { return 0, nil }
//
// assertFn(foo, []interface{} { 0, errors.New("") })
func AssertFn(fn Fn, returnValues []interface{}) error {
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		return NotAFuncError{ t: fnType }
	}

	if returnValues != nil {

		expectedRetValTypes := make([]reflect.Type, len(returnValues))
		for i, retVal := range returnValues {
			expectedRetValTypes[i] = reflect.TypeOf(retVal)
		}

		if len(returnValues) != fnType.NumOut() {
			return InvalidSignatureError{
				fn: fn,
				expectedRetVals: expectedRetValTypes,
			}
		}

		for i := 0; i < fnType.NumOut(); i++ {
			if !expectedRetValTypes[i].AssignableTo(fnType.Out(i)) {
				return InvalidSignatureError{
					fn: fn,
					expectedRetVals: expectedRetValTypes,
				}
			}
		}
	}

	return nil
}


// Injects a fn with the given providers and calls it immediately.
// Returns the values returned by fn and an error, if any.
//
// How it works:
// This function looks at fn's signature and figures out the types of
// parameters it wants. It then looks for providers with the same types, if
// a provider match is found for all parameters, fn is executed and it's return
// values are returned.
//
// If a provider isn't found for any parameters an error will be returned.
//
// Note: if fn returns an error it will NOT be in the error return value of
// this function, it's going to be in the returned []reflect.Value.
//
// Note: if you want to make a provider accessible as an interface, you
// have to pass a pointer to the interface as provider. Example:
// Assume myProvider implements MyIface
// If we want the handler to receive a MyIface parameter,
// providers has to be Providers{ &MyIface(myProvider) }.
func Inject(
	fn Fn,
	providers Providers,
) ([]reflect.Value, error) {

	fnType := reflect.TypeOf(fn)

	if fnType.Kind() != reflect.Func {
		return nil, errors.New("fn is not func")
	}

	// Create an array of arguments that will be used to call fn.
	// These arguments are found by matching the fn parameter types
	// against the provider types, when there is a match the respective provider
	// will be used as argument for the respective parameter.
	paramCount := fnType.NumIn()
	args := make([]reflect.Value, paramCount)

	for i := 0; i < paramCount; i++ {
		paramType := fnType.In(i)

		for _, provider := range providers {

			providerType := reflect.TypeOf(provider)

			if providerType.AssignableTo(paramType) {
				args[i] = reflect.ValueOf(provider)
			}
		}

		if !args[i].IsValid() {
			return nil, MissingProviderError{
				fn: fn,
				paramIdx: i,
			}
		}
	}

	// Call the handler
	returnValues := reflect.ValueOf(fn).Call(args)

	return returnValues, nil
}