package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/santokan/gator/internal/database"
)

func handlerFollow(s *state, c command) error {
	if len(c.Args) != 1 {
		return fmt.Errorf("usage: follow <feed_url>")
	}

	feedURL := c.Args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("unable to find feed with URL %s: %v", feedURL, err)
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("unable to get userID for user '%s': %v", s.cfg.CurrentUserName, err)
	}

	followID := uuid.New()
	now := time.Now()

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        followID,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to create feed follow: %v", err)
	}

	fmt.Printf("Following feed '%s' for user '%s'\n", feed.Name, s.cfg.CurrentUserName)
	return nil
}

func handlerFollowing(s *state, c command) error {
	if len(c.Args) != 0 {
		return fmt.Errorf("usage: following")
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("unable to get userID for user '%s': %v", s.cfg.CurrentUserName, err)
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("unable to get feed follows for user '%s': %v", s.cfg.CurrentUserName, err)
	}

	fmt.Println("List of feeds followed by the user:")
	for _, follow := range follows {
		fmt.Printf("* Feed name:        %s\n", follow.FeedName)
		fmt.Printf("* Feed URL:         %s\n", follow.FeedID)
		fmt.Printf("* Username:         %s\n", follow.UserName)
		fmt.Println()
	}

	return nil
}
