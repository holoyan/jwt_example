package core

import (
	"fmt"
	"net/http"
)

func Run()  {
	port := 8082
	fmt.Printf("Server started on port %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
