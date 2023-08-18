package route

import (
	"../handlers"
	"net/http"
)

func Register()  {
	http.HandleFunc("/", handlers.Index)
}