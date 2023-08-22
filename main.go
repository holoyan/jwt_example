package main

import (
	"JwtTask/core"
	"JwtTask/route"
)

func main() {
	route.Register()
	core.Load()
	core.Run()
	core.Close()
}
