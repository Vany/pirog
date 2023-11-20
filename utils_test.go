package pirog

import (
	"io"
	"os"
	"testing"
)

type S1 struct {
	I int
}

func TestPutToStruct(t *testing.T) {
	app := &struct {
		S *S1
		R io.Reader
	}{}
	PutToStruct(app, &S1{})
	PutToStruct(app, os.Stdin)
	if app.S == nil || app.R == nil {
		t.Fail()
	}
}
