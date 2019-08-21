package helpers

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// LoadYAML ...
func LoadYAML(filename string, cfg interface{}) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.UnmarshalStrict(b, cfg)
}
