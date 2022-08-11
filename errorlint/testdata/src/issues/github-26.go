package issues

import "fmt"

func main() {
	err := fmt.Errorf("%v %#[1]v", struct{ string }{})
	fmt.Println(err)
}
