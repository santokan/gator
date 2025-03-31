package main

import (
	"context"
	"fmt"
)

func handlerAgg(s *state, c command) error {
	if len(c.Args) != 0 {
		return fmt.Errorf("usage: agg")
	}

	rssfeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}

	fmt.Printf("RSS Feed: %+v\n", rssfeed)
	return nil
}
