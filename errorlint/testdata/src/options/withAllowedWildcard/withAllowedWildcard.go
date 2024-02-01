package testdata

import (
	"fmt"

	"example.com/pkg"
)

func Magic() {
	err := pkg.MagicOne()
	if err != pkg.ErrMagicOne {
		fmt.Println(err)
	}
}
