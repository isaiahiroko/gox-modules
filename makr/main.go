/*
Copyright © 2022 Isaiah Iroko <isaiahiroko@gmail.com>
*/
package main

import (
	"github.com/origine-run/makr/cmd"
	"github.com/origine-run/makr/pkg/store"
)

func main() {
	store.Open()
	defer store.Close()

	cmd.Execute()
}
