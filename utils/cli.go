package utils

import (
	"github.com/reddec/vault"
	"io/ioutil"
	"strings"
)

func ReadUrlFiles(wh *vault.SimpleWarehouse, files []string) error {
	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		for _, line := range strings.Split(string(content), "\n") {
			line = strings.TrimSpace(line)
			if len(line) == 0 || line[0] == '#' {
				continue
			}
			wh.Add(line)
		}
	}
	return nil
}
