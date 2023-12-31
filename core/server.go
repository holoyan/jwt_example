package core

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func Run()  {
	var port = os.Getenv("PORT")
	// string to int
	intPort, strErr := strconv.Atoi(port)
	if strErr != nil {
		// default port
		intPort = 8080
	}

	fmt.Println()
	fmt.Printf("Server started on port %d\n", intPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", intPort), nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
