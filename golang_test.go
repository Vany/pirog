package pirog

import (
	"context"
	"github.com/stretchr/testify/assert"
	"sync"
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

func TestFANOUT(t *testing.T) {
	c := make(chan int)
	g := FANOUT(c)

	c1, d1 := g()
	c2, d2 := g()
	defer d1()
	start, finish := sync.WaitGroup{}, sync.WaitGroup{}
	start.Add(2)
	finish.Add(2)
	go func() {
		start.Done()
		assert.Equal(t, 10, <-c1)
		finish.Done()
	}()
	go func() {
		start.Done()
		assert.Equal(t, 10, <-c2)
		finish.Done()
	}()
	start.Wait()
	go func() {
		c <- 10
	}()
	finish.Wait()
	d2()

	go func() {
		c <- 30
	}()

	assert.Equal(t, 30, <-c1)
	_, ok := <-c2
	assert.False(t, ok)
}

func TestFANIN(t *testing.T) {
	c := make(chan int)
	final := make(chan struct{})

	go func() {
		prg := []int{10, 20, 30, 40, 50, 60}
		for msg := range c {
			ev := prg[0]
			prg = prg[1:]
			assert.Equal(t, ev, msg)
		}
		close(final)
	}()

	g, done := FANIN(c)
	c1 := g()
	c2 := g()
	c1 <- 10
	c2 <- 20
	close(c1)
	c2 <- 30
	close(c2)
	<-time.After(100 * time.Millisecond)
	done()
	<-final
}

func TestCHANGEWATCHER(t *testing.T) {
	w := CHANGEWATCHER("main", "")
	assert.False(t, w(""))
	assert.True(t, w("lalala"))
	assert.False(t, w("lalala"))
	assert.True(t, w("lololo"))
	assert.False(t, w("lololo"))
}

func Test_SUBSCRIPTION(t *testing.T) {
	s := NewSUBSCRIPTION[int, bool]()
	ok := false
	start, finish := sync.WaitGroup{}, sync.WaitGroup{}
	start.Add(1)
	finish.Add(1)
	go func() { start.Done(); ok = <-s.Subscribe(10); finish.Done() }()
	start.Wait()
	s.Notify(10, true)
	finish.Wait()
	assert.True(t, ok)
}

func Test_REQUESTTYPE(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan REQUESTTYPE[int, int])

	go func() {
		r := <-c
		r.RESPOND(ctx, r.REQ+1)
		cancel()
	}()
	ok := false
	c <- REQUEST[int, int](10).
		THEN(ctx, func(ctx context.Context, i int) { ok = i == 11 })
	<-ctx.Done()
	assert.True(t, ok)
}
