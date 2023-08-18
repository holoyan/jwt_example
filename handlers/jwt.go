package handlers

import (
	"net/http"
	"fmt"
)

func Index(res http.ResponseWriter, req *http.Request)  {
	fmt.Fprintln(res, "Welcome to the Home Page!")
}
