package osfile

import (
	"io"
	"io/ioutil"
	"os"
)

// DefaultPerm is the default permissions for creating files
var DefaultPerm os.FileMode = 0777

// Driver implements cloudfile.Driver using local filesystem operations
type Driver struct{}

// Open returns a reader that reads the given resource
func (d Driver) Open(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

// Get reads the entire resource
func (d Driver) ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

// Put writes the entire resource
func (d Driver) WriteFile(path string, buf []byte) error {
	return ioutil.WriteFile(path, buf, DefaultPerm)
}
