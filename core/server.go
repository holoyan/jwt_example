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
	intPort, error := strconv.Atoi(port)
	if error != nil {
		// ... handle error
		intPort = 8085
	}

	fmt.Println()
	fmt.Printf("Server started on port %d\n", intPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", intPort), nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
