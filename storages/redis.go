package storages

import (
	"github.com/go-redis/redis"
	"os"
	"net/url"
	"fmt"
)

type redisProxy struct {
	client *redis.Client
}

func (rp *redisProxy) Put(uid string, data []byte) error {
	cmd := rp.client.SetNX(uid, data, 0)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	if !cmd.Val() {
		return os.ErrExist
	}
	return nil
}

func (rp *redisProxy) Del(uid string) (error) {
	return rp.client.Del(uid).Err()
}

func (rp *redisProxy) Get(uid string) ([]byte, error) {
	cmd := rp.client.Get(uid)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	data, err := cmd.Bytes()
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, os.ErrNotExist
	}
	return data, nil
}

func (rp *redisProxy) List() ([]string, error) {
	cmd := rp.client.Keys("*")
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return cmd.Val(), nil
}

func init() {
	Register("redis", func(uri *url.URL) []Storage {
		opts, err := redis.ParseURL(uri.String())
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed parse redis url", uri.String(), "-", err)
			return nil
		}
		return []Storage{&redisProxy{redis.NewClient(opts)}}
	})
}
