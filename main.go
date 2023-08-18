package main

import (
	"./core"
	"./route"
)

func main() {
	route.Register()

	core.Run()
}
