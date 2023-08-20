package main

import (
	"./core"
	"./route"
)

func main() {
	route.Register()
	core.Load()
	core.Run()
	core.Close()
}