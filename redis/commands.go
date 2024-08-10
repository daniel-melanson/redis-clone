package redis

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var ErrorLog *log.Logger

func init() {
	ErrorLog = log.New(os.Stderr, "(error) ", 0)
}

type NetworkBoundary uint64

const (
	Both NetworkBoundary = iota
	Client
	Server
)

type Command struct {
	Handler  func(conn net.Conn, args string)
	Name     string
	Help     string
	Boundary NetworkBoundary
}

type CommandRegistry interface {
	Get(name string) (command Command, exists bool)
	add(command Command)
}

type CommandStore struct {
	commands map[string]Command
}

func sanitize(raw string) (key string) {
	key = strings.Trim(strings.ToUpper(raw), " \n\r")
	return
}

func (store CommandStore) add(command Command) {
	key := sanitize(command.Name)
	store.commands[key] = command
}

func (store CommandStore) Get(name string) (command Command, exists bool) {
	key := sanitize(name)
	command, exists = store.commands[key]
	return
}

func Commands() (commands CommandRegistry) {
	commands = CommandStore{
		make(map[string]Command),
	}

	commands.add(Command{
		Name:     "help",
		Boundary: Client,
		Help:     "Display help.",
		Handler: func(_ net.Conn, args string) {
			if len(args) == 0 {
				fmt.Println(`redis-cli
  To get help about Redis commands type:
      "help <command>" for help on <command>
      "quit" to exit`)
			} else if args[0] == '@' {
				// TODO: Implement
				return
			} else {
				command, exists := commands.Get(args)

				if exists {
					fmt.Println(command.Help)
				} else {
					ErrorLog.Printf("unknown command '%s'", args)
				}
			}
		},
	})

	commands.add(Command{
		Name:     "quit",
		Boundary: Client,
		Help:     "Closes the client.",
		Handler: func(conn net.Conn, _ string) {
			if conn != nil {
				conn.Close()
			}
			os.Exit(0)
		},
	})

	commands.add(Command{})

	return
}
