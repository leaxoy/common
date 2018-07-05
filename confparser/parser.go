package confparser

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func Load(file string, v interface{}) error {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(buf, v)
}
