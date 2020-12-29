package parser

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func ParseFile(data interface{}, path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrapf(err, "ParseFile -> ioutil.ReadFile(%v)", path)
	}
	err = yaml.Unmarshal(file, data)
	if err != nil {
		return errors.Wrapf(err, "ParseFile -> yaml.Unmarshal(%v, %v)", file, data)
	}
	return nil
}
