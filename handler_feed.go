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

	// Automatically create a feed follow record
	followID := uuid.New()
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        followID,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    userID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to create feed follow: %v", err)
	}

	fmt.Printf("Feed added successfully: Name=%s, URL=%s\n", feed.Name, feed.Url)
	fmt.Printf("Automatically following feed '%s' for user '%s'\n", feed.Name, s.cfg.CurrentUserName)

	printFeed(feed)

	return nil
}

func handlerFeeds(s *state, c command) error {
	if len(c.Args) != 0 {
		return fmt.Errorf("usage: feeds")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get feeds: %v", err)
	}

	fmt.Println("List of feeds:")
	for _, feed := range feeds {
		fmt.Printf("* Feed name:        %s\n", feed.Name)
		fmt.Printf("* Feed URL:         %s\n", feed.Url)
		fmt.Printf("* Username:         %s\n", feed.Username)
		fmt.Println()
	}

	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
}
