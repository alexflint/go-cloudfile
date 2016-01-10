package cloudfile

import (
	"io"
	"strings"
	"sync"

	"github.com/alexflint/go-cloudfile/httpfile"
	"github.com/alexflint/go-cloudfile/osfile"
	"github.com/alexflint/go-cloudfile/s3file"
)

var (
	connectMutex  sync.Mutex
	connectOnce   = make(map[string]sync.Once)
	connectErrors = make(map[string]error)
)

// A Driver reads and writes files
type Driver interface {
	// Open returns a reader that reads the given resource
	Open(url string) (io.ReadCloser, error)
	// Get reads the entire resource
	ReadFile(url string) ([]byte, error)
	// Put writes the entire resource
	WriteFile(url string, buf []byte) error
}

// A Connector creates a driver
type Connector func() (Driver, error)

// Connectors is a map from URL prefixes to connectors that handle them. The first time
// a URL scheme is accessed, the connector is used to construct a driver, and from then
// on only that driver is used.
var Connectors = map[string]Connector{
	"s3:": func() (Driver, error) { return s3file.NewDriverFromEnv() },
}

// Drivers is a map from URL prefixes to drivers that handle them
var Drivers = map[string]Driver{
	"http:":  httpfile.Driver{},
	"https:": httpfile.Driver{},
}

// Fallback is the driver of last resort used when no prefix is matched
var Fallback Driver = osfile.Driver{}

// drive gets the driver for a given path
func drive(url string) (Driver, error) {
	// First try the drivers
	for prefix, driver := range Drivers {
		if strings.HasPrefix(url, prefix) {
			return driver, connectErrors[prefix]
		}
	}

	// Next try the connectors
	for prefix, connect := range Connectors {
		if strings.HasPrefix(url, prefix) {
			connectMutex.Lock()
			// Must check again to see whether the driver is nil to avoid races
			if _, present := Drivers[prefix]; !present {
				Drivers[prefix], connectErrors[prefix] = connect()
			}
			connectMutex.Unlock()

			if err := connectErrors[prefix]; err != nil {
				return nil, err
			}
			return Drivers[prefix], nil
		}
	}

	// Finally fall back to the default
	return Fallback, nil
}

// Open returns a reader that reads the given resource
func Open(url string) (io.ReadCloser, error) {
	d, err := drive(url)
	if err != nil {
		return nil, err
	}
	return d.Open(url)
}

// WriteFile writes to the given resource
func ReadFile(url string) ([]byte, error) {
	d, err := drive(url)
	if err != nil {
		return nil, err
	}
	return d.ReadFile(url)
}

// ReadFile reads from the given resource
func WriteFile(url string, buf []byte) error {
	d, err := drive(url)
	if err != nil {
		return err
	}
	return d.WriteFile(url, buf)
}
