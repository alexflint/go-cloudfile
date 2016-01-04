package main

import (
	"fmt"
	"io"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/alexflint/go-cloudfile"
)

func main() {
	var args struct {
		URL string `arg:"positional"`
	}
	arg.MustParse(&args)

	r, err := cloudfile.Open(args.URL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer r.Close()

	io.Copy(os.Stdout, r)
}
