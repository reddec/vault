package vault

import (
	"io"
	"github.com/reddec/vault/storages"
	"strconv"
	"strings"
)

const ChunkSize = 1 * 1024 * 1024 // 1 MB

type Net interface {
	Put(chunkUID string, data []byte, redundancy int) (error)
	Get(chunkUID string) ([]byte, error)
	Sync(redundancy int) (int, error)
	Del(chunkUID string) error // Even if delete returns successfully, it's not mean that no more chunks left. Reuse removed chunks NOT SAFE!!
}

type Warehouse interface {
	All() []storages.Storage
	Shuffle() []storages.Storage
}

func Write(net Net, id string, source io.Reader) (int64, error) {
	return WriteStream(net, id, 3, ChunkSize, source)
}

func WriteStream(net Net, id string, redundancy, chunkSize int, source io.Reader) (int64, error) {
	var chunk = make([]byte, chunkSize)
	var total int64
	var chunkNum int64
	for {
		size, err := source.Read(chunk)
		if size == 0 && err != nil {
			if err == io.EOF {
				break
			}
			return total, err
		}
		if size == 0 {
			continue
		}
		chunkId := chunkUid(id, chunkNum)

		err = net.Put(chunkId, chunk[:size], redundancy)
		if err != nil {
			return total, err
		}
		chunkNum++
		total += int64(size)
	}
	// final block
	err := net.Put(chunkUid(id, chunkNum), make([]byte, 0), redundancy)
	return total, err
}

func Delete(net Net, id string) error {
	var chunkNum int64
	for {
		err := net.Del(chunkUid(id, chunkNum))
		if err != nil {
			if strings.Contains(err.Error(), "exist") {
				break
			}
			return err
		}
		chunkNum++
	}
	return nil
}

func Read(net Net, id string, target io.Writer) (int64, error) {
	var total int64
	var chunkNum int64
	for {
		chunk, err := net.Get(chunkUid(id, chunkNum))
		if err != nil {
			return total, err
		}
		if len(chunk) == 0 {
			//final
			break
		}
		var offset = 0
		for offset < len(chunk) {
			written, err := target.Write(chunk[offset:])
			if err != nil {
				return total, err
			}
			offset += written
			total += int64(written)
		}
		chunkNum++
	}
	return total, nil
}

func chunkUid(id string, chunkNum int64) string { return id + "/" + strconv.FormatInt(chunkNum, 10) }
