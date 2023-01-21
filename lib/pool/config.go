package pool

import (
	"github.com/code-to-go/safepool.lib/core"
)

func Define(c Config) error {
	return sqlSave(c.Name, c)
}

func GetConfig(name string) (Config, error) {
	c, err := sqlLoad(name)
	if core.IsErr(err, "cannot load config for pool '%s'", name) {
		return Config{}, err
	}
	return c, nil
}
