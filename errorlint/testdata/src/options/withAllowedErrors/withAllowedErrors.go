package testdata

import (
	"fmt"
	"io"

	"example.com/pkg"
)

func CustomPackageCompare(r io.Reader) {
	err := pkg.Read(r)
	if err == io.EOF {
		fmt.Println(err)
	}
}
