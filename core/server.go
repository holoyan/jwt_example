package core

import (
	"context"
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

func Close()  {
	defer DbClient.Disconnect(context.Background())
}
