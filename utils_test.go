package pirog

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
	"time"
)

type S1 struct{}
type S2 struct{}
type I1 struct{}
type I2 struct{}
type I3 struct{}

func (s *S1) InitTest(ctx context.Context) error { TEST_VARIABLE++; return nil }
func (i *I1) InitTest(ctx context.Context) error { TEST_VARIABLE++; return nil }
func (i *I2) InitTest(ctx context.Context) error { TEST_VARIABLE++; return nil }
func (i *I3) InitTest(ctx context.Context) error { TEST_VARIABLE++; return nil }
func (i *I1) isI1I()                             {}
func (i *I2) isI2I()                             {}
func (i *I3) isI3I()                             {}

type I1I interface{ isI1I() }
type I2I interface{ isI2I() }
type I3I interface{ isI3I() }

type App struct {
	S1 *S1
	S2 *S2
	R  io.Reader
	I1 I1I
	I2 I2I
	I3 I3I
}

var TEST_VARIABLE = 0

func TestExecuteOnAllFields(t *testing.T) {
	app := &App{}
	app.S1 = &S1{}
	app.S2 = &S2{}
	app.R = os.Stdin
	app.I1 = &I1{}
	app.I2 = &I2{}

	assert.NoError(t, ExecuteOnAllFields(context.Background(), app, "InitTest"))
	assert.Equal(t, 2, TEST_VARIABLE)

	app.I3 = &I3{}
	assert.NoError(t, ExecuteOnAllFields(context.Background(), app, "InitTest"))
	assert.Equal(t, 5, TEST_VARIABLE)
}

func TestSEND(t *testing.T) {
	ctx, c := context.WithTimeout(context.Background(), time.Second)
	ch := make(chan string)
	gofin := make(chan struct{})
	gf := func() { gofin <- struct{}{} }
	summ := 0

	go func() { SEND(ctx, ch, "lala"); summ++; gf() }()
	assert.Equal(t, 0, summ)
	<-ch
	<-gofin
	assert.Equal(t, 1, summ)

	go func() { SEND(ctx, ch, "lolo"); summ++; gf() }()
	assert.Equal(t, 1, summ)
	c()
	<-gofin
	assert.Equal(t, 2, summ)
}
