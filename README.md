[![GoDoc](https://godoc.org/github.com/alexflint/go-cloudfile?status.svg)](https://godoc.org/github.com/alexflint/go-cloudfile)
[![Build Status](https://travis-ci.org/alexflint/go-cloudfile.svg?branch=master)](https://travis-ci.org/alexflint/go-cloudfile)
[![Coverage Status](https://coveralls.io/repos/alexflint/go-cloudfile/badge.svg?branch=master&service=github)](https://coveralls.io/github/alexflint/go-cloudfile?branch=master)

# cloudfile

A consistent way to work with remote (and local) files.

```go
func main() {
	url := "s3://my-bucket/path/to/file"
	// url := "http://example.com/path/to/file"
	// url := "/path/to/file"

	r, err := cloudfile.Open(url)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	...
}
```

Use `cloudfile` to write code that can doesn't care whether its data comes from
the local filesystem or a remote storage provider. All the information necessary
to locate a file is specified as a URI string.

### Motivation

Suppose you're writing a server that loads a data file. In the first version you
just load the data from the local filesystem:

    ./myserver --data /path/to/data

Later you want to load the data from an S3 bucket. Typically this would
require rewriting a bunch of code to load data using the AWS API rather than `os.Open`,
after which you would no longer be able to load from a local file without changing
the code back.

Using `cloudfile.Open`, you would simply pass in S3 URL and there would be no
need for any code changes:

    ./myserver --data s3://my-bucket/path/to/data

## Backends

`cloudfile` currently supports the following backends:

### Amazon S3

URLs look like `s3://BUCKET/KEY`

Credentials can be specified by:

 1. setting the `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables

 2. creating `~/.aws/credentials` with the following contents:

	```
    [default]
    aws_access_key_id=YOUR_ACCESS_KEY
	aws_secret_access_key=YOUR_SECRET_KEY
	```

 3. initializing the driver explicitly

	```go
	import (
		"github.com/alexflint/go-cloudfile"
		"github.com/alexflint/go-cloudfile/s3file"
	)

	func main() {
		cloudfile.Drivers["s3:"] = s3file.NewDriver("ACCESS_KEY", "SECRET_KEY", aws.USWest)
		...
	}
	```

### HTTP and HTTPS

Uses standard HTTP URLs: `https://example.com/path/to/resource`.

The resource is fetched with a `GET` using the default HTTP client, `http.DefaultClient`. HTTP resources are read-only.

### Local paths

Any path that does not have a known protocol prefix (`"s3:"`, `"http:"`, `"https:"`, etc) is 
interpreted as a local path and behavior will be equivalent to `os.Open`, `ioutil.ReadFile`, 
or `ioutil.WriteFile`.

### Custom backends

You can define your own backend by implementing `cloudfile.Driver`:

```go
type InMemoryDriver map[string][]byte

func (d InMemoryDriver) Open(path string) (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewBuffer(d[path])), nil
}

func (d InMemoryDriver) ReadFile(path string) ([]byte, error) {
	return d[path], nil
}

func (d InMemoryDriver) WriteFile(path string, buf []byte) error {
	d[path] = buf
	return nil
}
```

Register the driver with:
```go
cloudfile.Drivers["inmemory:"] = make(InMemoryDriver)
```

Then use it with:
```go
cloudfile.WriteFile("inmemory:foo", []byte("contents of file"))
buf, err := cloudfile.ReadFile("inmemory:foo")
if err != nil {
	fmt.Println(err)
	return
}
fmt.Println(string(buf))
```
