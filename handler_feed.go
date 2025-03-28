package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/santokan/gator/internal/database"
)

func handlerAddFeed(s *state, c command) error {
	if len(c.Args) != 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}

	feedID := uuid.New()
	now := time.Now()
	userID, err := s.db.GetUserIDbyUsername(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("unable to get userID if user '%s': %v", s.cfg.CurrentUserName, err)
	}

	feed, err := s.db.AddFeed(context.Background(), database.AddFeedParams{
		ID:        feedID,
		Name:      c.Args[0],
		Url:       c.Args[1],
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    userID,
	})
	if err != nil {
		return fmt.Errorf("unable to add feed: %v", err)
	}

	fmt.Printf("Feed added successfully: Name=%s, URL=%s\n", feed.Name, feed.Url)
	return nil
}
