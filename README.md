go-ini
======

INI parsing library for Go (golang).

This library supports read/write and implements io.ReadWriteCloser, io.ReadFrom and io.WriteTo for convenience of integration

View the API documentation [here](https://godoc.org/github.com/dombenson/go-ini).

N.B. that the current (v2) of this library is a substantial change from v1 (see https://github.com/vaughan0/go-ini). General usage remains unchanged, but support for direct access to internal structures has been dropped.

Usage
-----

Parse an INI file:

```go
import "github.com/dombenson/go-ini"

file, err := ini.LoadFile("myfile.ini")
```

Get (string) data from the parsed file:

```go
name, ok := file.Get("person", "name")
if !ok {
  panic("'name' variable missing from 'person' section")
}
```

Get (array) data from the parsed file:

```go
colours, ok := file.GetArr("apples", "colour")
if !ok {
  panic("'colours' array variable missing from 'apples' section")
}
```

Create a new file for writing:

```go
file := ini.NewFile()
```

Set a value in the file:

```go
file.Set("person", "name", "fred")
```

Write a file out:

```go
file.WriteTo(io.Writer)
```

File Format
-----------

INI files are parsed by go-ini line-by-line. Each line may be one of the following:

  * A section definition: [section-name]
  * A property: key = value
  * An array property: key[] = value
  * A comment: #blahblah _or_ ;blahblah
  * Blank. The line will be ignored.

Properties defined before any section headers are placed in the default section, which has
the empty string as it's key.

Example:

```ini
# I am a comment
; So am I!

[apples]
colour[] = red
colour[] = green
shape = applish

[oranges]
shape = square
colour = blue
```

Tests
-----
The tests in this package are written to use the subtest feature of Go 1.7. 

Attempting to run the tests with older go will yield a result like 
```
./ini_test.go:366: t.Run undefined (type *testing.T has no field or method Run)
```
See [https://github.com/mpvl/subtest](https://github.com/mpvl/subtest) if you need to run them on an older release of Go.
