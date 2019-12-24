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

**Caveats**:
* When using the `-errorf` lint, keep in mind that any errors wrapped by
  `fmt.Errorf` implicitly become part of your API as according to [Hyrum's
  Law](https://github.com/dwmkerr/hacker-laws#hyrums-law-the-law-of-implicit-interfaces).


## Examples

### fmt.Errorf wrapping verb
This lint must be enabled with the `-errorf` flag.
```go
// bad
fmt.Errorf("oh noes: %v", err)
// ^ non-wrapping format verb for fmt.Errorf. Use `%w` to format errors

// good
fmt.Errorf("oh noes: %w", err)
```

### Comparisons of errors
```go
// bad
err == ErrFoo
// ^ comparing with == will fail on wrapped errors. Use errors.Is to check for a specific error

// good
errors.Is(err, ErrFoo)
```

Switch statements are also checked.

### Type assertions of errors
```go
// bad
myErr, ok := err.(*MyError)
// ^ type assertion on error will fail on wrapped errors. Use errors.As to check for specific errors

// good
var me MyError
ok := errors.As(err, &me)
```

Type switch statements are also checked.
