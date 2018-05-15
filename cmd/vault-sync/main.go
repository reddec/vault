package main

import (
	"github.com/jessevdk/go-flags"
	"os"
	"github.com/reddec/vault"
	"fmt"
)

var config struct {
	URL        []string `short:"u" long:"url" env:"URLS" description:"storage url"`
	Redundancy int      `short:"r" long:"redundancy" env:"REDUNDANCY" description:"number of copies" default:"3"`
}

var (
	version     = "dev"
	commit      = "none"
	date        = "unknown"
	name        = "vault-sync"
	description = `
Re-distribute chunks with redundancy less then specified (heal net when some nodes down)
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

	network := vault.SimpleNet(&wr)

	items, err := network.Sync(config.Redundancy)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to sync: %v\n", err)
		os.Exit(2)
	}
	fmt.Println("fixed", items, "chunks")
}
