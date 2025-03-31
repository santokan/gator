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
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("unable to get userID if user '%s': %v", s.cfg.CurrentUserName, err)
	}

	feed, err := s.db.AddFeed(context.Background(), database.AddFeedParams{
		ID:        feedID,
		Name:      c.Args[0],
		Url:       c.Args[1],
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
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
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to create feed follow: %v", err)
	}

	fmt.Printf("Feed added successfully: Name=%s, URL=%s\n", feed.Name, feed.Url)
	fmt.Printf("Automatically following feed '%s' for user '%s'\n", feed.Name, s.cfg.CurrentUserName)

	printFeed(feed, user)

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

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	fmt.Println("List of feeds:")
	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("unable to get user for feed '%s': %v", feed.Name, err)
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
