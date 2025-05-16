package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	fmt.Println("Client starting...")

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
}
