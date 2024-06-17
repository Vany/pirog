package pirog

import (
	"context"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

type S1 struct {
	S2 *S2 `inject:""`
}
type S2 struct {
	S1   S1  `inject:""`
	SAny any `inject:"S1"`
}
type I1 struct{ I int }
type I2 struct {
	S1 S1 `inject:""`
}
type I3 struct {
	I1 *I1 `inject:"I1"`
	I2 *I2 `inject:"I2"`
	I3 *I3 `inject:"I3"`
	I4 I1I `inject:"I1"`
}

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
	S1 *S1 `injectable:""`
	S2 *S2 `injectable:""`
	R  io.Reader
	I1 I1I `injectable:""`
	I2 I2I `injectable:""`
	I3 I3I `injectable:""`
	I  int
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
	assert.Equal(t, 3, TEST_VARIABLE)

	app.I3 = &I3{}
	assert.NoError(t, ExecuteOnAllFields(context.Background(), app, "InitTest"))
	assert.Equal(t, 7, TEST_VARIABLE)
}

func TestInjectComponents(t *testing.T) {
	app := &App{}
	app.S1 = &S1{}
	app.S2 = &S2{}
	app.R = os.Stdin
	app.I1 = &I1{}
	app.I2 = &I2{}
	app.I3 = &I3{}

	InjectComponents(app)
	assert.NotNil(t, app.S1)
	assert.NotNil(t, app.S2)
	assert.NotNil(t, app.S2.S1.S2)
	assert.NotNil(t, app.S1.S2)

	assert.NotNil(t, app.I3.(*I3).I1)
	assert.NotNil(t, app.I3.(*I3).I2)
}
