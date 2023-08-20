package route

import (
	"../handlers"
	"net/http"
)

func Register()  {
	http.HandleFunc("/tokens", handlers.IssueTokens)
	http.HandleFunc("/tokens/update", handlers.UpdateTokens)
}