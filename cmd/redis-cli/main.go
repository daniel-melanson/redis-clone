package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	// "strings"
)

func connect_redis(host string, port int) (net.Conn, string) {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", addr)
	// Connection is made
	if err != nil {
		msg := fmt.Sprint(err)

		re := regexp.MustCompile(`^dial tcp [^:]+:[^:]+: .+: (.+)$`)
		m := re.FindAllStringSubmatch(msg, -1)
		// Matches a normal dial error
		if len(m) > 0 && len(m[0]) > 1 {
			return nil, fmt.Sprintf("Could not connect to Redis at %s: %s", addr, m[0][1])
		}

		return nil, msg
	}

	return conn, ""
}

func main() {
	var (
		host string
		port int
	)

	flag.StringVar(&host, "H", "127.0.0.1", "redis server host (default: 127.0.0.1)")
	flag.IntVar(&port, "P", 6379, "redis server host port (default: 6379)")

	flag.Parse()

	conn, err := connect_redis(host, port)

	var prefix string
	if conn == nil {
		fmt.Println(err)

		prefix = "not connected"
	} else {
		prefix = host
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s> ", prefix)
		if !scanner.Scan() {
			break
		}

		// line := scanner.Text()
		if conn == nil {
			fmt.Println(err)
		} else {
			conn.Write([]byte("+OK\r\n"))
		}
	}

	fmt.Println()
}
