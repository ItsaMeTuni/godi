package godi

import (
	"fmt"
	"reflect"
)

type MissingProviderError struct {
	fn Fn
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
	fn Fn
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