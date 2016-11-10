

#blue.Line ![blue logo](logo-512.png)

[![GoDoc](https://godoc.org/github.com/gernest/blue?status.svg)](https://godoc.org/github.com/gernest/blue) [![Go Report Card](https://goreportcard.com/badge/github.com/gernest/blue)](https://goreportcard.com/report/github.com/gernest/blue)

Generates influxdb line protocol from json objects.


**WARNING**:
- Do not use this in production.
- There is a fail dose of interface{}
usage, but don't worry this library DOES NOT use reflection.

After saying that you should also know, I tested this on a live influxdb instance
by hooking up a telegraf tail plugin to a named pipe that another commandline
library was writing to( after converting live stream of json input to the
influxdb line protocol using this library)

# Features

* Flexible
  - Allow custom functions to choose
   - The name of the measurement
   - Tags
   - Fields
   - Timestamps

* Works on arbitrary json input.


# The steps taken when processing the json inpu

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

The user of this library, has contol on alll the mentioned components. The
Options strct looks like this.

```go
type Options struct {
	KeyJoinFunc   func(a, b string) string
	IsTag         func(string) (string, bool)
	IsField       func(string) bool
	IsMeasurement func(string, interface{}) (string, bool)
	IsTimeStamp   func(string, interface{}) (time.Time, bool)
	Measurement   string
}
```

Specifying the Measurement field will make the IsMeasurent function to be
ignored.

## Outpu

Just call the `Line()` method of the returned `*Measurement` object and influxdb
line compliant string will be generated.


# Example Usage
