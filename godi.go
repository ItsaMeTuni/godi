package godi

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
)

func getFnPath(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

type MissingProviderError struct {
	fn interface{}
	paramIdx int
}

func (e MissingProviderError) Error() string {

	fnType := reflect.TypeOf(e.fn)

	return fmt.Sprintf(
		"Could not find a provider of type %s for parameter %d of %s",
		fnType.In(e.paramIdx).Name(),
		e.paramIdx,
		getFnPath(e.fn),
	)
}

type NotAFuncError struct {
	t reflect.Type
}

func (e NotAFuncError) Error() string {
	return fmt.Sprintf(
		"You provided something other than a func to godi.Inject. Provided type: %s.",
		e.t.Name(),
	)
}

type InvalidSignatureError struct {
	fn interface{}
	expectedRetVals []reflect.Type
}

func (e InvalidSignatureError) Error() string {

	fnType := reflect.TypeOf(e.fn)

	fnRetVals := make([]reflect.Type, fnType.NumOut())
	for i := 0; i < len(fnRetVals); i++ {
		fnRetVals[i] = fnType.Out(i)
	}

	return fmt.Sprintf(
		"Wrong signature. Func %s has signature %s, expected %s.",
		getFnPath(e.fn),
		fnRetVals,
		e.expectedRetVals,
	)
}

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
func AssertFn(fn interface{}, returnValues []interface{}) error {
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		// TODO: use custom error type
		return errors.New("fn is not a func")
	}

	if returnValues != nil {
		if len(returnValues) != fnType.NumOut() {
			// TODO: use custom error type
			return errors.New("invalid return values")
		}

		for i := 0; i < fnType.NumOut(); i++ {
			retValType := reflect.TypeOf(returnValues[i])
			if !fnType.Out(i).AssignableTo(retValType) {
				// TODO: use custom error type
				return errors.New("invalid return values")
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
// providers has to be []interface{}{ &MyIface(myProvider) }.
func Inject(
	fn interface{},
	providers []interface{},
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