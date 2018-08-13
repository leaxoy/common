package toml

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/leaxoy/common/loader"
)

func init() {
	loader.Register(new(l))
}

const Name = "toml"

type l struct{}

func (*l) Name() string {
	return Name
}

func (*l) Load(f string, v interface{}) error {
	buf, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	err = toml.Unmarshal(buf, v)
	if err != nil {
		return err
	}
	if verifier, ok := v.(loader.Verifier); ok {
		return verifier.Verify()
	}
	return nil
}
