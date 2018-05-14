package storages

import "testing"

func testStorage_Basic(stor Storage, t *testing.T) {
	err := stor.Put("a/b/c", []byte("data"))
	if err != nil {
		t.Fatal("put:", err)
	}

	data, err := stor.Get("a/b/c")
	if err != nil {
		t.Fatal("get:", err)
	}

	if string(data) != "data" {
		t.Errorf("save 'data' != '%v'", string(data))
	}

}

func TestBasicDummy(t *testing.T) {
	testStorage_Basic(NewDummy(), t)
}
