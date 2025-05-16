package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

var globalCache = make(map[string]interface{})

func handleGET() {

}

func handleSET() {

}

func handleDEL() {

}

func handleEXPIRE() {

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading from connection", err)
		}

		operation := strings.Split(strings.ToUpper(message), " ")[0]
		switch operation {
		case "GET":
			// Handle Get

			handleGET()
		case "SET":

			handleSET()
		case "DEL":
			handleDEL()
		case "EXPIRE":
			handleEXPIRE()
		default:
			log.Println("Unknown operation", operation)
		}
	}
}

func main() {
	fmt.Println("Server starting...")

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn)
	}
}
