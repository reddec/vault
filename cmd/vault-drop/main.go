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
	name        = "vault-drop"
	description = `
Remove object from network by ID. 
There is no guarantee that object will be deleted when not all nodes reachable.
Useful for regular cleaning
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

	err = vault.Delete(network, config.Args.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to delete: %v\n", err)
		os.Exit(2)
	}
	fmt.Fprintln(os.Stderr, "dropped")
}
