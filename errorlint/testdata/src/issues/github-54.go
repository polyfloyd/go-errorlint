package issues

import (
	"fmt"

	"golang.org/x/sys/unix"
)

func SwitchOnUnixErrors() {
	err := unix.Rmdir("somepath")
	switch err {
	case unix.ENOENT:
		return
	case unix.EPERM:
		return
	}
	fmt.Println(err)
}
