package yml

import (
	"io/ioutil"

	"github.com/leaxoy/common/loader"
	"gopkg.in/yaml.v2"
)

type Loader struct{}

func (*Loader) Load(f string, v interface{}) error {
	buf, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buf, v)
	if err != nil {
		return err
	}
	if verifier, ok := v.(loader.Verifier); ok {
		return verifier.Verify()
	}
	return nil
}
