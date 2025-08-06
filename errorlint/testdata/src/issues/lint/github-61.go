package issues

import (
	"errors"
)

func Issue61() {
	err := errors.Join(errors.New("err1"), errors.New("err2"))
	errs, ok := err.(interface{ Unwrap() []error })
	_ = errs
	_ = ok
}
