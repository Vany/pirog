package pirog

// GREP - This is filter, that leaves only that elements that trigerrs callback function to return true
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

// MAP - This is part of mapreduce and almost full copy of perl's map. It transforms input array to output array with callback function.
func MAP[IN any, OUT any](arr []IN, f func(IN) OUT) []OUT {
	accum := make([]OUT, len(arr))
	for i, v := range arr {
		accum[i] = f(v)
	}
	return accum
}

// KEYS - Returns full set of keys from map to use it further
func KEYS[K comparable, V any](in map[K]V) []K {
	acc := make([]K, 0, len(in))
	for k := range in {
		acc = append(acc, k)
	}
	return acc
}

// VALUES - Returns full set of values from map to use it further
func VALUES[K comparable, V any](in map[K]V) []V {
	acc := make([]V, 0, len(in))
	for k := range in {
		acc = append(acc, in[k])
	}
	return acc
}

// HAVEKEY - Just indicates do we have key in map, or no.
func HAVEKEY[K comparable, V any](in map[K]V, key K) bool {
	_, have := in[key]
	return have
}

// ANYKEY - Returns any arbitrary key from map.
func ANYKEY[K comparable, V any](in map[K]V) K {
	for k := range in {
		return k
	}
	var k K
	return k
}

// ANYWITHDRAW - Chooses arbitrary key from map, delete it and return.
func ANYWITHDRAW[K comparable, V any](in map[K]V) (K, V) {
	for k := range in {
		ret := in[k]
		delete(in, k)
		return k, ret
	}
	var k K
	var v V
	return k, v
}

// TERNARY - ternary operator
func TERNARY[T any](e bool, a, b T) T {
	if e {
		return a
	} else {
		return b
	}
}

// REDUCE - Takes array and applies callback function to aggregate object and each element of array. Starts from init.
func REDUCE[IN any, ACC any](init ACC, in []IN, f func(int, IN, *ACC)) *ACC {
	acc := new(ACC)
	*acc = init
	for i, el := range in {
		f(i, el, acc)
	}
	return acc
}

// EXPLODE - Explodes number to range of values.
func EXPLODE[T any](num int, f func(int) T) []T {
	acc := make([]T, num)
	for i := 0; i < num; i++ {
		acc[i] = f(i)
	}
	return acc
}

// FLATLIST - flaterns list of lists to just list
func FLATLIST[T any](arrs [][]T) []T {
	ret := make([]T, 0, len(arrs))
	for _, v := range arrs {
		ret = append(ret, v...)
	}
	return ret
}

// COALESCE - return first non zero(false) value
func COALESCE[T comparable](in ...T) T {
	var zero T
	for _, t := range in {
		if t != zero {
			return t
		}
	}
	return zero
}
