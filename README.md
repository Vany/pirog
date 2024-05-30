[![License](https://img.shields.io/github/license/vany/pirog)](./LICENSE)
[![tag](https://img.shields.io/github/tag/vany/pirog.svg)](https://github.com/vany/pirog/releases)
[![Test](https://github.com/Vany/pirog/actions/workflows/test.yml/badge.svg)](https://github.com/Vany/pirog/actions/workflows/test.yml)
[![goreport](https://goreportcard.com/badge/github.com/vany/pirog)](https://goreportcard.com/badge/github.com/vany/pirog)


# pirog
Golang mapreduce primitives and other cool stuff from perl and javascript.

Main idea is to use commonly known and well proven constructions introduced in other languages.
Constructions have distinguishable appearance and pretend to be a part of the language rather than just functions.

## Useage
Just import it as is, but if you want and have enough dare, use special symbol to package reference like
```go
import . "github.com/vany/pirog"
```
Then just use it. 

### MAP(array, callback) array
This is part of mapreduce and almost full copy of perl's map. It transforms input array to output array with callback function.
```go
type Person struct {
    FirstName  string
    SecondName string
}

people := []Person{{"Vasya", "Pupkin"}, {"Vasya", "Lozhkin"}, {"Salvador", "Dalí"}}
out := MAP(people, func(p Person) string{
	return fmt.Sprintf("%s %s", p.FirstName, p.SecondName)
})
```
`out` now is []string containing concatenated names.

### GREP(array, callback) array
This is filter, that leaves only that elements that trigerrs callback function to return true
```go
fakePeople := GREP(out, func(s string) bool {
    return strings.HasSuffix(s, "Pupkin")
})
```
`fakePeople` now is []string and contains just `"Vasya Pupkin"`

### REDUCE(init, array, callback) aggregate
Takes array and applies callback function to aggregate object and each element of array. Starts from init.
```go
x := REDUCE(1+0i, EXPLODE(6, func(i int) complex128 {
    return 0 + 1i
}), func(i int, in complex128, acc *complex128) {
    *acc *= in
})
// rounds dot for 3π so result will be -1
```

### EXPLODE(number, callback) array
Explodes number to range of values
```go
smallEnglishLetters := EXPLODE(26, func(in int) string { return string([]byte{byte('a' + in)}) }) {
```

### KEYS(map) array
Returns just full set of keys from map to use it further
```go
artistsMap := map[string]string{
	"Vasya":"Lozhkin",
	"Salvador":"Dalí",
}
AllLozhkins := GREP(KEYS(artistsMap), func(in string) string }{
	return artistsMap[in] == "Lozhkin" 
})
```
`AllLozhkins` will be `[]string{"Vasya"}`

### VALUES(map) array
Returns full set of values from map to use it further

### HAVEKEY(map, key) bool
Just indicates do we have key in map, or no.

### ANYKEY(map) key
Returns any arbitrary key from map.

### ANYWITHDRAW(map) key, value
Chooses arbitrary key from map, delete it and return.

### FLATLIST(list of lists) list
Flaterns list of lists, used when you have MAP in MAP, but need flat list outside.


### MUST(err)
Validates err for nil and panics other way. When you in CLI or sure that here can not be an error.
```go
MUST(json.NewEncoder(buff).Encode(object))

```

### MUST2 - MUST5 
Same as must, but returns values number at the end of name is number of the arguments.
```go
files := GREP(MUST2(os.ReadDir(".")), func (in os.DirEntry) bool {
	return !in.IsDir()
})
```

### SWAPPER(array) func
Same as reflect.Swapper(), generates function of two int params to swap values in specified array
```go
arr := []os.File{f1,f2,f3}
swapFiles := SWAPPER(arr)
swapFiles(1,2)
```

### TYPEOK(interface type conversion) bool
Returns just ok part from conversion, used for checking interface type
```go
v := any(os.File{})
if TYPEOK(v.(os.File)) { ... }
```

### SEND(ctx, chan, val)
Send to unbuffered chan, exit if context canceled
```
go func() {SEND(ctx, chan, "value"); print("continue execution")}()
cancel()
```

### NBSEND(chan, val) bool
Send to unbuffered chan, nonblocking
```
if NBSEND(chan, "value") { ... }

```

### COPYCHAN(chan) chan
Creates copy of chan, all events put in base chan will be copyed to copy. All chan events will be handled properly.
If copy is closed there must stop copying routine, if original chan will be closed all copies will be closed.
```go
func serveClient(original chan T) {
	c, d := COPYCHAN(original)
	defer d()
    w.UseChannel(c)
    ...
}
```


## General purpose functions
Set of functions, you always want to have.

### ToJson(any)string
Returns json representation of argument or dies.
```go
jsonPeople := MAP(people, func(p Person) string{ return ToJson(p) })
```
`jsonPeople` becomes slice of strings contained json representation of `people` array elements.

### ExecuteOnAllFields(ctx, storage, "method_name") error
Executes `method_name(ctx)` on all non nil interface fields in storage, used to initialize application.
```go
app := ... {
	Connection: fancy.NewConnection(cfg)
}
ExecuteOnAllFields(ctx, app, "InitStage1")
```

Requests and pull requests are [welcome](https://github.com/Vany/pirog/issues).
