package issues

import (
	"errors"
	"fmt"
)

var err1 = errors.New("1")

func Issue57() {
	err := err1
	var authErr error

	authErr = err
	err = authErr

	if err == err1 { // want `comparing with == will fail on wrapped errors. Use errors.Is to check for a specific error`
		fmt.Println("err1")
	}
}
