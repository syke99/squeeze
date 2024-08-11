# squeeze
[![Go Reference](https://pkg.go.dev/badge/github.com/syke99/squeeze.svg)](https://pkg.go.dev/github.com/syke99/squeeze)
[![Go Reportcard](https://goreportcard.com/badge/github.com/syke99/squeeze)](https://goreportcard.com/report/github.com/syke99/squeeze)

Easily flatten a slice of structs of the same type by defining and applying rules in a simple, generic, and type-safe manner

# What problem does squeeze solve?

Imagine, for example, you have a producer-consumer pattern with many producers, and one consumer. And in this hypothetical scenario,
each producer produces a `result` struct for the consumer to use. But, say you want this consumer to also be a producer of sorts, i.e.
you want it to condense all of these `result`s into a single `result` for another consumer to consume and process. Unfortunately, in Go,
this isn't exactly an easy task. It often involves a lot of messy code using the `reflect` package to loop through the fields, etc. etc.
This is where `squeeze` comes in handy. Instead, you can define functions, or `Rule`s that pertain to specific fields and tie them to the
field you want them to compare using a `json` tag of the corresponding name (`squeeze` _is_ case-sensitive). This way, you can define `Rule`s
for each field in a struct in a type-safe, generic manner. Below is a quick example.

# Example

Say you have a struct that we'll call `MyStruct` with, among others a field named `MyIntField` of type `int`:

```go
type MyStruct struct {
	MyIntField int `json:"MyIntField"`
	... // other fields
}
```

Then, say there's several possibilities that could determine what you would like the final value of `MyIntField` to be. It's as simple as
first defining your `Rule`:

```go
func whichMyIntFieldDoIUse(left, right MyStruct) squeeze.Result {
	// Implement your logic here for each field here. If you want squeeze
	// to use the value from the left struct, return squeeze.Left;
	// or, if you want squeeze to use the value from the right struct,
	// return squeeze.Right
}
```
(note: `Rule` functions must match the signature of: `func Rule[T any] (left, right T) squeeze.Result`)

Once that's complete, add it and any other field(s') `Rule`s to a `Rules` map and pass both your slice of structs and your `Rules` to
`squeeze.Squeeze()`:

```go

first := MyStruct{MyIntField: 12}
secont := MyStruct{MyIntField: 34}
third := MyStruct{MyIntField: 56}

myStructs := []MyStruct{
	first,
	second,
	third,
	...
}

rules := squeeze.Rules[MyStruct]{
	"MyIntField": whichMyIntFieldDoIUse,
	...
}

result, err := squeeze.Squeeze([]MyStruct{myStructs, rules)
```

Put it all together and it should look a little something like this:

```go
package main

import (
	"github.com/syke99/squeeze"
)

type MyStruct struct {
	MyIntField int `json:"MyIntField"`
	... // other fields
}

func whichMyIntFieldDoIUse(left, right MyStruct) squeeze.Result {
	// Implement your logic here for each field here. If you want squeeze
	// to use the value from the left struct, return squeeze.Left;
	// or, if you want squeeze to use the value from the right struct,
	// return squeeze.Right
}

func main() {
	first := MyStruct{MyIntField: 12}
	secont := MyStruct{MyIntField: 34}
	third := MyStruct{MyIntField: 56}

	myStructs := []MyStruct{
		first,
		second,
		third,
		...
	}
	
	rules := squeeze.Rules[MyStruct]{
		"MyIntField": whichMyIntFieldDoIUse,
		...
	}

	result, err := squeeze.Squeeze([]MyStruct{myStructs, rules)
}
```
