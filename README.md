go-errorlint
============

go-errorlint is a source code linter for Go software that can be used to find
code that will cause problems with the error wrapping scheme introduced in Go
1.13.

Error wrapping allows for extra context in errors without sacrificing type
information about the error's cause.

For details on Go error wrapping, see: https://golang.org/pkg/errors/


## Examples

### fmt.Errorf wrapping verb
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
