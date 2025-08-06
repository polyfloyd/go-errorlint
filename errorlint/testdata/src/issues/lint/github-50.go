package issues

// Regression test for: https://github.com/polyfloyd/go-errorlint/issues/50

type AssertError struct {
	reason string
}

func (e *AssertError) Error() string {
	return e.reason
}

func (e *AssertError) Is(target error) bool {
	terr, ok := target.(*AssertError)
	return ok && terr.reason == e.reason
}

type SwitchError struct {
	reason string
}

func (e *SwitchError) Error() string {
	return e.reason
}

func (e *SwitchError) Is(target error) bool {
	switch terr := target.(type) {
	case *SwitchError:
		return terr.reason == e.reason
	default:
		return false
	}
}
