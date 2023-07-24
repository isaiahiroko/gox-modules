package logger

import (
	"io"

	"github.com/origine-run/core/store"
	"github.com/origine-run/core/utils"
)

type logger struct {
	name string
	w    io.Writer
}

func (l logger) Write(buf []byte) (int, error) {
	logger := store.New(l.name)

	logger.Create(utils.Now(), string(buf))

	n, err := l.w.Write(buf)
	if err != nil {
		return n, err
	}
	if n != len(buf) {
		return n, io.ErrShortWrite
	}
	return len(buf), nil
}

func (l logger) Add(message string) {
	log := store.New(l.name)
	log.Create(utils.Now(), message)
}

func New(name string, w io.Writer) *logger {
	return &logger{name, w}
}
