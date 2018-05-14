package storages

import (
	"sync"
	"os"
	"strconv"
	"net/url"
	"github.com/pkg/errors"
)

type Dummy struct {
	id   string
	lock sync.RWMutex
	data map[string][]byte
}

func (d *Dummy) Put(uid string, data []byte) error {
	cp := make([]byte, len(data))
	copy(cp, data)
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.data == nil {
		d.data = make(map[string][]byte)
	}
	_, exists := d.data[uid]
	if exists {
		return errors.Errorf("%v already exists", uid)
	}
	d.data[uid] = cp
	return nil
}
func (d *Dummy) Del(uid string) error {
	if d.data == nil {
		return nil
	}
	d.lock.Lock()
	defer d.lock.Unlock()
	delete(d.data, uid)
	return nil
}

func (d *Dummy) Get(uid string) ([]byte, error) {
	if d.data == nil {
		return nil, os.ErrNotExist
	}
	d.lock.RLock()
	defer d.lock.RUnlock()
	data, ok := d.data[uid]
	if !ok {
		return nil, os.ErrNotExist
	}
	cp := make([]byte, len(data))
	copy(cp, data)
	return cp, nil
}

func (d *Dummy) List() ([]string, error) {
	if d.data == nil {
		return nil, nil
	}
	var ans []string
	d.lock.RLock()
	defer d.lock.RUnlock()
	for k := range d.data {
		ans = append(ans, k)

	}
	return ans, nil
}

func (d *Dummy) URI() string {
	return "dummy://" + d.id
}

var dummyCache struct {
	cached map[string]*Dummy
	lock   sync.Mutex
}

func NewDummy() *Dummy {
	dummyCache.lock.Lock()
	defer dummyCache.lock.Unlock()
	if dummyCache.cached == nil {
		dummyCache.cached = make(map[string]*Dummy)
	}
	id := "auto" + strconv.Itoa(len(dummyCache.cached))
	dm := &Dummy{id: id}
	dummyCache.cached[id] = dm
	return dm
}
func RemoveCachedDummy(id string) {
	if dummyCache.cached == nil {
		return
	}
	dummyCache.lock.Lock()
	defer dummyCache.lock.Unlock()
	delete(dummyCache.cached, id)
}
func init() {
	Register("dummy", func(uri *url.URL) []Storage {
		id := uri.Host
		dummyCache.lock.Lock()
		defer dummyCache.lock.Unlock()
		dm, ok := dummyCache.cached[id]
		if ok {
			return []Storage{dm}
		}
		if dummyCache.cached == nil {
			dummyCache.cached = make(map[string]*Dummy)
		}
		dm = &Dummy{id: id}
		dummyCache.cached[id] = dm
		return []Storage{dm}
	})
}
