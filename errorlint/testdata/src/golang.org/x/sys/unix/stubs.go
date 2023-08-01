package unix

const ENOENT = err(iota)

type err int

func (err) Error() string {
	return "unix err"
}

func Rmdir(string) error {
	return ENOENT
}
