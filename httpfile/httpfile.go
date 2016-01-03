package httpfile

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	ErrWriteNotSupported = errors.New("WriteFile is not supported for http[s]")
)

// Driver implements cloudfile.Driver using local filesystem operations
type Driver struct{}

// Open returns a reader that reads the given resource
func (d Driver) Open(path string) (io.ReadCloser, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// ReadFile reads the entire resource
func (d Driver) ReadFile(path string) ([]byte, error) {
	r, err := d.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return ioutil.ReadAll(r)
}

// WriteFile writes the entire resource
func (d Driver) WriteFile(path string, buf []byte) error {
	return ErrWriteNotSupported
}
