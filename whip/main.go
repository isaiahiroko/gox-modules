/*
Copyright © 2022 Isaiah Iroko <isaiahiroko@gmail.com>
*/
package main

import (
	"github.com/origine-run/whip/cmd"
	"github.com/origine-run/whip/pkg/store"
	"github.com/origine-run/whip/pkg/task"
)

func main() {
	// be ready to run task when added
	// runs forever, until the service is shutdown
	go task.Run()

	// open data store
	store.Open()
	defer store.Close()

	cmd.Execute()
}
