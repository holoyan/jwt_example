package handlers

import (
	"fmt"
	"net/http"
)

func IssueTokens(res http.ResponseWriter, req *http.Request)  {
	fmt.Fprintln(res, "Access and refresh token")
}

func UpdateTokens(res http.ResponseWriter, req *http.Request)  {
	fmt.Fprintln(res, "Update tokens")

	fmt.Println("Document inserted successfully!")
}
