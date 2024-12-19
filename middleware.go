package main

import (
	"context"
	"fmt"

	"github.com/bdibon/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, c command, user database.User) error) func(s *state, c command) error {
	return func(s *state, c command) error {
		currentUser, err := s.db.GetUserByName(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("couldn't find user %s: %w", s.cfg.CurrentUserName, err)
		}
		return handler(s, c, currentUser)
	}
}
