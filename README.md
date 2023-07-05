![Test](https://github.com/Vany/pirog/actions/workflows/test.yml/badge.svg)
# pirog
Golang mapreduce primitives and other cool stuff from perl and javascript.

Main idea is to use commonly known and well proven constructions introduced in other languages.
Constructions have distinguishable appearance and pretend to be a part of the language rather than just functions.

## Useage
Just import it as is, but if you want and have enough dare, use special symbol to package reference like
```go
import π "github.com/Vany/pirog.git"
```
Then just use it. 

### MAP(array, callback) array
This is part of mapreduce and almost full copy of perl's map. It transforms input array to output array with callback function.
```go
type Person struct {
    FirstName  string
    SecondName string
}

persons := []Person{{"Vasya", "Pupkin"}, {"Vasya", "Lozhkin"}, {"Salvador", "Dalí"}}
out := π.MAP(persons, func(p Person) string{
	return fmt.Sprintf("%s %s", p.FirstName, p.SecondName)
})
```
`out` now is []string containing concatenated names.

### GREP(array, callback) array
This is filter, that leaves only that elements that trigerrs callback function to return true
```go
fakePersons := π.GREP(out, func(s string) bool {
    return strings.HasSuffix(s, "Pupkin")
})
```
`fakePersons` now is []string and contains just `"Vasya Pupkin"`

### KEYS(map) array
Returns just full set of keys from map to use it further
```go
artistsMap := map[string]string{
	"Vasya":"Lozhkin",
	"Salvador":"Dalí",
}
AllLozhkins := π.GREP(π.KEYS(artistsMap), func(in string) string }{
	return if artistsMap[in] == "Lozhkin" 
})
```
`AllLozhkins` will be `[]string{"Vasya"}`


### HAVEKEY(map, key) bool
Just indicates do we have key in map, or no.

### ANYKEY(map) key
Just returns any arbitrary key from map.

### MUST(err)
Validates err for nil and panics other way. When you in CLI or sure that here can not be an error.
```go
π.MUST(json.NewEncoder(buff).Encode(object))

```

## General purpose functions
Set of functions, you always want to have.

### ToJson(any)string
Returns json representation of argument or dies.
```go
jsonPersons := π.MAP(persons, func(p Person) string{ return π.ToJson(p) })
```
`jsonPersons` becomes slice of strings contained json representation of `persons` array elements.


Requests and pull requests are [welcome](https://github.com/Vany/pirog/issues).