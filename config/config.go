// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

type Config struct {
	Listen        string `config:"listen"`
	Context       string `config:"context"`
	NameUnderRoot bool   `config:"name_under_root"`
}

var DefaultConfig = Config{
	Listen:        ":8080",
	Context:       "/prometheus",
	NameUnderRoot: false,
}
