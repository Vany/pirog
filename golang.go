package pirog

import (
	"context"
	"fmt"
	"runtime"
)

func MUST(err error) {
	if err != nil {
		panic(err)
	}
}

func MUST2[T1 any](a1 T1, err error) T1 {
	MUST(err)
	return a1
}

func MUST3[T1 any, T2 any](a1 T1, a2 T2, err error) (T1, T2) {
	MUST(err)
	return a1, a2
}

func MUST4[T1 any, T2 any, T3 any](a1 T1, a2 T2, a3 T3, err error) (T1, T2, T3) {
	MUST(err)
	return a1, a2, a3
}

func MUST5[T1 any, T2 any, T3 any, T4 any](a1 T1, a2 T2, a3 T3, a4 T4, err error) (T1, T2, T3, T4) {
	MUST(err)
	return a1, a2, a3, a4
}

// SWAPPER - same as reflect.Swapper, but template
func SWAPPER[T any](slice []T) func(i, j int) {
	return func(i, j int) {
		mid := slice[i]
		slice[i] = slice[j]
		slice[j] = mid
	}
}

// TYPEOK - strip value from explicit interface type conversion
func TYPEOK[T any](_ T, ok bool) bool {
	return ok
}

// SEND - Sends to channel obeying cancel of context
func SEND[T any](ctx context.Context, ch chan<- T, val T) {
	select {
	case <-ctx.Done():
	case ch <- val:
	}
}

// NBSEND - Sends to channel nonblockingly
func NBSEND[T any](ch chan<- T, val T) bool {
	select {
	case ch <- val:
		return true
	default:
		return false
	}
}

// RECV - receive from chan obeying context
func RECV[T any](ctx context.Context, ch <-chan T) (T, bool) {
	var NIL T
	select {
	case <-ctx.Done():
		return NIL, false
	case val, ok := <-ch:
		return val, ok
	}
}

// WAIT - for message on chan and do action once, obeying context
func WAIT[T any](ctx context.Context, ch <-chan T, cb func(T)) {
	select {
	case <-ctx.Done():
	case val, ok := <-ch:
		if ok {
			cb(val)
		}
	}
}

// FANOUT - returns function that returns channels attached to source chan
func FANOUT[T any](src <-chan T) (
	generator func() (tap <-chan T, destructor func()),
) {
	chans := make(map[chan T]struct{})
	go func() {
		for msg := range src {
			for c := range chans {
				c <- msg
			}
		}
		for c := range chans {
			close(c)
		}
	}()

	return func() (tap <-chan T, destructor func()) {
		ret := make(chan T)
		chans[ret] = struct{}{}
		return ret, func() {
			delete(chans, ret)
			close(ret)
		}
	}
}

// FANIN - returns channel generator that push all incoming into one channel
func FANIN[T any](src chan T) (generator func() chan T, destructor func()) {
	done := make(chan struct{})
	return func() chan T {
			tap := make(chan T)
			go func() {
				for {
					select {
					case <-done:
						return
					case msg, ok := <-tap:
						if !ok {
							return
						}
						src <- msg
					}
				}
			}()
			return tap
		}, func() {
			close(done)
			close(src)
		}
}

type CHANGEWATCHERFUNC[T comparable] func(n T) bool

// CHANGEWATCHER - was variabe changed from previous call
func CHANGEWATCHER[T comparable](name string, o T) CHANGEWATCHERFUNC[T] {
	return func(n T) bool {
		if o != n {
			if DEBUG {
				_, file, line, _ := runtime.Caller(1)
				fmt.Printf("<=>%s<=> %s:%d ", name, file, line)
				println(n)
			}
			o = n
			return true
		}
		return false
	}
}
