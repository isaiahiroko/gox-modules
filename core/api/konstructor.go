package api

import (
	"fmt"

	"github.com/origine-run/portr/pkg/http"
)

type Konstructor struct {
	host string
}

func (k *Konstructor) GetInstances() (http.JSON, error) {
	res, err := http.Request("GET", fmt.Sprintf("%s/instances", k.host), map[string]any{
		"page":    1,
		"perPage": 12,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
