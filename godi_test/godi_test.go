package godi_test

import (
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

