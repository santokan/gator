package main

import "fmt"

func handlerLogin(s *state, c command) error {
	if len(c.Args) != 1 {
		return fmt.Errorf("usage: login <username>")
	}

	username := c.Args[0]
	err := s.config.SetUser(username)
	if err != nil {
		return fmt.Errorf("error setting user: %v", err)
	}
	fmt.Printf("User '%s' has been set.\n", username)
	return nil
}
