package main

import (
	"github.com/jessevdk/go-flags"
	"os"
	"github.com/reddec/vault"
	"fmt"
)

var config struct {
	URL []string `short:"u" long:"url" env:"URL" description:"storage url"`
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

	err = vault.Delete(network, config.Args.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to delete: %v\n", err)
		os.Exit(2)
	}
	fmt.Fprintln(os.Stderr, "dropped")
}
