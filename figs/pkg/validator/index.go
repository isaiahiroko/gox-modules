package validator

type Validator interface {
	Validate(source []byte, schema []byte) error
}
