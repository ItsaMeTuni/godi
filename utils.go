package godi

import (
	"reflect"
	"runtime"
)

// Get fn's package path and name like github.com/ItsaMeTuni/godi.Inject
func getFnPath(fn Fn) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}