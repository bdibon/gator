package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bdibon/gator/internal/database"
	"github.com/bdibon/gator/internal/rss"
	"github.com/google/uuid"
)

func handlerAgg(_ *state, _ command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("couldn't fetch rss feed: %w", err)
	}

	fmt.Printf("%#v\n", feed)
	return nil
}

func handlerAddFeed(s *state, c command) error {
	if len(c.args) != 2 {
		return fmt.Errorf("expected 2 arguments got %d", len(c.args))
	}
	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't find user %s: %w", s.cfg.CurrentUserName, err)
	}

	name, url := c.args[0], c.args[1]
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    currentUser.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't save feed to db: %w", err)
	}
	fmt.Printf("Sucessfully created feed: %#v\n", feed)
	return nil
}
