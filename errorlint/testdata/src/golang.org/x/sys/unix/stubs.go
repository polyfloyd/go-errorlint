package unix

import "syscall"

const (
	EPERM  = syscall.Errno(0x1)
	ENOENT = syscall.Errno(0x2)
)

func Rmdir(string) error {
	return ENOENT
}

func Kill(int, syscall.Signal) error {
	return EPERM
}
