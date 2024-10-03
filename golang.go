package pirog

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"reflect"
	"runtime"
	"strings"
	"sync"
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
	return func(i, j int) { slice[i], slice[j] = slice[j], slice[i] }
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
	var zero T
	select {
	case <-ctx.Done():
		return zero, false
	case val, ok := <-ch:
		return val, ok
	}
}

// NBRECV - receive from chan obeying context
func NBRECV[T any](ch <-chan T) (T, bool) {
	var zero T
	select {
	case val, ok := <-ch:
		return val, ok
	default:
		return zero, false
	}
}

// WAIT - for message on chan and do action once, obeying context
func WAIT[T any](ctx context.Context, ch <-chan T, cb func(T)) {
	go func() {
		select {
		case <-ctx.Done():
		case val, ok := <-ch:
			if ok {
				cb(val)
			}
		}
	}()
}

// FANOUT - returns function that returns channels attached to source chan
func FANOUT[T any](src <-chan T) (
	generator func() (tap <-chan T, destructor func()),
) {
	chans := make(map[chan T]struct{})
	go func() {
		for msg := range src {
			for c := range chans {
				NBSEND(c, msg)
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

type COMMENTABLETYPE[T any] struct{ T T }

func COMMENTABLE[T any](in T) COMMENTABLETYPE[T] { return COMMENTABLETYPE[T]{T: in} }
func (C *COMMENTABLETYPE[T]) MarshalYAML() (interface{}, error) {
	t := reflect.TypeOf(C.T)
	v := reflect.ValueOf(C.T)
	if t.Kind() == reflect.Ptr {
		return COMMENTABLE(v.Elem().Interface()), nil
	}

	b := new(strings.Builder)

	switch t.Kind() {
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			cmt := ""
			if comment := f.Tag.Get("comment"); comment != "" {
				cmt = "# " + comment + "\n"
			}

			b.WriteString(cmt)
			b.WriteString(f.Name + ": ")
			yaml.NewEncoder(b).Encode(COMMENTABLE(v.Field(i).Interface()))
		}
		n := yaml.Node{}
		n.SetString(b.String())
		return n, nil

	default:
		return C.T, nil
	}
}

type SUBSCRIPTION[A comparable, T any] struct {
	sync.Mutex
	M map[A][]chan T
}

func NewSUBSCRIPTION[A comparable, T any]() *SUBSCRIPTION[A, T] {
	return &SUBSCRIPTION[A, T]{M: make(map[A][]chan T)}
}

func (s *SUBSCRIPTION[A, T]) Close(id A) {
	s.Lock()
	for _, c := range s.M[id] {
		close(c)
	}
	delete(s.M, id)
	s.Unlock()
}

func (s *SUBSCRIPTION[A, T]) Notify(id A, data T) {
	s.Lock()
	if !HAVEKEY(s.M, id) {
		s.M[id] = nil
	}
	var zero A
	for _, c := range append(s.M[id], s.M[zero]...) {
		_ = NBSEND(c, data)
	}
	s.Unlock()
}

func (s *SUBSCRIPTION[A, T]) Subscribe(id A) chan T {
	c := make(chan T)
	s.Lock()
	s.M[id] = append(s.M[id], c)
	s.Unlock()
	return c
}

func (s *SUBSCRIPTION[A, T]) UnSubscribe(id A, c chan T) {
	s.Lock()
	s.M[id] = GREP(s.M[id], func(in chan T) bool { return in != c })
	if s.M[id] == nil {
		delete(s.M, id)
	}
	s.Unlock()
}

func (s *SUBSCRIPTION[A, T]) Has(id A) bool {
	s.Lock()
	h := HAVEKEY(s.M, id)
	s.Unlock()
	return h
}

type REQUESTTYPE[REQ any, RES any] struct {
	REQ REQ
	RES chan RES
}

func REQUEST[REQ any, RES any](req REQ) REQUESTTYPE[REQ, RES] {
	return REQUESTTYPE[REQ, RES]{REQ: req, RES: make(chan RES)}
}

func (r REQUESTTYPE[REQ, RES]) RESPOND(ctx context.Context, res RES) REQUESTTYPE[REQ, RES] {
	SEND(ctx, r.RES, res)
	return r
}

func (r REQUESTTYPE[REQ, RES]) SEND(ctx context.Context, c chan REQUESTTYPE[REQ, RES]) REQUESTTYPE[REQ, RES] {
	SEND(ctx, c, r)
	return r
}

func (r REQUESTTYPE[REQ, RES]) THEN(ctx context.Context, f func(context.Context, RES)) REQUESTTYPE[REQ, RES] {
	go func() { WAIT(ctx, r.RES, func(r RES) { f(ctx, r) }) }()
	return r
}
