package issues

import (
	"database/sql"
	"errors"
	"fmt"
)

type Foo struct {
	Err error
}

func Issue42(db *sql.DB) {
	var err1 error
	foo := Foo{}
	foo.Err, err1 = errors.New("foo"), errors.New("bar")

	var i int
	err := db.QueryRow(`SELECT 2`).Scan(&i)
	if err == err1 { // want `comparing with == will fail on wrapped errors. Use errors.Is to check for a specific error`
		fmt.Println("no rows!")
	}
}
