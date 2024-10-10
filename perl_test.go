package pirog

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
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

func TestEXPLODEREDUCE(t *testing.T) {
	x := REDUCE(1+0i, EXPLODE(6, func(i int) complex128 {
		return 0 + 1i
	}), func(i int, in complex128, acc *complex128) {
		*acc *= in
	})
	if *x == -1 {
		return
	}
	t.Error("failed")
}

func Test_MAP(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7}
	res := MAP(arr, func(in int) string {
		return fmt.Sprintf("%d", in)
	})

	assert.Equal(t, len(arr), len(res))
	assert.Equal(t, res[5], "6")
}
