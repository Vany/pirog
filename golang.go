package pirog

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
