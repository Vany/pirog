![Test](https://github.com/Vany/pirog/actions/workflows/test.yml/badge.svg)
# pirog
Golang mapreduce primitives and other cool stuff from perl and javascript.

Main idea is to use commonly known and well proven constructions introduced in other languages.
Constructions have distinguishable appearance and pretend to be a part of the language rather than just functions.

## Useage
Just import it as is, but if you want and have enough dare, use special symbol to package reference like
```go
import . "github.com/Vany/pirog.git"
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

### KEYS(map) array
Returns just full set of keys from map to use it further
```go
artistsMap := map[string]string{
	"Vasya":"Lozhkin",
	"Salvador":"Dalí",
}
AllLozhkins := GREP(KEYS(artistsMap), func(in string) string }{
	return if artistsMap[in] == "Lozhkin" 
})
```
`AllLozhkins` will be `[]string{"Vasya"}`


### HAVEKEY(map, key) bool
Just indicates do we have key in map, or no.

### ANYKEY(map) key
Returns any arbitrary key from map.

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

## General purpose functions
Set of functions, you always want to have.

### ToJson(any)string
Returns json representation of argument or dies.
```go
jsonPeople := MAP(people, func(p Person) string{ return ToJson(p) })
```
`jsonPeople` becomes slice of strings contained json representation of `people` array elements.


Requests and pull requests are [welcome](https://github.com/Vany/pirog/issues).
