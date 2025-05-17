package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var globalCache = make(map[string]interface{})
var ttlMap = make(map[string]time.Time)
var mu sync.RWMutex

func handleGET(args []string, conn net.Conn) {
	if len(args) != 2 {
		conn.Write([]byte("ERR usage: GET <key>\n"))
		return
	}
	key := args[1]
	if time.Now().After(ttlMap[key]) {
		delete(ttlMap, key)
		delete(globalCache, key)
		conn.Write([]byte("Key expired\n"))
		return
	}
	mu.RLock()
	value, ok := globalCache[key]
	mu.RUnlock()
	if !ok {
		conn.Write([]byte("nil\n"))
		return
	}
	conn.Write([]byte(fmt.Sprintf("%v\n", value)))
}

func handleSET(args []string, conn net.Conn) {
	if len(args) != 3 {
		conn.Write([]byte("ERR usage: SET <key> <value>\n"))
		return
	}
	key := args[1]
	value := args[2]
	mu.Lock()
	globalCache[key] = value
	mu.Unlock()
	conn.Write([]byte("OK\n"))
}

func handleDEL(args []string, conn net.Conn) {
	if len(args) != 2 {
		conn.Write([]byte("ERR usage: DEL <key>\n"))
		return
	}
	key := args[1]
	mu.Lock()
	delete(globalCache, key)
	mu.Unlock()
	conn.Write([]byte("Deleted key\n"))
}

func handleEXPIRE(args []string, conn net.Conn) {
	if len(args) != 3 {
		conn.Write([]byte("ERR usage: EXPIRE <key> <seconds>\n"))
		return
	}
	key := args[1]
	seconds, err := strconv.Atoi(args[2])
	if err != nil {
		conn.Write([]byte("ERR usage: EXPIRE <key> <seconds>\n"))
		return
	}
	mu.Lock()
	ttlMap[key] = time.Now().Add(time.Duration(seconds) * time.Second)
	mu.Unlock()
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
		args := strings.Fields(message)

		switch operation {
		case "GET":
			handleGET(args, conn)
		case "SET":
			handleSET(args, conn)
		case "DEL":
			handleDEL(args, conn)
		case "EXPIRE":
			handleEXPIRE(args, conn)
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
