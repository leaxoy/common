package yml

import (
	"io/ioutil"

	"github.com/leaxoy/common/loader"
	"gopkg.in/yaml.v2"
)

func init() {
	loader.Register(new(l))
}

const Name = "yaml"

type l struct{}

func (*l) Name() string {
	return Name
}

func (*l) Load(f string, v interface{}) error {
	buf, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buf, v)
	if err != nil {
		return err
	}
	return nil
}
