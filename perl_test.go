package pirog

import (
	"errors"
	"testing"
)

func TestMUST(t *testing.T) {
	var mustNil bool
	defer func() {
		if !mustNil {
			t.Error("MUST() reacted to nil")
		}
		e := recover()
		if e == nil {
			t.Error("MUST() not reacted to error")
		}
	}()
	MUST(nil)
	mustNil = true
	MUST(errors.New(""))
}
