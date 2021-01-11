package parser

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func ParseFile(data interface{}, path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrapf(err, "parser: ParseFile -> ioutil.ReadFile(%v)", path)
	}
	err = yaml.Unmarshal(file, data)
	if err != nil {
		return errors.Wrapf(err, "parser: ParseFile -> yaml.Unmarshal(%v, %v)", file, data)
	}
	return nil
}
