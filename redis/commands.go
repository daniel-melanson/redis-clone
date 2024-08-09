package redis

type Command struct {
	name string
}

type CommandRegistry interface {
	Get(name string) (command Command, exists bool)
	add(command Command)
}

type CommandStore struct {
	commands map[string]Command
}

func (s CommandStore) add(command Command) {
}

func (s CommandStore) Get(name string) (command Command, exists bool) {
	return Command{}, false
}

func Commands() (commands CommandRegistry) {
	commands = CommandStore{
		make(map[string]Command),
	}

  

  return
}
