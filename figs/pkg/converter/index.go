package converter

type Converter interface {
	Convert(source []byte) ([]byte, error)
}
