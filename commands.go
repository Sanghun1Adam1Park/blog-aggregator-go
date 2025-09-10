package main

import (
	"fmt"
)

type commands struct {
	commandMap map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandMap[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if f, exist := c.commandMap[cmd.name]; exist {
		if err := f(s, cmd); err != nil {
			return err
		}
		return nil
	} else {
		return fmt.Errorf("there is no such commad '%s'", cmd.name)
	}
}
