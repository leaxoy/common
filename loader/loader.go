package loader

type Loader interface {
	Name() string
	Load(string, interface{}) error
}

var loaderMap = make(map[string]Loader)

func Register(loader Loader) {
	loaderMap[loader.Name()] = loader
}

func Get(name string) Loader {
	l, ok := loaderMap[name]
	if ok {
		return l
	}
	return nil
}
