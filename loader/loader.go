package loader

type Loader interface {
	Load(string, interface{}) error
}

type Verifier interface {
	Verify() error
}

func LoadWithLoader(f string, v interface{}, loader Loader) error {
	return loader.Load(f, v)
}
