package kv

var (
	buffer = 1000
)

const (
	batchSize = 10
)

type Store struct {
	container chan []byte
}

func (s *Store) Push(values ...[]byte) {
	go func() {
		for _, value := range values {
			s.container <- value
		}
	}()
}

func (s *Store) Pop(len int) {
	batch := [batchSize][]byte{}

	for f := range s.container {
		batch[0] = f
	L:
		for i := 1; i < batchSize; i++ {
			select {
			case batch[i] = <-s.container:
			default:
				// continue
				break L
			}

		}

		//do sth with batch
	}
}

func New() *Store {
	return &Store{
		container: make(chan []byte, buffer),
	}
}
