package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/daniel-melanson/redis-clone/redis"
)

var ErrorLog *log.Logger

func makeConnection(host string, port int) (net.Conn, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return nil, fmt.Errorf("could not connect to Redis at %s:%d: %w", host, port, err)
	}

	return conn, nil
}

func splitLine(line string) (commandName string, rawArgs string) {
	line = strings.Trim(line, " \r\n")
	first_word_index := strings.Index(line, " ")
	if first_word_index != -1 {
		commandName = line[:first_word_index]
		rawArgs = line[first_word_index+1:]
	} else {
		commandName = line
		rawArgs = ""
	}

	return
}

func init() {
	ErrorLog = log.New(os.Stderr, "(error) ", 0)
}

func main() {
	var (
		host string
		port int
	)

	flag.StringVar(&host, "H", "127.0.0.1", "redis server host (default: 127.0.0.1)")
	flag.IntVar(&port, "P", 6379, "redis server host port (default: 6379)")

	flag.Parse()

	conn, connErr := makeConnection(host, port)
	connName := host
	if connErr != nil {
		connName = "not connected"
		fmt.Println(connErr)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s> ", connName)
		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		commandName, rawArgs := splitLine(line)

		command, exists := redis.Registry.Get(commandName)

		if !exists {
			ErrorLog.Printf("ERR unknown command '%s' with arguments: %s\n", commandName, rawArgs)
		} else if command.Boundary == redis.Client {
			command.Handler(conn, rawArgs)
		} else if conn == nil {
			fmt.Println(connErr)
		} else {
			command.Handler(conn, rawArgs)
		}
	}

	fmt.Println()
}
