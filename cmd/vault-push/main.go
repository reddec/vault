package main

import (
	"github.com/jessevdk/go-flags"
	"os"
	"github.com/reddec/vault"
	"fmt"
)

var config struct {
	URL        []string `short:"u" long:"url" env:"URLS" description:"storage url"`
	ChunkSize  int      `short:"c" long:"chunk" env:"CHUNK" description:"chunk size in kb" default:"1024"`
	Redundancy int      `short:"r" long:"redundancy" env:"REDUNDANCY" description:"number of copies" default:"3"`
	Args struct {
		ID string
	} `positional-args:"yes" required:"yes"`
}


func main() {
	_, err := flags.Parse(&config)
	if err != nil {
		os.Exit(1)
	}

	wr := vault.SimpleWarehouse{}
	for _, url := range config.URL {
		wr.Add(url)
	}

	network := vault.SimpleNet(&wr)

	amount, err := vault.WriteStream(network, config.Args.ID, config.Redundancy, config.ChunkSize*1024, os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write: %v\n", err)
		os.Exit(2)
	}
	fmt.Println("written", amount, "bytes")
}
