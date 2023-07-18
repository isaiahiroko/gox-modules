package validator

import (
	"encoding/json"
	"log"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

type Json struct{}

func (j *Json) Validate(source []byte, schema []byte) error {
	val, err := jsonschema.CompileString("schema.json", string(schema))
	if err != nil {
		return err
	}

	var v interface{}
	if err := json.Unmarshal(source, &v); err != nil {
		log.Fatal(err)
	}

	err = val.Validate(v)
	if err != nil {
		return err
	}

	return nil
}

func NewJson() *Json {
	return &Json{}
}
