package main

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/bdibon/gator/internal/database"
	"github.com/bdibon/gator/internal/rss"
	"github.com/google/uuid"
)

func handlerFollow(s *state, c command) error {
	if len(c.args) < 1 {
		return errors.New("missing argument: feed url")
	}

	currentUser, err := s.db.GetUserByName(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't find current user: %w", err)
	}

	feedUrl := c.args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)
	if err != nil {
		if err != sql.ErrNoRows {
			return fmt.Errorf("couldn't check if feed already exists: %w", err)
		}
	}

	if err == sql.ErrNoRows {
		newFeed, err := rss.FetchFeed(context.Background(), feedUrl)
		if err != nil {
			return fmt.Errorf("couldn't fetch new feed data: %w", err)
		}

		feed, err = s.db.CreateFeed(context.Background(), database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      newFeed.Channel.Title,
			Url:       newFeed.Channel.Link,
			UserID:    currentUser.ID,
		})
		if err != nil {
			return fmt.Errorf("couldn't create new feed: %w", err)
		}
	}
	return createFeedFollows(s, feed, currentUser)
}

func createFeedFollows(s *state, feed database.Feed, user database.User) error {
	feedFollows, err := s.db.CreateFeedFollows(context.Background(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldnt' create new feed_follows: %w", err)
	}

	fmt.Printf("%s now following %s\n", feedFollows.Username, feedFollows.Feedname)
	return nil
}

func handlerFollowing(s *state, c command) error {
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		if err != sql.ErrNoRows {
			return fmt.Errorf("couldn't retrieve feed_follows: %w", err)
		}
		fmt.Printf("%s doesn't follow any feed\n", s.cfg.CurrentUserName)
	}

	writer := bufio.NewWriter(os.Stdout)
	writer.WriteString(fmt.Sprintf("%s follows:\n", s.cfg.CurrentUserName))
	for _, feedFollow := range feedFollows {
		writer.WriteString(fmt.Sprintf("* %s\n", feedFollow.Feedname))
	}
	writer.Flush()
	return nil
}
