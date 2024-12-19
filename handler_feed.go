package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/bdibon/gator/internal/database"
	"github.com/bdibon/gator/internal/rss"
	"github.com/google/uuid"
)

func handlerAgg(s *state, c command) error {
	if len(c.args) < 1 {
		return errors.New("missing argument: <time_between_req>")
	}

	timeBetweenRequests, err := time.ParseDuration(c.args[0])
	if err != nil {
		return fmt.Errorf("couldn't parse <time_between_req>: %w", err)
	}

	ticker := time.NewTicker(timeBetweenRequests)
	for range ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return fmt.Errorf("error while scrapping feeds: %w", err)
		}
	}

	return nil
}

func handlerAddFeed(s *state, c command, user database.User) error {
	if len(c.args) != 2 {
		return errors.New("expected 2 arguments <feedname, feedurl>")
	}

	name, url := c.args[0], c.args[1]
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't save feed to db: %w", err)
	}
	fmt.Printf("Sucessfully created feed: %#v\n", feed)
	return createFeedFollows(s, feed, user)
}

func handlerFeeds(s *state, _ command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't retrieve feeds from db: %w", err)
	}

	ownerCache := map[string]string{}
	writer := bufio.NewWriter(os.Stdout)
	for _, feed := range feeds {
		ownerName, ok := ownerCache[feed.UserID.String()]
		if !ok {
			usr, err := s.db.GetUserById(context.Background(), feed.UserID)
			if err != nil {
				return fmt.Errorf("no matching user for id %s: %w", feed.UserID.String(), err)
			}
			ownerCache[feed.UserID.String()] = usr.Name
			ownerName = usr.Name
		}

		writer.WriteString(fmt.Sprintf("* Name: %s\n", feed.Name))
		writer.WriteString(fmt.Sprintf("* URL: %s\n", feed.Url))
		writer.WriteString(fmt.Sprintf("* Owner: %s\n", ownerName))
		writer.WriteString("\n")
	}
	writer.Flush()
	return nil
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get next feed: %w", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID:        feed.ID,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("couldn't mark feed \"%s\" as fetched: %w", feed.Name, err)
	}

	freshFeed, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("couldn't fetch feed with url %s: %w", feed.Url, err)
	}
	printFeedItem(freshFeed)
	return nil
}

func printFeedItem(feed *rss.RSSFeed) {
	writer := bufio.NewWriter(os.Stdout)
	writer.WriteString(fmt.Sprintf("%s's items:\n", feed.Channel.Title))
	for _, item := range feed.Channel.Item {
		writer.WriteString(fmt.Sprintf("\t* %s\n", item.Title))
	}
	writer.Flush()
}
