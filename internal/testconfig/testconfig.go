package testconfig

type Config struct {
	Tests []Test `yaml:"tests"`
}

type Test struct {
	Memo         string `yaml:"memo"`
	Postcode     int    `yaml:"postcode"`
	Address      string `yaml:"address"`
	Organisation string `yaml:"organisation"`
	Description  string `yaml:"description"`
	State        string `yaml:"state"`
}
