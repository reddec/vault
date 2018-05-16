package main

import (
	"github.com/jessevdk/go-flags"
	"os"
	"github.com/reddec/vault"
	"fmt"
	"github.com/reddec/vault/utils"
)

var config struct {
	URL      []string `short:"u" long:"url" env:"URL" description:"storage url"`
	URLFiles []string `short:"U" long:"url-file" env:"URL_FILES" description:"file with lines of storage URLs"`
	Args struct {
		ID string
	} `positional-args:"yes" required:"yes"`
}

var (
	version     = "dev"
	commit      = "none"
	date        = "unknown"
	name        = "vault-pull"
	description = `
Fetch object from net and write it to stdout
`
)

func main() {
	parser := flags.NewParser(&config, flags.Default)
	parser.ShortDescription = fmt.Sprint(name, " part of VAULT tools - distributes master-less object storage")
	parser.LongDescription = fmt.Sprint("version: ", version, " commit: ", commit, " date: ", date, "\n", description)

	_, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}

	wr := vault.SimpleWarehouse{}
	for _, url := range config.URL {
		wr.Add(url)
	}

	err = utils.ReadUrlFiles(&wr, config.URLFiles)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read urls file: %v\n", err)
		os.Exit(2)
	}

	network := vault.SimpleNet(&wr)

	amount, err := vault.Read(network, config.Args.ID, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read: %v\n", err)
		os.Exit(3)
	}
	fmt.Fprintln(os.Stderr, "read", amount, "bytes")
}
