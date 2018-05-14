package storages

import "net/url"

type Storage interface {
	// put one chunk to storage. If chunk already exists - error
	Put(uid string, data []byte) error
	Get(uid string) ([]byte, error)
	Del(uid string) error // if no such chunk, delete must return error with 'exist' word
	List() ([]string, error)
}

type StorageResolver func(uri *url.URL) []Storage

var _registry = make(map[string]StorageResolver)

func Register(protocol string, resolver StorageResolver) {
	_registry[protocol] = resolver
}

func Resolve(uri string) []Storage {
	u, e := url.Parse(uri)
	if e != nil {
		return nil
	}
	resolver, ok := _registry[u.Scheme]
	if !ok {
		return nil
	}
	return resolver(u)
}
