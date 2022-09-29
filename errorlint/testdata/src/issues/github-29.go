package issues

import "fmt"

func Issue29() {
	err := fmt.Errorf("%v %#[2]v", struct{ string }{})
	fmt.Println(err)
}
