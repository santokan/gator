package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/santokan/gator/internal/database"
)

func handlerFollow(s *state, c command, user database.User) error {
	if len(c.Args) != 1 {
		return fmt.Errorf("usage: follow <feed_url>")
	}

	feedURL := c.Args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("unable to find feed with URL %s: %v", feedURL, err)
	}

	followID := uuid.New()
	now := time.Now()

	feedFollowRow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        followID,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to create feed follow: %v", err)
	}

	fmt.Println("Feed follow created successfully:")
	printFollowFeed(feedFollowRow.UserName, feedFollowRow.FeedName)
	return nil
}

func handlerFollowing(s *state, c command, user database.User) error {
	if len(c.Args) != 0 {
		return fmt.Errorf("usage: following")
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("unable to get feed follows for user '%s': %v", s.cfg.CurrentUserName, err)
	}

	if len(follows) == 0 {
		fmt.Printf("User '%s' is not following any feeds.\n", s.cfg.CurrentUserName)
		return nil
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

func printFollowFeed(username, feedname string) {
	fmt.Printf("* User:        %s\n", username)
	fmt.Printf("* Feed:			   %s\n", feedname)
}
