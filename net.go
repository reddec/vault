package vault

import (
	"github.com/pkg/errors"
	"math/rand"
	"fmt"
	"os"
)

type netImpl struct {
	wr Warehouse
}

func SimpleNet(wr Warehouse) Net { return &netImpl{wr} }

func (ni *netImpl) Put(chunkUID string, data []byte, redundancy int) (error) {
	var nodes = ni.wr.All()
	var ok bool
	var i int
	for i < redundancy && len(nodes) > 0 {
		j := rand.Intn(len(nodes))
		node := nodes[j]

		// remove
		nodes[j] = nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]

		err := node.Put(chunkUID, data) // put on exists must return error
		if err == nil {
			ok = true
			i++
		}
	}
	if !ok {
		return errors.Errorf("put %v: no suitable hosts", chunkUID)
	}
	return nil
}

func (ni *netImpl) Get(chunkUID string) ([]byte, error) {
	var nodes = ni.wr.All()
	for len(nodes) > 0 {
		j := rand.Intn(len(nodes))
		node := nodes[j]

		data, err := node.Get(chunkUID)
		if err == nil {
			return data, nil
		}

		// remove
		nodes[j] = nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]

	}
	return nil, errors.Errorf("get %v: no suitable hosts", chunkUID)
}

func (ni *netImpl) Sync(redundancy int) (int, error) {
	var nodes = ni.wr.All()
	chunkIndex := map[string]int{}
	// get index
	for _, node := range nodes {
		chunks, err := node.List()
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed get index", err)
			continue
		}
		for _, chunk := range chunks {
			chunkIndex[chunk] = chunkIndex[chunk] + 1
		}
	}
	var synced int
	// redistribute chunks with less redundancy than expected
	for chunk, count := range chunkIndex {
		if count >= redundancy {
			continue
		}
		data, err := ni.Get(chunk)
		if err != nil {
			fmt.Println("failed get chunk", chunk, err)
			continue
		}

		err = ni.Put(chunk, data, redundancy-count)
		if err != nil {
			fmt.Println("failed redistribute chunk", chunk, err)
		}
		synced++
	}
	return synced, nil
}

func (ni *netImpl) Del(chunkUID string) error {
	var nodes = ni.wr.All()
	var success bool
	var lastErr error
	for _, node := range nodes {
		err := node.Del(chunkUID)
		if err != nil {
			lastErr = err
		} else {
			success = true
		}
	}
	if !success {
		return lastErr
	}
	return nil
}
