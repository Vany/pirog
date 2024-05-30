package pirog

import "context"

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

// COPYCHAN - returns function that returns channels attached to source chan
func COPYCHAN[T any](src chan T) (
	generator func() (tap chan T, destructor func()),
	destructor func(),
) {
	var chans map[chan T]struct{}
	done := make(chan struct{})

	go func() {
		for {
			select {
			case msg := <-src:
				for c, _ := range chans {
					c <- msg
				}
			case <-done:
				for c, _ := range chans {
					close(c)
					delete(chans, c)
				}
				return
			}
		}
	}()

	return func() (tap chan T, destructor func()) {
			ret := make(chan T)
			chans[ret] = struct{}{}
			return ret, func() {
				delete(chans, ret)
			}
		},
		func() {
			done <- struct{}{}
		}
}
