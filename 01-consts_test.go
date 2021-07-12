package main

import (
	"testing"
)

func TestAbcConst(t *testing.T) {
	const myConst = "abc"
	if myConst != "abc" {
		t.Error("Const was changed")
	}
}

func TestAbcVar(t *testing.T) {
	myVar := "abc"
	const myConst = "cba"
	if myVar == myConst {
		t.Error("Variable was initialized wrong")
	}
	myVar = myConst
	if myVar != myConst {
		t.Error("Variable wasn't changed, something is wrong")
	}
}
