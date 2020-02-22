package main

import "flag"

type Options struct {
	Help bool
	Args []string
}

//ParseFlags parse flags
func ParseFlags() *Options {
	opts := new(Options)
	flag.BoolVar(&opts.Help, "help", false, "Print help")
	flag.Parse()
	opts.Args = flag.Args()
	return opts
}
