package vault

import (
	"sync"
	"github.com/reddec/vault/storages"
	"math/rand"
)

type SimpleWarehouse struct {
	lock     sync.RWMutex
	storages []string
}

func (sw *SimpleWarehouse) All() []storages.Storage {
	var ans []storages.Storage
	sw.lock.RLock()
	defer sw.lock.RUnlock()
	for _, s := range sw.storages {
		ans = append(ans, storages.Resolve(s)...)
	}
	return ans
}

func (sw *SimpleWarehouse) Shuffle() []storages.Storage {
	all := sw.All()
	rand.Shuffle(len(all), func(i, j int) {
		all[i], all[j] = all[j], all[i]
	})
	return all
}

func (sw *SimpleWarehouse) Add(uri string) {
	sw.lock.Lock()
	defer sw.lock.Unlock()
	sw.storages = append(sw.storages, uri)
}

func (sw *SimpleWarehouse) Remove(uri string) {
	sw.lock.Lock()
	defer sw.lock.Unlock()
	for i, u := range sw.storages {
		if u == uri {
			sw.storages = append(sw.storages[:i], sw.storages[i+1:]...)
			break
		}
	}
}
