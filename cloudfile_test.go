package cloudfile

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockDriver struct {
	lastOpen, lastRead, lastWrite string
}

func (d *mockDriver) Open(url string) (io.ReadCloser, error) {
	d.lastOpen = url
	return nil, nil
}

func (d *mockDriver) ReadFile(url string) ([]byte, error) {
	d.lastRead = url
	return nil, nil
}

func (d *mockDriver) WriteFile(url string, buf []byte) error {
	d.lastWrite = url
	return nil
}

func TestConnector(t *testing.T) {
	mock := &mockDriver{}
	Connectors["foo:"] = func() (Driver, error) {
		return mock, nil
	}
	Open("foo:bar")
	assert.Equal(t, "foo:bar", mock.lastOpen)
}

func TestOpen(t *testing.T) {
	mock := &mockDriver{}
	Drivers["foo:"] = mock
	Open("foo:bar")
	assert.Equal(t, "foo:bar", mock.lastOpen)
}

func TestReadFile(t *testing.T) {
	mock := &mockDriver{}
	Drivers["foo:"] = mock
	ReadFile("foo:bar")
	assert.Equal(t, "foo:bar", mock.lastRead)
}
func TestWriteFile(t *testing.T) {
	mock := &mockDriver{}
	Drivers["foo:"] = mock
	WriteFile("foo:bar", nil)
	assert.Equal(t, "foo:bar", mock.lastWrite)
}
