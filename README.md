# cloudfile

A consistent way to work with remote files.

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

Use `cloudfile` to write code that can doesn't care whether its data comes from
local or cloud storage. All the information necessary to locate
a file is specified as a string.

For example, suppose you're writing a server that loads a data file. 
In the first implementation you load the data from a file on
the local filesystem:

    ./myserver --data /path/to/data

Later you want to load the data from an S3 bucket. Typically this would
require rewriting a bunch of code to load data using the AWS API, after which you
would no longer be able to load from a local file without changing the code back.

Using `cloudfile`, you would simply pass in S3 URL and there would be no
need for any code changes:

    ./myserver --data s3://my-bucket/path/to/data

## Backends

`cloudfile` currently supports the following backends:

### Amazon S3

URLs look like `s3://BUCKET/KEY`

Credentials can be specified in one of the following ways:

 1. by setting the `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables

 2. by creating `~/.aws/credentials` with the following contents:

     [default]
     aws_access_key_id=YOUR_ACCESS_KEY
	 aws_secret_access_key=YOUR_SECRET_KEY

 3. by assigning to global variables:

     import "github.com/alexflint/go-cloudfile/s3file"

     ...

     s3file.AccessKey = "YOUR_ACCESS_KEY"
     s3file.SecretKey = "YOUR_SECRET_KEY"

### HTTP[S]

Uses standard HTTP URLs: `https://example.com/path/to/resource`. The resource is fetched with a `GET` using the default HTTP client, `http.DefaultClient`. HTTP resources are read-only.

### Local paths

Any path that doesn't begin with a known prefix (`s3:`, `http:`, `https:`, etc) is 
interpreted as a local path and behavior will be equivalent to `os.Open`, `ioutil.ReadFile`, 
or `ioutil.WriteFile`.
