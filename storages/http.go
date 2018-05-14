package storages

import (
	"net/http"
	"net/url"
	"bytes"
	"io/ioutil"
	"github.com/pkg/errors"
	"strings"
	"encoding/hex"
)

// http or https
//
// GET    /       - list (plain text, each line - chunk id)
// GET    /:chunk - get chunk (escaped). Must return 200 on success
// DELETE /:chunk - delete chunk (escaped). Must return 200 on success
// POST   /:chunk - put chunk (escaped). Must return non 201 code if already exists

type httpProxy struct {
	baseUrl string
	client  *http.Client
}

func (hp *httpProxy) Put(chunk string, data []byte) error {

	res, err := hp.client.Post(hp.baseUrl+"/"+hex.EncodeToString([]byte(chunk)), "application/octet-stream", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	msg, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusCreated {
		return errors.Errorf("non 201 code: %v %v %v", res.StatusCode, res.Status, string(msg))
	}
	return nil
}

func (hp *httpProxy) Get(chunk string) ([]byte, error) {
	res, err := hp.client.Get(hp.baseUrl + "/" + hex.EncodeToString([]byte(chunk)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("non 200 code: %v %v", res.StatusCode, res.Status)
	}
	return data, nil
}

func (hp *httpProxy) Del(chunk string) (error) {
	req, err := http.NewRequest(http.MethodDelete, hp.baseUrl+"/"+hex.EncodeToString([]byte(chunk)), nil)
	if err != nil {
		return err
	}
	res, err := hp.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	msg, _ := ioutil.ReadAll(res.Body) // for keep-alives
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.Errorf("non 200 code: %v %v %v", res.StatusCode, res.Status, string(msg))
	}
	return nil
}

func (hp *httpProxy) List() ([]string, error) {
	res, err := hp.client.Get(hp.baseUrl + "/")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("non 200 code: %v %v %v", res.StatusCode, res.Status, string(data))
	}
	var uids = make([]string, 0)
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if len(line) != 0 {
			val, err := hex.DecodeString(line)
			if err == nil {
				uids = append(uids, string(val))
			}
		}
	}
	return uids, nil
}
func httpFactory(uri *url.URL) []Storage {
	return []Storage{&httpProxy{uri.String(), new(http.Client)}}
}
func init() {
	Register("http", httpFactory)
	Register("https", httpFactory)
}
