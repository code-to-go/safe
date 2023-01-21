package pool

func Save(name string, config Config) error {
	return sqlSave(name, config)
}

func Load(name string) (Config, error) {
	return sqlLoad(name)
}
