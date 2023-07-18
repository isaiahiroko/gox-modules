package store

type Store interface {
	Store(data []byte) error
}
