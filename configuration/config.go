package configuration

import (
	"os"

	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Unifi  []*UnfiConfig   `yaml:"unifi,omitempty"`
	Ruckus []*RuckusConfig `yaml:"ruckus,omitempty"`
}

type UnfiConfig struct {
	ApiUrl  string `yaml:"api_url"`
	ApiUser string `yaml:"api_user"`
	ApiPass string `yaml:"api_password"`
}

type RuckusConfig struct {
	Host string `yaml:"host"`
	Site string `yaml:"host"`
}

// Load reads config file
func Load(path string) (*Config, error) {
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	defer f.Close()

	return loadFile(f)
}

func loadFile(f *os.File) (*Config, error) {
	b, err := ioutil.ReadAll(f)

	if err != nil {
		return nil, err
	}

	return unmarshal(b)
}

func unmarshal(b []byte) (*Config, error) {
	var c Config
	err := yaml.Unmarshal(b, &c)

	if err != nil {
		return nil, err
	}

	return &c, nil
}
