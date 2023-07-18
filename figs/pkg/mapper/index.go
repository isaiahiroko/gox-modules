package mapper

type Mapper interface {
	Map(source []byte) ([]byte, error)
}
