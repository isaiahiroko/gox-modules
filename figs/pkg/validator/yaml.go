package validator

import (
	"log"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"gopkg.in/yaml.v3"
)

type Yaml struct{}

func (j *Yaml) Validate(source []byte, schema []byte) error {
	val, err := jsonschema.CompileString("schema.yaml", string(schema))
	if err != nil {
		return err
	}

	var v interface{}
	if err := yaml.Unmarshal(source, &v); err != nil {
		log.Fatal(err)
	}

	err = val.Validate(v)
	if err != nil {
		return err
	}

	return nil
}

func NewYaml() *Yaml {
	return &Yaml{}
}
