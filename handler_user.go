package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/santokan/gator/internal/database"
)

func handlerLogin(s *state, c command) error {
	if len(c.Args) != 1 {
		return fmt.Errorf("usage: login <username>")
	}

	username := c.Args[0]
	err := s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("error setting user: %w", err)
	}

	_, err = s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("unable to login with user '%s': %w", username, err)
	}
	fmt.Printf("User '%s' has been set.\n", username)
	return nil
}

func handlerRegister(s *state, c command) error {
	if len(c.Args) != 1 {
		return fmt.Errorf("usage: register <username>")
	}

	username := c.Args[0]
	_, err := s.db.GetUser(context.Background(), username)
	if err == nil { // If no error, it means a user with that name exists
		return fmt.Errorf("user with username '%s' already exists", username)
	}

	if err != sql.ErrNoRows {
		return fmt.Errorf("error checking username: %w", err)
	}

	userID := uuid.New()
	now := time.Now()

	newUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        userID,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      username,
	})
	if err != nil {
		return fmt.Errorf("failed to create user '%s': %w", username, err)
	}

	s.cfg.CurrentUserName = newUser.Name

	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("error setting user: %w", err)
	}

	log.Printf("User '%s' has been successfully registered.", username)
	printUser(newUser)

	return nil
}

func handlerUsers(s *state, c command) error {
	if len(c.Args) != 0 {
		return fmt.Errorf("usage: users")
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching users: %w", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf(" * %v (current)\n", user)
		} else {
			fmt.Printf(" * %v\n", user)
		}
	}

	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}
