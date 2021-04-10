package godi_test

import (
	"errors"
	"github.com/ItsaMeTuni/godi"
	"testing"
)

func TestInjectConcrete(t *testing.T) {
	type Provider1 struct {}

	ran := false
	fn := func(p1 Provider1) {
		ran = true
	}
	providers := []interface{}{ Provider1{} }

	retVals, err := godi.Inject(fn, providers)
	if err != nil {
		t.Fatal(err)
	}

	if len(retVals) != 0 {
		t.Fatal(retVals)
	}

	if !ran {
		t.Fatal("fn didnt run")
	}
}

type Iface interface { Foo() }
type Provider struct {}
func (_ Provider) Foo() {}

func TestInjectInterface(t *testing.T) {
	ran := false

	fn := func(p Iface) {
		_ = p.(Provider) // panic if p is nil or holds a nil
		ran = true
	}

	providers := []interface{}{ Iface(Provider{}) }

	retVals, err := godi.Inject(fn, providers)
	if err != nil {
		t.Fatal(err)
	}

	if len(retVals) != 0 {
		t.Fatal(retVals)
	}

	if !ran {
		t.Fatal("fn didnt run")
	}
}


func TestInjectDeepPtr(t *testing.T) {
	ran := false

	fn := func(p ****Provider) {
		_ = ****p // panic if p is nil or holds a nil
		ran = true
	}

	a := Provider{}
	b := &a
	c := &b
	d := &c

	providers := []interface{}{ &d }

	retVals, err := godi.Inject(fn, providers)
	if err != nil {
		t.Fatal(err)
	}

	if len(retVals) != 0 {
		t.Fatal(retVals)
	}

	if !ran {
		t.Fatal("fn didnt run")
	}
}

func TestAssertFn(t *testing.T) {
	fn := func() (int, string, error) { return 0, "", nil }

	err := godi.AssertFn(fn, []interface{}{ 0, "", errors.New("") })
	if err != nil {
		t.Fatal(err)
	}
}

func TestAssertFnMissingOne(t *testing.T) {
	fn := func() (int, error) { return 0, nil }

	err := godi.AssertFn(fn, []interface{}{ 0, "", errors.New("") })
	if _, ok := err.(godi.InvalidSignatureError); !ok {
		t.Fatal("expected InvalidSignatureError error")
	}
}

func TestAssertFnOneExtra(t *testing.T) {
	fn := func() (int, string, uint, error) { return 0, "", 0, nil }

	err := godi.AssertFn(fn, []interface{}{ 0, "", errors.New("") })
	if _, ok := err.(godi.InvalidSignatureError); !ok {
		t.Fatal("expected InvalidSignatureError error")
	}
}

func TestAssertNotFn(t *testing.T) {
	fn := ""

	err := godi.AssertFn(fn, []interface{}{ 0, "", errors.New("") })
	if _, ok := err.(godi.NotAFuncError); !ok {
		t.Fatal("expected NotAFuncError error")
	}
}