package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/alexflint/go-cloudfile"
)

func main() {
	var args struct {
		URL string `arg:"positional"`
	}
	arg.MustParse(&args)

	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = cloudfile.WriteFile(args.URL, buf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
