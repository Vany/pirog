package pirog

func GREP[T any](arr []T, f func(T) bool) []T {
	var accum []T
	for _, x := range arr {
		if !f(x) {
			continue
		}
		accum = append(accum, x)
	}
	return accum
}

func MAP[IN any, OUT any](arr []IN, f func(IN) OUT) []OUT {
	var accum []OUT
	for _, in := range arr {
		accum = append(accum, f(in))
	}
	return accum
}

func KEYS[K comparable, V any](in map[K]V) []K {
	acc := make([]K, 0, len(in))
	for k := range in {
		acc = append(acc, k)
	}
	return acc
}

func HAVEKEY[K comparable, V any](in map[K]V, key K) bool {
	for k := range in {
		if k == key {
			return true
		}
	}
	return false
}

func ANYKEY[K comparable, V any](in map[K]V) K {
	for k := range in {
		return k
	}
	var k K
	return k
}

func TERNARY[T any](e bool, a, b T) T {
	if e {
		return a
	} else {
		return b
	}
}

func MUST(err error) {
	if err != nil {
		panic(err)
	}
}
