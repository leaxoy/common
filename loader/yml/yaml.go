package yml

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Loader struct{}

func (l *Loader) Load(f string, v interface{}) error {
	buf, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buf, v)
	if err != nil {
		return err
	}
	return l.Verify(v)
}

func (*Loader) Verify(interface{}) error {
	return nil
}
