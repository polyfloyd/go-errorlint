package testdata

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
)

func CompareErrIndirect(r io.Reader) {
	var buf [4096]byte
	_, err := r.Read(buf[:])
	eof := io.EOF

	// Do not bother to check for comparing to aliased std errors, I have never
	// seen any code that does this in the wild. This makes the checker a bit simpler to write.
	// Supporting this use case is acceptable, patches welcome :)
	if err == eof { // want `comparing with == will fail on wrapped errors. Use errors.Is to check for a specific error`
		fmt.Println(err)
	}
}

func CompareAssignIndirect(r io.Reader) {
	var buf [4096]byte
	_, err1 := r.Read(buf[:])
	err2 := err1
	err3 := err2
	if err3 == io.EOF {
		fmt.Println(err3)
	}
}

func CompareInline(db *sql.DB) {
	var i int
	row := db.QueryRow(`SELECT 1`)
	if row.Scan(&i) == sql.ErrNoRows {
		fmt.Println("no rows!")
	}
}

func IoReadEOF(r io.Reader) {
	var buf [4096]byte
	_, err := r.Read(buf[:])
	if err == io.EOF {
		fmt.Println(err)
	}
}

func OsFileReadEOF(fd *os.File) {
	var buf [4096]byte
	_, err := fd.Read(buf[:])
	if err == io.EOF {
		fmt.Println(err)
	}
}

func IoPipeWriterWrite(w *io.PipeWriter) {
	var buf [4096]byte
	_, err := w.Write(buf[:])
	if err == io.ErrClosedPipe {
		fmt.Println(err)
	}
}

func IoReadAtLeast(r io.Reader) {
	var buf [4096]byte
	_, err := io.ReadAtLeast(r, buf[:], 8192)
	if err == io.ErrShortBuffer {
		fmt.Println(err)
	}
	if err == io.ErrUnexpectedEOF {
		fmt.Println(err)
	}
}

func IoReadFull(r io.Reader) {
	var buf [4096]byte
	_, err := io.ReadFull(r, buf[:])
	if err == io.ErrUnexpectedEOF {
		fmt.Println(err)
	}
}

func SqlRowScan(db *sql.DB) {
	var i int
	row := db.QueryRow(`SELECT 1`)
	err := row.Scan(&i)
	if err == sql.ErrNoRows {
		fmt.Println("no rows!")
	}
}

// https://github.com/polyfloyd/go-errorlint/issues/13
type CompressedFile struct {
	reader  *bytes.Reader
	zipPath string
}

func (c CompressedFile) Read(p []byte) (int, error) {
	n, err := c.reader.Read(p)
	if err == io.EOF {
		return n, io.EOF
	}
	return n, fmt.Errorf("can't read from reader")
}

func HTTPErrServerClosed(s *http.Server) error {
	if err := s.Serve(nil); err != http.ErrServerClosed {
		return err
	}
	if err := s.ServeTLS(nil, "", ""); err != http.ErrServerClosed {
		return err
	}
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	if err := s.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
		return err
	}
	if err := http.Serve(nil, nil); err != http.ErrServerClosed {
		return err
	}
	if err := http.ServeTLS(nil, nil, "", ""); err != http.ErrServerClosed {
		return err
	}
	if err := http.ListenAndServe("", nil); err != http.ErrServerClosed {
		return err
	}
	if err := http.ListenAndServeTLS("", "", "", nil); err != http.ErrServerClosed {
		return err
	}
	return nil
}
