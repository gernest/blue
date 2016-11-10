

#blue.Line ![blue logo](logo-512.png)

[![GoDoc](https://godoc.org/github.com/gernest/blue?status.svg)](https://godoc.org/github.com/gernest/blue) [![Go Report Card](https://goreportcard.com/badge/github.com/gernest/blue)](https://goreportcard.com/report/github.com/gernest/blue)  [![Build Status](https://travis-ci.org/gernest/blue.svg?branch=master)](https://travis-ci.org/gernest/blue) [![Coverage Status](https://coveralls.io/repos/github/gernest/blue/badge.svg?branch=master)](https://coveralls.io/github/gernest/blue?branch=master)

Generates influxdb line protocol from json objects.


**WARNING**:Do not use this in production

# Features

* Flexible:Allow custom functions to choose
   - The name of the measurement
   - Tags
   - Fields
   - Timestamps

* Works on arbitrary json input.
* Tags and Fields are sorted. This is per recommendation from influxdb
  documentation.


# The steps taken when processing the json input

## Flattening
The json object is flattened to key, value pair.This is important because
influxdb line protocol is based on key, value semantics.

This brings up the problem of nested objects like arrays of objects etc. The
best strategy I could come up with is to join the keys that point to a value.

For example
```json
{
  "top":{
    "level":1
  }
}
```

That can easily be flattered to `top_level=1` , where by we used underscore `_` to
join the keys. This rises another problem though maybe the
`names_of_keys_are_like_this` so to make it flexible the user of this library
can specify the joining strategy of the keys, whereby you can define a function
of the form `func(a,b string)string` that will be used to join keys.


## Filtering

After the json bject has been flattered to key, value pairs it is filtered to
pic what the user thinks is important. The main components are,

- Measurement name
- Tags
- Fields

The user of this library, has control on all the mentioned components. The
Options struct looks like this.

```go
type Options struct {
	//This is the function that joins the keys when flattering the json object.
	//The first argument is the top level key(although this might be not the
	//case for deeply nested objects) and the second is the current key.
	//
	// The returned string is the key that will be used. This is implementation
	// specific, you can do whatever the hell you want with this.
	KeyJoinFunc func(a, b string) string

	// Checks  if the key is a tag. If the second returned value is true then
	// the first returned value is going to be used as key.
	//
	// Something to consider here the key might be a joint keys for the deeply
	// nested objects. The keys will be supposedly joined by the KeyJoinFunc
	// implementation.
	IsTag func(string) (string, bool)

	// Checks if the aregumen( a key) belongs to a Field. Implementations should
	// return true if the key is a field and false otherwise.
	IsField func(string) bool

	//IsMeasurement this determines if the given key is a key with the
	//measurement name.
	//
	//It is up to the implementaion to return the measurement name, true. If the
	//second returned value is false then the key is assumend to be not
	//measurement name.
	//
	// If the Measurement field is set i.e not empty then this function is nver
	// going to be used.
	IsMeasurement func(string, interface{}) (string, bool)

	// Determines the timestamp of the measurement. Timestamp is set only once.
	IsTimeStamp func(string, interface{}) (time.Time, bool)
	Measurement string
}
```

Specifying the Measurement field will make the IsMeasurent function to be
ignored.

# Installation
```shell
go get github.com/gernest/blue
```


# Usage
```go

package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gernest/blue"
)

const src = `
{
	"timestamp":1434055562000000000,
	"tags":{
		"noun":"broke",
		"verb":"bankrupt"
	},
	"anger_level":80,
	"this":{
		"is":{
			"nested":{
				"value": 1
			}
		}
	}
}
`

func main() {
	o := blue.Options{
		KeyJoinFunc: func(a, b string) string {
			if a == "" {
				return b
			}
			return a + "_" + b
		},
		IsTag: func(key string) (string, bool) {
			prefix := "tags_"
			if strings.HasPrefix(key, prefix) {
				return strings.TrimPrefix(key, prefix), true
			}
			return "", false
		},
		IsField: func(key string) bool {
			return true
		},
		Measurement: "ordinary",
	}
	m, err := blue.Line([]byte(src), o)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(m)
	//ordinary,noun=broke,verb=bankrupt anger_level=80,this_is_nested_value=1 1434055562000000000
}
```

# Author
Geofrey Ernest [@gernesti](https://twitter.com/gernesti)

Logo  was  made by [@dr-okra](https://github.com/dr-okra)
