package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/bdibon/gator/internal/database"
)

func handleBrowse(s *state, c command, user database.User) error {
	var limit int32 = 2
	if len(c.args) > 0 {
		customLimit, err := strconv.Atoi(c.args[0])
		if err != nil {
			return fmt.Errorf("invalid integer: %s", c.args[0])
		}
		limit = int32(customLimit)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		ID:    user.ID,
		Limit: limit,
	})
	if err != nil {
		return fmt.Errorf("couldn't retrieve posts from db: %w", err)
	}

	writer := bufio.NewWriter(os.Stdout)
	writer.WriteString(fmt.Sprintf("Hey %s, here's your posts:\n\n", user.Name))
	for _, post := range posts {
		writer.WriteString(fmt.Sprintf("%s from %s \n\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName))
		writer.WriteString(fmt.Sprintf("--- %s ---\n", post.Title))
		writer.WriteString(fmt.Sprintf("\t%s\n\n", post.Description))
	}
	writer.Flush()

	return nil
}
