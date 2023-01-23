package shared

import "log"

func AddCommand(commands map[string]func() int, cmd string, action func() int) {
	commands[cmd] = action
}

func DispatchCommand(commands map[string]func() int, cmd string) int {
	if f, ok := commands[cmd]; ok {
		return f()
	} else {
		log.Printf("Unsupported command")
		return 400
	}
}