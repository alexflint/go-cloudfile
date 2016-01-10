package s3file

import (
	"errors"
	"io"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

var (
	defaultRegion      = aws.USWest
	defaultContentType string
	defaultPerm        s3.ACL
)

var (
	ErrInvalidScheme = errors.New("invalid scheme for S3 path (should be 's3')")
	ErrInvalidRegion = errors.New("invalid AWS region specified")
)

// Check for an AWS_REGION environment variable during startup
func init() {
}

// DriverOptions represents options for creating a driver
type DriverOptions struct {
	AccessKey string
	SecretKey string
	Region    aws.Region
}

// Driver implements cloudfile.Driver
type Driver struct {
	Region aws.Region
	Auth   aws.Auth
}

// NewDriverFromEnv creates a driver for S3 paths using credentials from environment variables
func NewDriverFromEnv() (*Driver, error) {
	region := defaultRegion
	if regionName := os.Getenv("AWS_REGION"); regionName != "" {
		if r, found := aws.Regions[regionName]; found {
			region = r
		} else {
			log.Printf("Invalid AWS_REGION '%s'. Ignoring.")
		}
	}

	return NewDriver("", "", region)
}

// NewDriver creates a driver for S3 paths
func NewDriver(accessKey, secretKey string, region aws.Region) (*Driver, error) {
	// Authenticate -- will fall back to ~/.aws then to environment variables
	auth, err := aws.GetAuth(accessKey, secretKey)
	if err != nil {
		return nil, err
	}

	return &Driver{
		Region: region,
		Auth:   auth,
	}, nil
}

// resolve returns the bucket corresponding to the host portion of a path
func (d *Driver) resolve(URL string) (*s3.Bucket, string, error) {
	url, err := url.Parse(URL)
	if err != nil {
		return nil, "", err
	}
	if url.Scheme != "s3" {
		return nil, "", ErrInvalidScheme
	}

	// Note that s3.New doesn't do any real work
	client := s3.New(d.Auth, d.Region)

	// S3 keys don't include the leading "/" in the URI
	path := strings.TrimPrefix(url.Path, "/")
	return client.Bucket(url.Host), path, nil
}

// Open returns a reader that reads the given S3 path (e.g. "s3://my-bucket/path/to/myfile")
func (d *Driver) Open(url string) (io.ReadCloser, error) {
	bucket, path, err := d.resolve(url)
	if err != nil {
		return nil, err
	}
	return bucket.GetReader(path)
}

// ReadFile reads from the given S3 path (e.g. "s3://my-bucket/path/to/myfile")
func (d *Driver) ReadFile(url string) ([]byte, error) {
	bucket, path, err := d.resolve(url)
	if err != nil {
		return nil, err
	}
	return bucket.Get(path)
}

// WriteFile writes to the given S3 path (e.g. "s3://my-bucket/path/to/myfile")
func (d *Driver) WriteFile(url string, buf []byte) error {
	bucket, path, err := d.resolve(url)
	if err != nil {
		return err
	}
	return bucket.Put(path, buf, defaultContentType, defaultPerm)
}
