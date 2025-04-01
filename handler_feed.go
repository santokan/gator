package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/santokan/gator/internal/database"
)

func handlerAddFeed(s *state, c command, user database.User) error {
	if len(c.Args) != 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}

	feedID := uuid.New()
	now := time.Now()

	feed, err := s.db.AddFeed(context.Background(), database.AddFeedParams{
		ID:        feedID,
		Name:      c.Args[0],
		Url:       c.Args[1],
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to add feed: %w", err)
	}

	fmt.Println("Feed added successfully:")

	// Automatically create a feed follow record
	followID := uuid.New()
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        followID,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to create feed follow: %w", err)
	}

	fmt.Println("")
	printFeed(feed, user)
	fmt.Println("")
	fmt.Println("Feed follow created successfully!")

	return nil
}

func handlerFeeds(s *state, c command) error {
	if len(c.Args) != 0 {
		return fmt.Errorf("usage: feeds")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	fmt.Println("List of feeds:")
	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("unable to get user for feed '%s': %w", feed.Name, err)
		}
		printFeed(feed, user)
	}

	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:          %s\n", user.Name)
}
