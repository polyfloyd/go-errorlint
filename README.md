go-errorlint
============

[![Build Status](https://github.com/polyfloyd/go-errorlint/workflows/CI/badge.svg)](https://github.com/polyfloyd/go-errorlint/actions)

go-errorlint is a source code linter for Go software that can be used to find
code that will cause problems with the error wrapping scheme introduced in Go
1.13.

Error wrapping allows for extra context in errors without sacrificing type
information about the error's cause.

For details on Go error wrapping, see: https://golang.org/pkg/errors/


## Usage
go-errorlint accepts a set of package names similar to golint:
```
go-errorlint ./...
```
If there are one or more results, the exit status is set to `1`.


## Examples

### fmt.Errorf wrapping verb
This lint is disabled by default. Use the `-errorf` flag to toggle.
```go
// bad
fmt.Errorf("oh noes: %v", err)
// ^ non-wrapping format verb for fmt.Errorf. Use `%w` to format errors

// good
fmt.Errorf("oh noes: %w", err)
```

**Caveats**:
* When using the `-errorf` lint, keep in mind that any errors wrapped by
  `fmt.Errorf` implicitly become part of your API as according to [Hyrum's
  Law](https://github.com/dwmkerr/hacker-laws#hyrums-law-the-law-of-implicit-interfaces).

### Comparisons of errors
This lint is enabled by default. Use the `-comparison` flag to toggle.
```go
// bad
err == ErrFoo
// ^ comparing with == will fail on wrapped errors. Use errors.Is to check for a specific error

// bad
switch err {
case ErrFoo:
}
// ^ switch on an error will fail on wrapped errors. Use errors.Is to check for specific errors

// good
errors.Is(err, ErrFoo)
```

Errors returned from standard library functions that explicitly document that
an unwrapped error is returned are allowed by the linter. Notable cases are
`io.EOF` and `sql.ErrNoRows`.

**Caveats**:
* Comparing the error returned from `(io.Reader).Read` to `io.EOF` without
  `errors.Is` is considered valid as this is
  [explicitly documented](https://golang.org/pkg/io/#Reader) behaviour.
  However, nothing stops 3rd party implementations from still wrapping
  `io.EOF`, causing this linter to not detect such cases.

### Type assertions of errors
This lint is enabled by default. Use the `-asserts` flag to toggle.
```go
// bad
myErr, ok := err.(*MyError)
// ^ type assertion on error will fail on wrapped errors. Use errors.As to check for specific errors

// bad
switch err.(type) {
case *MyError:
}
// ^ type switch on error will fail on wrapped errors. Use errors.As to check for specific errors

// good
var me MyError
ok := errors.As(err, &me)
```
