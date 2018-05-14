package vault

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func genWarehouse() Warehouse {
	whr := &SimpleWarehouse{}
	whr.Add("dummy://node0")
	whr.Add("dummy://node1")
	whr.Add("dummy://node2")

	//whr.Add("redis://127.0.0.1:6379")
	//whr.Add("redis://127.0.0.1:6380")
	//whr.Add("redis://127.0.0.1:6381")
	return whr
}

func TestWriteRead(t *testing.T) {
	whr := genWarehouse()
	net := &netImpl{wr: whr}
	err := net.Put("qwe", []byte("1234567890"), 2)
	assert.NoError(t, err)
	data, err := net.Get("qwe")
	assert.NoError(t, err)
	assert.Equal(t, "1234567890", string(data))

	var redundancy = 0
	for _, nd := range whr.All() {
		items, err := nd.List()
		assert.NoError(t, err)
		if len(items) > 0 {
			redundancy++
		}
	}

	assert.Equal(t, 2, redundancy)

	err = net.Sync(3)
	assert.NoError(t, err)

	redundancy = 0
	for _, nd := range whr.All() {
		items, err := nd.List()
		assert.NoError(t, err)
		if len(items) > 0 {
			redundancy++
		}
	}

	assert.Equal(t, 3, redundancy)


}
