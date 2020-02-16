package main

import (
	"flag"
	"fmt"
)

func main() {
	fmt.Printf("EPubFixer v1\n")
	opts := ParseFlags()
	if opts.Help {
		flag.Usage()
		return
	}
	var err error
	for _, filename := range opts.Args {
		err = ProcessFile(filename)
		if err != nil {
			fmt.Printf("Error while fixing '%s': %s", filename, err)
		}
	}
}
