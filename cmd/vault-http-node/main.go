package main

import (
	"net/http"
	"github.com/jessevdk/go-flags"
	"os"
	"github.com/boltdb/bolt"
	"fmt"
	"io/ioutil"
)

var config struct {
	Listen string `short:"l" long:"listen" env:"LISTEN" description:"Listen address" default:":9911"`
	DbPath string `short:"d" long:"db" env:"DB" description:"Path to directory with database" default:"data"`
}

var (
	version     = "dev"
	commit      = "none"
	date        = "unknown"
	name        = "vault-http-node"
	description = `
Server simple HTTP node storage.
Keeps chunks in BoltDB.
`
)

func run() error {
	db, err := bolt.Open(config.DbPath, 0755, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()
		chunkId := request.URL.Path[1:]
		if chunkId == "" {
			var ans = make([]string, 0)
			err := db.View(func(tx *bolt.Tx) error {
				bucket := tx.Bucket([]byte("main"))
				if bucket == nil {
					return nil
				}
				c := bucket.Cursor()
				for k, _ := c.First(); k != nil; k, _ = c.Next() {
					ans = append(ans, string(k))
				}
				return nil
			})
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			writer.WriteHeader(http.StatusOK)
			for _, k := range ans {
				writer.Write([]byte(k))
				writer.Write([]byte("\n"))
			}
			return
		}

		switch request.Method {
		case http.MethodGet:
			var payload []byte

			err := db.View(func(tx *bolt.Tx) error {
				bucket := tx.Bucket([]byte("main"))
				if bucket == nil {
					return os.ErrNotExist
				}
				payload = bucket.Get([]byte(chunkId))
				if err != nil {
					return err
				}
				if payload == nil {
					return os.ErrNotExist
				}
				return nil
			})
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			writer.WriteHeader(http.StatusOK)
			writer.Write(payload)
		case http.MethodPost:
			err := db.Update(func(tx *bolt.Tx) error {
				bucket, err := tx.CreateBucketIfNotExists([]byte("main"))
				if err != nil {
					return err
				}
				payload := bucket.Get([]byte(chunkId))
				if err != nil {
					return err
				}
				if payload != nil {
					return os.ErrExist
				}
				data, err := ioutil.ReadAll(request.Body)
				if err != nil {
					return err
				}
				return bucket.Put([]byte(chunkId), data)
			})
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			writer.WriteHeader(http.StatusCreated)
		case http.MethodDelete:
			err := db.Update(func(tx *bolt.Tx) error {
				bucket := tx.Bucket([]byte("main"))
				if bucket == nil {
					return os.ErrNotExist
				}
				payload := bucket.Get([]byte(chunkId))
				if err != nil {
					return err
				}
				if payload == nil {
					return os.ErrNotExist
				}
				return bucket.Delete([]byte(chunkId))
			})
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			writer.WriteHeader(http.StatusOK)
		}

	})

	return http.ListenAndServe(config.Listen, nil)
}

func main() {
	parser := flags.NewParser(&config, flags.Default)
	parser.ShortDescription = fmt.Sprint(name, " part of VAULT tools - distributes master-less object storage")
	parser.LongDescription = fmt.Sprint("version: ", version, " commit: ", commit, " date: ", date, "\n", description)

	_, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}
	err = run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
