package main

import (
	"flag"
	"strconv"

	"dc_assignment.com/sidecar/v2/sidecarroutes"
)

var (
	appPortNumber     = flag.Int("appPortNumber", 8080, "App Port Number")
	sideCarPortNumber = flag.Int("sideCarPortNumber", 0, "Sidecar Port Number")
)

func main() {
	r := sidecarroutes.SetupRouter(strconv.Itoa(*appPortNumber))

	r.Run(":" + strconv.Itoa(*sideCarPortNumber))
}
