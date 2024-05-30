package pirog

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

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
