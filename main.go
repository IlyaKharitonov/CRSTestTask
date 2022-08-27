package main

import (
	"fmt"
	"log"
	"net/http"

	"testTask/api"
)

func main() {
	api.RegisterHandlers()
	log.Printf("RegisterHandlers start")

	addr := "127.0.0.1:1616"
	fmt.Printf("Server Start in addr %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Printf("Server error: %s", err)
	}

	fmt.Println("Server stop")

}
