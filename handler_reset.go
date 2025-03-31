package main

import (
	"context"
	"fmt"
	"log"
)

func handlerReset(s *state, c command) error {
	if len(c.Args) != 0 {
		return fmt.Errorf("usage: reset")
	}
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error resetting database: %w", err)
	}
	log.Println("Database has been reset.")
	return nil
}
