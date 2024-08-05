package testdata

import (
	"archive/tar"
	"bytes"
	"context"
	"database/sql"
	"debug/elf"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
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

func CompareAssignMultiple(a, b io.Reader) error {
	var err error
	var buf []byte
	_, err = a.Read(buf)
	_, err = b.Read(buf)
	if err == io.EOF {
		return io.EOF
	}
	return fmt.Errorf("can't read from reader")
}

func CompareAssignMultipleWithUnsafe(a, b io.Reader) error {
	var err error
	var buf []byte
	_, err = a.Read(buf)
	_, err = b.Read(buf)
	err = errors.New("asdf")
	if err == io.EOF { // want `comparing with == will fail on wrapped errors. Use errors.Is to check for a specific error`
		return io.EOF
	}
	return fmt.Errorf("can't read from reader")
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

func IoReadCloserEOF(r io.ReadCloser) {
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
	if err == io.EOF {
		fmt.Println(err)
	}
	if err == io.ErrUnexpectedEOF {
		fmt.Println(err)
	}
}

func IoReaderAt(r io.ReaderAt) {
	var buf [4096]byte
	_, err := r.ReadAt(buf[:], 0)
	if err == io.EOF {
		fmt.Println(err)
	}
}

func IoLimitedReader(r *io.LimitedReader) {
	var buf [4096]byte
	_, err := r.Read(buf[:])
	if err == io.EOF {
		fmt.Println(err)
	}
}

func IoSectionReader(r *io.SectionReader) {
	var buf [4096]byte
	_, err := r.Read(buf[:])
	if err == io.EOF {
		fmt.Println(err)
	}
	_, err = r.ReadAt(buf[:], 0)
	if err == io.EOF {
		fmt.Println(err)
	}
}

func ElfOpen() {
	_, err := elf.Open("file")
	if err == io.EOF {
		fmt.Println(err)
	}
}

func ElfNewFile(r io.ReaderAt) {
	_, err := elf.NewFile(r)
	if err == io.EOF {
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

func TarHeader(r io.Reader) {
	reader := tar.NewReader(r)
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		_ = header
	}
}

func CompareUnixErrors() {
	if err := unix.Rmdir("somepath"); err != unix.ENOENT {
		fmt.Println(err)
	}
	if err := unix.Kill(1, syscall.SIGKILL); err != unix.EPERM {
		fmt.Println(err)
	}
}

func ContextErr(ctx context.Context) error {
	if err := ctx.Err(); err == context.DeadlineExceeded {
		return err
	}
	if err := ctx.Err(); err == context.Canceled {
		return err
	}
	return nil
}

func JSONReader(r io.Reader) {
	err := json.NewDecoder(r).Decode(nil)
	if err == io.EOF {
		fmt.Println(err)
	}
}

func CSVReader(r io.Reader) {
	_, err := csv.NewReader(r).Read()
	if err == io.EOF {
		fmt.Println(err)
	}
}

func MIMEMultipartReader(r io.Reader, boundary string, raw bool) {
	var err error
	if raw {
		_, err = multipart.NewReader(r, boundary).NextRawPart()
	} else {
		_, err = multipart.NewReader(r, boundary).NextPart()
	}
	if err == io.EOF {
		fmt.Println(err)
	}
}

func MIMEMultipartReadFrom(r io.Reader, boundary string, maxMemory int64) {
	_, err := multipart.NewReader(r, boundary).ReadForm(maxMemory)
	if err == multipart.ErrMessageTooLarge {
		fmt.Println(err)
	}
}

func MimeParseMediaType(contentType string) {
	_, _, err := mime.ParseMediaType(contentType)
	if err == mime.ErrInvalidMediaParameter {
		fmt.Println(err)
	}
}

func SyscallErrors() {
	err := syscall.Chmod("/dev/null", 0666)
	if err == syscall.EINVAL {
		fmt.Println(err)
	}
}
