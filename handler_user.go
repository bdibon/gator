package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bdibon/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing argument: username")
	}
	username := cmd.args[0]

	_, err := s.db.GetUser(
		context.Background(),
		username,
	)
	if err != nil {
		return fmt.Errorf("%s doesn't exist", username)
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("error setting user: %w", err)
	}
	fmt.Printf("username was set to %s\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing argument: username")
	}
	username := cmd.args[0]

	_, usrCheck := s.db.GetUser(
		context.Background(),
		username,
	)
	if usrCheck == nil {
		return errors.New("user already exists")
	}

	usr, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      username,
		})
	if err != nil {
		return fmt.Errorf("create user failed: %w", err)
	}
	fmt.Println("user was created:")
	printUser(usr)

	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Printf("switched to new user: %s\n", usr.Name)
	return nil
}

func printUser(usr database.User) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(" * ID:		%v\n", usr.ID))
	sb.WriteString(fmt.Sprintf(" * Name:	%v\n", usr.Name))
	fmt.Print(sb.String())
}
